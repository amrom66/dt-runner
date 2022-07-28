/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"

	"dt-runner/pkg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/robfig/cron/v3"
)

var kubeconfig string
var serverHost string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server is used to run dt-runner as a daemon server",
	Long:  `dt-runner will listen on a web port, which will be triggered by gitlab webhook.`,
	Run: func(cmd *cobra.Command, args []string) {

		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})

		http.HandleFunc("/", pkg.GitlabHook)
		// c is a schedule plan
		c := cron.New()
		c.AddFunc("* * * * *", pkg.StartPod)
		c.Start()
		go pkg.Watch(kubeconfig)

		if serverHost == "" {
			serverHost = pkg.GetLocalIpV4()
		} else {
			glog.Info("flag server.host is set, will use flag from command line")
			viper.Set("server.host", serverHost)
		}
		port := strings.Join([]string{":", strconv.Itoa(viper.GetInt("server.port"))}, "")

		glog.Infof("dt-runner is running on http://%s%s, with token: %s\n", serverHost, port, viper.GetString("webhook.token"))
		http.ListenAndServe(port, nil)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	home, err := os.UserHomeDir()
	if err != nil {
		glog.Errorf("get user home dir error: %v", err)
		os.Exit(1)
	}
	serverCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", home+"/.kube/config",
		"kubeconfig file(default is $HOME/.kube/config)")

	serverCmd.Flags().StringVar(&serverHost, "server.host", "", "server host")

	serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
