package pkg

import (
	"log"

	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"

	corev1 "k8s.io/api/core/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	appsv1 "dt-runner/api/apps/v1"

	crdclientset "dt-runner/generated/clientset/versioned"
	crdlister "dt-runner/generated/listers/apps/v1"

	informers "dt-runner/generated/informers/externalversions"

	crdscheme "dt-runner/generated/clientset/versioned/scheme"
)

type modelController struct {
	kubeclientset          kubernetes.Interface
	apiextensionsclientset apiextensionsclientset.Interface
	informer               cache.SharedIndexInformer
	crdclientset           crdclientset.Interface
	lister                 crdlister.ModelLister
	recorder               record.EventRecorder
	workqueue              workqueue.RateLimitingInterface
}

type ciController struct {
	kubeclientset          kubernetes.Interface
	apiextensionsclientset apiextensionsclientset.Interface
	informer               cache.SharedIndexInformer
	crdclientset           crdclientset.Interface
	lister                 crdlister.CiLister
	recorder               record.EventRecorder
	workqueue              workqueue.RateLimitingInterface
}

//
func newCiController(config *rest.Config) *ciController {
	klog.Infoln("Creating ci controller.")
	kubeClient := kubernetes.NewForConfigOrDie(config)
	apiextensionsClient := apiextensionsclientset.NewForConfigOrDie(config)
	ciClient := crdclientset.NewForConfigOrDie(config)
	informerFactory := informers.NewSharedInformerFactory(ciClient, 0)
	informer := informerFactory.Apps().V1().Cis()
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ci := obj.(*appsv1.Ci)
			klog.Infof("Added: %v", ci.Name)
			klog.Info("repo :%v", ci.Spec.Repo, " will be watched")
		},
		UpdateFunc: func(old, new interface{}) {
			ci := old.(*appsv1.Ci)
			klog.Infof("Updates: %v", ci.Name)

		},
		DeleteFunc: func(obj interface{}) {
			ci := obj.(*appsv1.Ci)
			klog.Infof("Deleted: %v", ci.Name)
			klog.Info("repo :%v", ci.Spec.Repo, " will not be watched")
		},
	})
	informerFactory.Start(wait.NeverStop)
	utilruntime.Must(appsv1.AddToScheme(crdscheme.Scheme))

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(crdscheme.Scheme, corev1.EventSource{Component: "ci-controller"})
	workqueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return &ciController{
		kubeclientset:          kubeClient,
		apiextensionsclientset: apiextensionsClient,
		crdclientset:           ciClient,
		informer:               informer.Informer(),
		lister:                 informer.Lister(),
		recorder:               recorder,
		workqueue:              workqueue,
	}
}

func newModelController(config *rest.Config) *modelController {
	klog.Infoln("Creating model controller.")
	kubeClient := kubernetes.NewForConfigOrDie(config)
	apiextensionsClient := apiextensionsclientset.NewForConfigOrDie(config)
	modelClient := crdclientset.NewForConfigOrDie(config)
	informerFactory := informers.NewSharedInformerFactory(modelClient, 0)
	informer := informerFactory.Apps().V1().Models()
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			model := obj.(*appsv1.Model)
			klog.Infof("Added: %v", model.Name)
		},
		UpdateFunc: func(old, new interface{}) {
			model := old.(*appsv1.Model)
			klog.Infof("Updates: %v", model.Name)
		},
		DeleteFunc: func(obj interface{}) {
			model := obj.(*appsv1.Model)
			klog.Infof("Deleted: %v", model.Name)
		},
	})
	informerFactory.Start(wait.NeverStop)
	utilruntime.Must(appsv1.AddToScheme(crdscheme.Scheme))

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(crdscheme.Scheme, corev1.EventSource{Component: "model-controller"})
	workqueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	return &modelController{
		kubeclientset:          kubeClient,
		apiextensionsclientset: apiextensionsClient,
		crdclientset:           modelClient,
		informer:               informer.Informer(),
		lister:                 informer.Lister(),
		recorder:               recorder,
		workqueue:              workqueue,
	}
}

func (ciController *ciController) run() {
	defer utilruntime.HandleCrash()
	defer ciController.workqueue.ShutDown()
	timeoutCh := make(chan struct{})
	if ok := cache.WaitForCacheSync(timeoutCh, ciController.informer.HasSynced); !ok {
		klog.Fatalln("Timeout expired during waiting for caches to sync.")
	}
	select {}
}

func (modelController *modelController) run() {
	defer utilruntime.HandleCrash()
	defer modelController.workqueue.ShutDown()
	timeoutCh := make(chan struct{})
	if ok := cache.WaitForCacheSync(timeoutCh, modelController.informer.HasSynced); !ok {
		klog.Fatalln("Timeout expired during waiting for caches to sync.")
	}
	select {}
}

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
	ciController := newCiController(config)
	modelController := newModelController(config)

	ciController.run()
	modelController.run()

	klog.Infoln("Starting custom controller.")

}
