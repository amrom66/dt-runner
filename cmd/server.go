/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/webhooks/v6/gitlab"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server is used to run dt-runner as a daemon server",
	Long:  `dt-runner will listen on a web port, which will be triggered by gitlab webhook.`,
	Run: func(cmd *cobra.Command, args []string) {

		secret := viper.GetString("webhook.token")
		hook, _ := gitlab.New(gitlab.Options.Secret(secret))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		})
		port := strings.Join([]string{":", strconv.Itoa(viper.GetInt("server.port"))}, "")
		fmt.Printf("dt-runner is running on port:%s, with token:%s\n", port, secret)
		http.ListenAndServe(port, nil)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
