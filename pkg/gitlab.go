package pkg

import (
	"context"
	crdclientset "dt-runner/generated/clientset/versioned"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/webhooks/gitlab"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// ci.name:pod
var jobCache = make(map[string]corev1.Pod)

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
	klog.Info("dtjob name: ", dtJob.name)
	pod, err := GeneratePod(dtJob)
	if err != nil {
		klog.Errorln("GeneratePod error, %s", err)
	}
	if dtJob.name != "" && pod.Name != "" {
		jobCache[dtJob.ci] = pod
	}
}

// startPod is used to start a pod if confirmed
// startPod will keep check key:value in jobCache
func StartPod() {

	ciclient := crdclientset.NewForConfigOrDie(config)
	cilist, err := ciclient.AppsV1().Cis(DefaultNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln("list cis error, ", err)
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
		jobCache[ci.Name] = pod
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Errorln("build clientset error, ", err)
		return
	}
	klog.Info("begin checking ci and pods")
	pods, err := clientset.CoreV1().Pods(DefaultNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln("list pods error, ", err)
		return
	}

	podsClient := clientset.CoreV1().Pods(DefaultNamespace)
	items := make(map[string]string)
	for _, v := range pods.Items {
		names := strings.Split(v.Name, "-")
		if len(names) == 2 {
			items[names[1]] = v.Name
		}
	}
	for k, v := range jobCache {
		klog.InfoS("check job: ", k)
		if _, ok := items[k]; ok {
			klog.InfoS("pod existes, %s", v.Name)
			continue
		}
		klog.InfoS("pod info, %s", v)
		result, err := podsClient.Create(context.TODO(), &v, metav1.CreateOptions{})
		if err != nil {
			klog.Error("pod create error, %s", err.Error())
			continue
		}
		jobCache[k] = *result
	}
	klog.Info("end checking ci and pods")
}

// generateDtJob is used to generate DtJob
func generateDtJob(payload interface{}) (dtJob DtJob, err error) {
	var name = RandomString(6)
	klog.Info("generateDtJob", name)
	dtJob.name = name
	switch payload.(type) {
	case gitlab.PushEventPayload:
		push := payload.(gitlab.PushEventPayload)
		dtJob.branch = "main"
		dtJob.httpurl = push.Repository.URL
		dtJob.sshurl = push.Project.SSHURL
		dtJob.ref = push.Ref
		dtJob.checkoutSHA = push.CheckoutSHA
	case gitlab.TagEventPayload:
		tag := payload.(gitlab.TagEventPayload)
		dtJob.branch = "main"
		dtJob.checkoutSHA = tag.CheckoutSHA
		dtJob.httpurl = tag.Repository.URL
		dtJob.ref = tag.Ref
		dtJob.sshurl = tag.Project.SSHURL
	case gitlab.SystemHookPayload:
		fmt.Println("system hook")
		return DtJob{}, nil
	}
	klog.Info("dtjob name: ", dtJob.name)
	return dtJob, nil
}
