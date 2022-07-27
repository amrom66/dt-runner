package pkg

import (
	"context"
	crdclientset "dt-runner/generated/clientset/versioned"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/webhooks/gitlab"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// gitlab event -> jobCache
// schedule task StartPod is used to start pod from jobCache

// job.name:pod
var jobCache = make(map[string]corev1.Pod)

var ciCache *cache.Cache

func init() {
	ciCache = cache.New(1*time.Minute, 2*time.Minute)
}

// GitlabHook is used to hook gitlab
func GitlabHook(w http.ResponseWriter, r *http.Request) {
	secret := viper.GetString("webhook.token")
	hook, _ := gitlab.New(gitlab.Options.Secret(secret))
	payload, err := hook.Parse(r, gitlab.PushEvents, gitlab.TagEvents, gitlab.SystemHookEvents)
	if err != nil {
		klog.Errorln("start gitlabhook error, %s", err)
		return
	}
	dtJob, error := generateDtJob(payload)
	if error != nil {
		klog.Errorln("generateDtJob error, %s", err)
	}

	dtJobs := attachDtJob(dtJob)

	for _, job := range dtJobs {
		pod, err := GeneratePod(job)
		if err != nil {
			klog.Errorln("GeneratePod error, %s", err)
		}
		//todo 矫正，连续触发限流 1分钟1次触发
		value, found := ciCache.Get(job.ci)
		if found {
			klog.Info("ci still in cache, this happens may because of quick trigger: ", value)
		} else if job.name != "" && pod.Name != "" {
			ciCache.Set(job.ci, job.name, cache.DefaultExpiration)
			klog.Info("ci will be cached: ", job.ci)
			jobCache[job.name] = pod
			//ci update
			UpdateCi(DefaultNamespace, job.ci, pod.Name, "SCHEDULED")
		}
	}
}

// startPod is used to start a pod if confirmed
// startPod will keep check key:value in jobCache
func StartPod() {

	//每个ci的默认pod
	ciclient := crdclientset.NewForConfigOrDie(config)
	cilist, err := ciclient.AppsV1().Cis(DefaultNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln("list cis error: ", err)
		return
	}
	for _, ci := range cilist.Items {
		klog.Info("ci name: ", ci.Name)
		dtjob := DtJob{
			name:   "init",
			ci:     ci.Name,
			branch: "main",
		}
		pod, err := GeneratePod(dtjob)
		if err != nil {
			klog.Errorln("GeneratePod error: ", err)
			continue
		}
		klog.Info("Generate init pod: ", pod.Name)
		// 存入缓存中
		jobCache[dtjob.name] = pod
	}

	// 矫正已有的pod数据
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Errorln("build clientset error, ", err)
		return
	}
	pods, err := clientset.CoreV1().Pods(DefaultNamespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: DefaultLabelDtRunner + "=true",
	})
	if err != nil {
		klog.Errorln("list pods error, ", err)
		return
	}

	//poditems 是pod数据快照
	poditems := make(map[string]struct{})
	for _, v := range pods.Items {
		poditems[v.Name] = struct{}{}
	}

	//启动pod 如果pod已经存在，则跳过
	podsClient := clientset.CoreV1().Pods(DefaultNamespace)
	for k, v := range jobCache {
		klog.Info("check job: ", k)
		if _, ok := poditems[v.Name]; ok {
			klog.Info("pod existes: ", v.Name)
			continue
		}
		_, err := podsClient.Create(context.TODO(), &v, metav1.CreateOptions{})
		if err != nil {
			klog.Error("pod create error, %s", err.Error())
			continue
		}
		UpdateCi(DefaultNamespace, v.Labels[DefaultLabelDtRunnerCi], v.Name, "STARTED")
		klog.Info("pod start success", v.Name)
		delete(jobCache, k)
	}
	klog.Info("end checking ci and pods")
}

// generateDtJob is used to generate DtJob
func generateDtJob(payload interface{}) (dtJob DtJob, err error) {
	switch payload.(type) {
	case gitlab.PushEventPayload:
		klog.Info("event type: ", "push")
		push := payload.(gitlab.PushEventPayload)
		dtJob.branch = "main"
		dtJob.httpurl = push.Project.GitHTTPURL
		dtJob.sshurl = push.Project.SSHURL
		dtJob.ref = push.Ref
		dtJob.checkoutSHA = push.CheckoutSHA
		dtJob.project = push.Project.Name
	case gitlab.TagEventPayload:
		klog.Info("event type: ", "tag")
		tag := payload.(gitlab.TagEventPayload)
		dtJob.branch = "main"
		dtJob.checkoutSHA = tag.CheckoutSHA
		dtJob.httpurl = tag.Project.GitHTTPURL
		dtJob.ref = tag.Ref
		dtJob.sshurl = tag.Project.SSHURL
		dtJob.project = tag.Project.Name
	case gitlab.SystemHookPayload:
		klog.Info("event type: ", "system")
		klog.Info("system hook")
		return DtJob{}, errors.New("system hook")
	}
	return dtJob, nil
}

// attach ci info to dtjob
func attachDtJob(dtJob DtJob) []DtJob {
	var dtJobs []DtJob
	ciList := ListCis(DefaultNamespace).Items
	for _, ci := range ciList {
		if ci.Spec.Repo == dtJob.httpurl {
			var name = RandStringRunes(6)
			dtJob.name = name
			dtJob.ci = ci.Name
			dtJobs = append(dtJobs, dtJob)
		}
	}
	return dtJobs
}
