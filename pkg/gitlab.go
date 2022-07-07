package pkg

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/webhooks/gitlab"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var jobCache = make(map[string]corev1.Pod)

// GitlabHook is used to hook gitlab
func GitlabHook(w http.ResponseWriter, r *http.Request) {

	secret := viper.GetString("webhook.token")
	hook, _ := gitlab.New(gitlab.Options.Secret(secret))
	payload, err := hook.Parse(r, gitlab.PushEvents, gitlab.TagEvents, gitlab.SystemHookEvents)

	if err != nil {
		log.Println(err)
		return
	}
	dtJob, error := generateDtJob(payload)
	if error != nil {
		log.Println(error)
	}
	fmt.Printf("dtjob generate, %s", dtJob.name)

	pod, err := GeneratePod(dtJob)
	if err != nil {
		fmt.Println(err)
	}
	jobCache[dtJob.name] = pod
}

// startPod is used to start a pod if confirmed
// startPod will keep check key:value in podCache
func StartPod() {
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods(DefaultNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	fmt.Println("jobCache", len(jobCache))
	for k, v := range jobCache {
		fmt.Printf(k, v.Name)
	}
}

// generateDtJob is used to generate DtJob
func generateDtJob(payload interface{}) (dtJob DtJob, err error) {
	dtJob.name = RandomString(6)
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
