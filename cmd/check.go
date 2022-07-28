/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check config for dt-runner",
	Long: `check config for dt-runner:
1. leader election
2. check gitlab service
3. check runc service
4. check workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkConfig(); err != nil {
			glog.Errorln(err)
			os.Exit(1)
		}
		glog.Infoln()
		glog.Infoln("all config checked success, now you can start dt-runer by `dt-runner server`")
		glog.Infoln()
	},
}

func checkConfig() error {
	sysType := runtime.GOOS
	if sysType == "windows" {
		log.Println("Target OS is linux")
		os.Exit(1)
	}
	name := viper.GetString("app.name")
	log.Println("app name is set to:", name)
	token := viper.GetString("webhook.token")
	if token == "" {
		return fmt.Errorf("token is not set")
	}

	runtime := viper.GetString("server.runtime")
	if runtime != "runc" && runtime != "kubernetes" {
		fmt.Println("runtime is not supported, will rollback to default kubernetes")
		viper.Set("server.ruuntime", "kubernetes")
	}
	log.Println("runtime is set to:", runtime)

	port := viper.GetInt32("server.port")
	if port == 0 || port > 65535 {
		log.Println("port is not set or invalid, will rollback to default 9001")
		viper.Set("server.port", 9001)
	}
	log.Println("port is set to:", port)

	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
