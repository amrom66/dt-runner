package pkg

import (
	"context"
	"fmt"
	"log"

	crdclient "dt-runner/generated/clientset/versioned"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Watch is used to start kubernetes client and watch crd resources
func Watch(kubeconfig string) {
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

	crdClient, err := crdclient.NewForConfig(config)
	if err != nil {
		log.Panicln(err)
	}
	cis, err := crdClient.AppsV1().Cis("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Panicln(err)
	}
	for ci := range cis.Items {
		fmt.Println(ci)
	}
}
