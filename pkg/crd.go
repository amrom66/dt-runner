package pkg

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	crdClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	crdClient.CCo
}
