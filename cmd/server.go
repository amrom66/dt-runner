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
			playload, err := hook.Parse(r, gitlab.PushEvents, gitlab.TagEvents)
			if err != nil {
				log.Println(err)
				return
			}
			switch playload {
			case gitlab.PushEvents:
				fmt.Println("push event")
			case gitlab.TagEvents:
				fmt.Println("tag event")
			default:
				fmt.Println("unknown event")
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
