/*
Copyright 2018 The vegamcache-operator Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"time"

	"log"

	vegamclient "github.com/sch00lb0y/vegamcache-operator/pkg/client/clientset/versioned"
	vegaminformer "github.com/sch00lb0y/vegamcache-operator/pkg/client/informers/externalversions"
	vegamcontroller "github.com/sch00lb0y/vegamcache-operator/pkg/controller"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/sample-controller/pkg/signals"
)

//const configFile = "/home/schoolgirl/.kube/config"

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.Parse()
	stopCh := signals.SetupSignalHandler()
	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		log.Fatalf("error on creating config from file: %v", err)
	}
	vegamClient, err := vegamclient.NewForConfig(config)
	if err != nil {
		log.Fatalf("error on creating vegam client: %v", err)
	}
	vegamInformer := vegaminformer.NewSharedInformerFactory(vegamClient, time.Second*30)

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("error on creating kuberentes client: %v", err)
	}
	sharedInformer := informers.NewSharedInformerFactory(kubeClient, time.Second*30)
	vegamController := vegamcontroller.NewController(vegamInformer, sharedInformer)
	go sharedInformer.Core().V1().Pods().Informer().Run(stopCh)
	go vegamInformer.Vegamcacheoperator().V1alpha1().VegamCaches().Informer().Run(stopCh)
	// let them sync
	time.Sleep(5)
	if err := vegamController.Run(stopCh); err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
