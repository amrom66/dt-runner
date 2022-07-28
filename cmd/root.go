/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// config for whole program
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dt-runner",
	Short: "runner for dtwave ci",
	Long: `
Dt runner is a runner for dtwave ci.
It will run for self-hosted gitlab service, 
Working as a centeralized compling service.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .dt-runner.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		pwd, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(pwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dt-runner")
	}
	viper.SetDefault("app.name", "dt-runner")
	viper.SetDefault("app.log", "/var/log/dt-runner")
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", 9001)
	viper.SetDefault("server.runtime", "kubernetes")
	if err := viper.ReadInConfig(); err == nil {
		logPath := viper.GetString("app.log")
		_, err := os.Stat(logPath)
		if err != nil {
			err := os.MkdirAll(logPath, 0766)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		defer glog.Flush()
		flag.Set("log_dir", logPath)
		flag.Set("alsologtostderr", "true")
		flag.Parse()
		glog.Infoln("Using config file:", viper.ConfigFileUsed())
	}
}
