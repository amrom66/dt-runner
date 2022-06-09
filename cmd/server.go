/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"dt-runner/api"
	"dt-runner/pkg"

	"github.com/go-playground/webhooks/v6/gitlab"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

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

		// init kubernetes client
		var config *rest.Config
		var err error
		if kubeconfig == "" {
			log.Printf("using in-cluster configuration")
			config, err = rest.InClusterConfig()
		} else {
			log.Printf("using configuration from '%s'", kubeconfig)
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		}
		if err != nil {
			panic(err)
		}
		// api.AddToScheme(scheme.Scheme)
		// crdConfig := *config
		// crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: api.GroupName, Version: api.GroupVersion}
		// crdConfig.APIPath = "/apis"
		// crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
		// crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

		// myClient, err := rest.UnversionedRESTClientFor(&crdConfig)
		// if err != nil {
		// 	panic(err)
		// }
		myClient, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		// pods, err := myClient.CoreV1().Pods("test-ns").List(context.TODO(), metav1.ListOptions{})
		// if err != nil {
		// 	panic(err)
		// }

		ci := api.Ci{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-job",
				Namespace: "default",
			},
			Spec: api.CiSpec{
				Repo:   "https://github.com/linjinbao666/vrmanager.git",
				Model:  "model-sample",
				Branch: "main",
			},
		}
		model := api.Model{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "model-sample",
				Namespace: "default",
			},
			Spec: api.ModelSpec{
				Tasks: []api.Task{
					{
						Name:  "build",
						Image: "maven:latest",
						Args:  []string{"mvn", "clean", "package"},
					},
				},
				Variables: map[string]string{
					"MAVEN_OPTS": "-DproxySet=true -DproxyHost=192.168.90.110 -DproxyPort=1087",
				},
			},
		}

		jobs, err := myClient.BatchV1().Jobs("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Panicln(err)
		}
		for _, job := range jobs.Items {
			if job.Name == strings.Join([]string{ci.Name, model.Name}, "-") {
				fmt.Println("job found")
				deletePolicy := metav1.DeletePropagationForeground
				err := myClient.BatchV1().Jobs("default").Delete(context.TODO(), job.Name, metav1.DeleteOptions{
					PropagationPolicy: &deletePolicy,
				})
				if err != nil {
					log.Panicln(err)
				}
				fmt.Println("job deleted")
			}
		}

		job, err := pkg.GenerateJob(myClient, ci, model)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Printf("job started %s\n", job.Name)

		http.ListenAndServe(port, nil)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panicln(err)
	}
	serverCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", home+"/.kube/config", "kubeconfig file(default is $HOME/.kube/config)")
	serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
