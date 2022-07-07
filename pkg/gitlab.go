package pkg

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-playground/webhooks/gitlab"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

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
	pod, err := GeneratePod(dtJob)
	if err != nil {
		klog.Errorln("GeneratePod error, %s", err)
	}
	if dtJob.name != "" && pod.Name != "" {
		jobCache[dtJob.name] = pod
	}
}

// startPod is used to start a pod if confirmed
// startPod will keep check key:value in podCache
func StartPod() {
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

	items := make(map[string]corev1.Pod)
	for _, v := range pods.Items {
		items[v.Name] = v
	}
	for k, v := range jobCache {
		klog.InfoS("check job %s", k)
		if _, ok := items[v.Name]; ok {
			klog.InfoS("pod existes, %s", v.Name)
			continue
		}
		result, err := podsClient.Create(context.TODO(), &v, metav1.CreateOptions{})
		if err != nil {
			klog.Error("pod create error, %s", err.Error())
			continue
		}
		jobCache[k] = *result
	}
	klog.Info("end checking ci and pods")
}

// cleanPod is used to clean pod on cron
//todo
func cleanPod() {

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
	return dtJob, nil
}
