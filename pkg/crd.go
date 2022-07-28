package pkg

import (
	"context"
	"time"

	"github.com/golang/glog"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var config *rest.Config

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

// newCiController creates a new ciController
func newCiController(config *rest.Config) *ciController {
	glog.Infoln("Creating ci controller.")
	kubeClient := kubernetes.NewForConfigOrDie(config)
	apiextensionsClient := apiextensionsclientset.NewForConfigOrDie(config)
	ciClient := crdclientset.NewForConfigOrDie(config)
	informerFactory := informers.NewSharedInformerFactory(ciClient, 0)
	informer := informerFactory.Apps().V1().Cis()
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ci := obj.(*appsv1.Ci)
			glog.Infof("Added: %v", ci.Name)
			glog.Info("repo: ", ci.Spec.Repo, " will be watched")
		},
		UpdateFunc: func(old, new interface{}) {
			ci := old.(*appsv1.Ci)
			glog.Infof("Updates: %v", ci.Name)
		},
		DeleteFunc: func(obj interface{}) {
			ci := obj.(*appsv1.Ci)
			glog.Infof("Deleted: %v", ci.Name)
			glog.Info("repo :%v", ci.Spec.Repo, " will not be watched")
		},
	})
	informerFactory.Start(wait.NeverStop)
	utilruntime.Must(appsv1.AddToScheme(crdscheme.Scheme))

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
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
	glog.Infoln("Creating model controller.")
	kubeClient := kubernetes.NewForConfigOrDie(config)
	apiextensionsClient := apiextensionsclientset.NewForConfigOrDie(config)
	modelClient := crdclientset.NewForConfigOrDie(config)
	informerFactory := informers.NewSharedInformerFactory(modelClient, 0)
	informer := informerFactory.Apps().V1().Models()
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			model := obj.(*appsv1.Model)
			glog.Infof("Added: %v", model.Name)
		},
		UpdateFunc: func(old, new interface{}) {
			model := old.(*appsv1.Model)
			glog.Infof("Updates: %v", model.Name)
		},
		DeleteFunc: func(obj interface{}) {
			model := obj.(*appsv1.Model)
			glog.Infof("Deleted: %v", model.Name)
		},
	})
	informerFactory.Start(wait.NeverStop)
	utilruntime.Must(appsv1.AddToScheme(crdscheme.Scheme))

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
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
		glog.Fatalln("Timeout expired during waiting for caches to sync.")
	}
	select {}
}

func (modelController *modelController) run() {
	defer utilruntime.HandleCrash()
	defer modelController.workqueue.ShutDown()
	timeoutCh := make(chan struct{})
	if ok := cache.WaitForCacheSync(timeoutCh, modelController.informer.HasSynced); !ok {
		glog.Fatalln("Timeout expired during waiting for caches to sync.")
	}
	select {}
}

// Watch is used to start kubernetes client and watch crd resources
func Watch(kubeconfig string) {
	// init kubernetes client
	var err error
	if kubeconfig == "" {
		glog.Info("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		glog.Infof("using configuration from %s", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		panic(err)
	}

	ciController := newCiController(config)
	modelController := newModelController(config)

	ciController.run()
	modelController.run()

	glog.Infoln("Starting custom controller.")

}

func ListModels(namespace string) *appsv1.ModelList {
	modelclient := crdclientset.NewForConfigOrDie(config)
	modelList, err := modelclient.AppsV1().Models(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		glog.Info("list model error", err)
		return &appsv1.ModelList{}
	}
	return modelList
}

func ListCis(namespace string) *appsv1.CiList {
	ciclient := crdclientset.NewForConfigOrDie(config)
	cilist, err := ciclient.AppsV1().Cis(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		glog.Info("list cis error", err)
		return &appsv1.CiList{}
	}
	return cilist
}

func GetCi(namespace string, name string) *appsv1.Ci {
	ciclient := crdclientset.NewForConfigOrDie(config)
	ci, err := ciclient.AppsV1().Cis(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		glog.Info("get ci error", err)
		return &appsv1.Ci{}
	}
	return ci
}

func UpdateCi(namespace string, name string, podName string, status string) {
	ciclient := crdclientset.NewForConfigOrDie(config)
	ci, err := ciclient.AppsV1().Cis(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		glog.Info("get ci error", err)
	}
	ci.Status.Histroy = append(ci.Status.Histroy, appsv1.Histroy{
		CiName:  ci.Name,
		PodName: podName,
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Status:  status,
	})

	_, err = ciclient.AppsV1().Cis(namespace).UpdateStatus(context.TODO(), ci, metav1.UpdateOptions{})
	if err != nil {
		glog.Info("update ci error", err)
	}
	glog.Info("update ci finished: ", ci.Name)
}

func GetModel(namespace string, name string) *appsv1.Model {
	modelClient := crdclientset.NewForConfigOrDie(config)
	model, err := modelClient.AppsV1().Models(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		glog.Info("get model err", model)
		return &appsv1.Model{}
	}
	return model
}
