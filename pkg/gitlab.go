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
	switch payload.(type) {
	case gitlab.PushEventPayload:
		fmt.Println("push event playload")
		push := payload.(gitlab.PushEventPayload)
		fmt.Printf("%+v", push)
	case gitlab.TagEventPayload:
		fmt.Println("tag event playload")
		tag := payload.(gitlab.TagEventPayload)
		fmt.Printf("%+v", tag)
	case gitlab.SystemHookPayload:
		fmt.Println("system event playload")
	default:
		fmt.Println("unknown event playload")
	}
}
