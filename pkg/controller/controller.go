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
	"time"

	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"

	"log"

	vegamcacheapi "github.com/sch00lb0y/vegamcache-operator/pkg/apis/vegamcache/v1alpha1"
	vegaminformer "github.com/sch00lb0y/vegamcache-operator/pkg/client/informers/externalversions"
	vegamlister "github.com/sch00lb0y/vegamcache-operator/pkg/client/listers/vegamcache/v1alpha1"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	listerappsv1 "k8s.io/client-go/listers/apps/v1"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
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
	podLister            listercorev1.PodLister
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
		vegamLister:          vegamInformer.Lister(),
		vegamHasSynced:       vegamInformer.Informer().HasSynced,
		deploymentLister:     deploymentInformer.Lister(),
		replicasetLister:     replicasetInformer.Lister(),
		podQueue:             workqueue.NewNamedRateLimitingQueue(workqueue.DefaultItemBasedRateLimiter(), "deployment-queue"),
		vegamQueue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultItemBasedRateLimiter(), "vegam-queue"),
		podLister:            podInformer.Lister(),
		podInformerHasSynced: podInformer.Informer().HasSynced,
	}
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			key, err := cache.MetaNamespaceKeyFunc(pod)
			if err != nil {
				log.Printf("unable to create key for obj %v : %v", pod, err)
				return
			}
			log.Println("new pod is added to the queue")
			ctrl.podQueue.Add(key)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldPod := old.(*v1.Pod)
			newPod := new.(*v1.Pod)
			if newPod.DeletionTimestamp != nil {
				key, err := cache.MetaNamespaceKeyFunc(newPod)
				if err != nil {
					log.Printf("unable to create key for obj %v : %v", newPod, err)
					return
				}
				log.Println("updated pod is added to the queue")
				ctrl.podQueue.Add(key)
				return
			}
			if oldPod.Status.Phase != v1.PodRunning && newPod.Status.Phase == v1.PodRunning {
				key, err := cache.MetaNamespaceKeyFunc(newPod)
				if err != nil {
					log.Printf("unable to create key for obj %v : %v", newPod, err)
					return
				}
				log.Println("updated pod is added to the queue")
				ctrl.podQueue.Add(key)
			}
		},
	})
	vegamcacheFactory.Vegamcacheoperator().V1alpha1().VegamCaches().
		Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newVegam := obj.(*vegamcacheapi.VegamCache)
			key, err := cache.MetaNamespaceKeyFunc(newVegam)
			if err != nil {
				log.Printf("unable to create key for obj %v : %v", newVegam, err)
				return
			}
			log.Println("vegam custom resource is added to the queue")
			ctrl.vegamQueue.Add(key)
		},
		DeleteFunc: func(obj interface{}) {
			newVegam := obj.(*vegamcacheapi.VegamCache)
			key, err := cache.MetaNamespaceKeyFunc(newVegam)
			if err != nil {
				log.Printf("unable to create key for obj %v : %v", newVegam, err)
				return
			}

			log.Println("deleted vegam custom resource is added to the queue")
			ctrl.vegamQueue.Add(key)
		},
	})
	return ctrl
}

func (c *vegamController) runPodWorker() {
	for c.processPod() {

	}
}

func (c *vegamController) processPod() bool {
	key, shutdown := c.podQueue.Get()
	if shutdown {
		log.Println("Shutting down")
		return false
	}
	defer c.podQueue.Done(key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key.(string))
	if err != nil {
		log.Printf("unable to split namespace and name for key %v", err)
		return true
	}
	pod, err := c.podLister.Pods(namespace).Get(name)
	if err != nil {
		log.Printf("error in retriving pod %v", err)
		return true
	}
	val, ok := pod.Labels["vegam"]
	if !ok {
		log.Printf("pod doesn't have vegam label so, ignoring this pod")
		return true
	}
	var podLabels labels.Set
	podLabels = map[string]string{"vegam": val}
	vegamCaches, err := c.vegamLister.List(podLabels.AsSelector())
	if len(vegamCaches) != 1 {
		log.Printf("label doens't not matching with any vegamcache, so ignoring this pod")
		return true
	}
	hitIp(pod.Status.PodIP)
	if err != nil {
		log.Printf("unable to list vegam caches from selectors %v", err)
	}
	return true
}

func (c *vegamController) Run(stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.vegamQueue.ShutDown()
	defer c.podQueue.ShutDown()
	log.Println("Vegam controller started")

	if !cache.WaitForCacheSync(stopCh, c.podInformerHasSynced, c.vegamHasSynced) {
		return fmt.Errorf("timeout on sync")
	}
	go wait.Until(c.runPodWorker, time.Second, stopCh)
	<-stopCh
	return nil
}
