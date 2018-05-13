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

package controller

import (
	"sync"

	"github.com/golang/glog"

	vegamcacheapi "github.com/sch00lb0y/vegamcache-operator/pkg/apis/vegamcache/v1alpha1"
	vegaminformer "github.com/sch00lb0y/vegamcache-operator/pkg/client/informers/externalversions"
	vegamlister "github.com/sch00lb0y/vegamcache-operator/pkg/client/listers/vegamcache/v1alpha1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	listerappsv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/kubernetes/pkg/controller"
)

type serverInfo struct {
	serverPort uint
	vegamPort  uint
}
type clusterInfo struct {
	info map[string]serverInfo // key pod uid
	sync.Mutex
}
type vegamController struct {
	vegamHasSynced       cache.InformerSynced
	podInformerHasSynced cache.InformerSynced
	vegamLister          vegamlister.VegamCacheLister
	deploymentLister     listerappsv1.DeploymentLister
	replicasetLister     listerappsv1.ReplicaSetLister
	clusterData          clusterInfo
	podQueue             workqueue.RateLimitingInterface
	vegamQueue           workqueue.RateLimitingInterface
}

func NewController(vegamcacheFactory vegaminformer.SharedInformerFactory, sharedInformer informers.SharedInformerFactory) *vegamController {
	vegamInformer := vegamcacheFactory.Vegamcacheoperator().V1alpha1().VegamCaches()
	deploymentInformer := sharedInformer.Apps().V1().Deployments()
	replicasetInformer := sharedInformer.Apps().V1().ReplicaSets()
	podInformer := sharedInformer.Core().V1().Pods()
	ctrl := &vegamController{
		vegamLister:      vegamInformer.Lister(),
		vegamHasSynced:   vegamInformer.Informer().HasSynced,
		deploymentLister: deploymentInformer.Lister(),
		replicasetLister: replicasetInformer.Lister(),
		podQueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "deployment-queue"),
		vegamQueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "vegam-queue"),
	}
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ctrl.podQueue.Add(obj)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldPod := old.(*v1.Pod)
			newPod := new.(*v1.Pod)
			if newPod.DeletionTimestamp != nil {
				key, err := controller.KeyFunc(newPod)
				if err != nil {
					glog.Errorf("unable to create key for obj %v : %v", newPod, err)
					return
				}
				ctrl.podQueue.Add(key)
				return
			}
			if oldPod.Status.Phase != v1.PodRunning && newPod.Status.Phase == v1.PodRunning {
				key, err := controller.KeyFunc(newPod)
				if err != nil {
					glog.Errorf("unable to create key for obj %v : %v", newPod, err)
					return
				}
				ctrl.podQueue.Add(key)
			}
		},
	})
	vegamcacheFactory.Vegamcacheoperator().V1alpha1().VegamCaches().
		Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newVegam := obj.(*vegamcacheapi.VegamCache)
			key, err := controller.KeyFunc(newVegam)
			if err != nil {
				glog.Errorf("unable to create key for obj %v : %v", newVegam, err)
				return
			}
			ctrl.vegamQueue.Add(key)
		},
		DeleteFunc: func(obj interface{}) {
			newVegam := obj.(*vegamcacheapi.VegamCache)
			key, err := controller.KeyFunc(newVegam)
			if err != nil {
				glog.Errorf("unable to create key for obj %v : %v", newVegam, err)
				return
			}
			ctrl.vegamQueue.Add(key)
		},
		UpdateFunc: func(_ interface{}, obj interface{}) {
			newVegam := obj.(*vegamcacheapi.VegamCache)
			key, err := controller.KeyFunc(newVegam)
			if err != nil {
				glog.Errorf("unable to create key for obj %v : %v", newVegam, err)
				return
			}
			ctrl.vegamQueue.Add(key)
		},
	})
	return ctrl
}

func (c *vegamController) Run(stop chan<- struct{}) {

}
