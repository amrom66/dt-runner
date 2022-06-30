package pkg

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/webhooks/gitlab"
	"github.com/spf13/viper"
)

// GitlabHook is used to hook gitlab
func GitlabHook(w http.ResponseWriter, r *http.Request) {

	secret := viper.GetString("webhook.token")
	hook, _ := gitlab.New(gitlab.Options.Secret(secret))
	payload, err := hook.Parse(r, gitlab.PushEvents, gitlab.TagEvents, gitlab.SystemHookEvents)

	if err != nil {
		log.Println(err)
		return
	}
	dtjob, error := generateDtJob(payload)
	if error != nil {
		log.Println(error)
	}
	fmt.Println(dtjob)

}

// generateDtJob is used to generate dtjob struct
func generateDtJob(payload interface{}) (DtJob, error) {
	fmt.Println("generateDtJob")

	switch payload.(type) {
	case gitlab.PushEventPayload:
		push := payload.(gitlab.PushEventPayload)

	case gitlab.TagEventPayload:
		tag := payload.(gitlab.TagEventPayload)
	case gitlab.SystemHookPayload:

	default:
		fmt.Println("unknown event playload")
	}
}
