package main

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/hasura/gitkube/pkg/signals"
	vegamclient "github.com/sch00lb0y/vegamcache-operator/pkg/client/clientset/versioned"
	vegaminformer "github.com/sch00lb0y/vegamcache-operator/pkg/client/informers/externalversions"
	"k8s.io/apimachinery/pkg/labels"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

const configFile = "/home/schoolgirl/.kube/config"

func main() {
	stopCh := signals.SetupSignalHandler()
	config, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		glog.Fatalf("error on creating config from file: %v", err)
	}
	if err != nil {
		glog.Fatalf("error on creating kuberentes client: %v", err)
	}
	vegamClient, err := vegamclient.NewForConfig(config)
	vegamInformer := vegaminformer.NewSharedInformerFactory(vegamClient, time.Second*30)
	var label labels.Set
	label = map[string]string{"vegam": "app1"}
	go vegamInformer.Vegamcacheoperator().V1alpha1().VegamCaches().Informer().Run(stopCh)
	//	go vegamInformer.Start(stopCh)
	tic := time.Tick(time.Second * 5)
	for range tic {
		c, err := vegamInformer.Vegamcacheoperator().V1alpha1().VegamCaches().Lister().List(label.AsSelector())
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print(c[0])
	}
	if !cache.WaitForCacheSync(stopCh, vegamInformer.Vegamcacheoperator().V1alpha1().VegamCaches().Informer().HasSynced) {
		fmt.Errorf("timeout on sync")
	}

	// kubeClient, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	glog.Fatalf("error on creating kuberentes client: %v", err)
	// }
	// sharedInformer := informers.NewSharedInformerFactory(kubeClient, time.Second*30)
	// sharedInformer.Start(stopCh)
	// sync := sharedInformer.Core().V1().Pods().Informer().HasSynced
	// sharedInformer.Core().V1().Pods().Informer().Run(stopCh)
	// cache.WaitForCacheSync(stopCh, sync)
	// c, _ := sharedInformer.Core().V1().Pods().Lister().List(labels.Everything())
	// fmt.Print(c)
	<-stopCh
}
