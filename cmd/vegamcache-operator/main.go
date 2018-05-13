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
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/hasura/gitkube/pkg/signals"
	vegamcacheapi "github.com/sch00lb0y/vegamcache-operator/pkg/apis/vegamcache/v1alpha1"
	vegamclient "github.com/sch00lb0y/vegamcache-operator/pkg/client/clientset/versioned"
	vegaminformer "github.com/sch00lb0y/vegamcache-operator/pkg/client/informers/externalversions"
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
	vegamClient, err := vegamclient.NewForConfig(config)
	if err != nil {
		glog.Fatalf("error on creating vegam client: %v", err)
	}
	vegamInformer := vegaminformer.NewSharedInformerFactory(vegamClient, time.Second*30)
	vegamInformer.Vegamcacheoperator().V1alpha1().VegamCaches().Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			obj, _ = obj.(*vegamcacheapi.VegamCache)
			fmt.Print(obj)
		},
	})
	vegamInformer.Start(stopCh)
}
