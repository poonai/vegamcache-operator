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

	vegaminformer "github.com/sch00lb0y/vegamcache-operator/pkg/client/informers/externalversions"
	vegamlister "github.com/sch00lb0y/vegamcache-operator/pkg/client/listers/vegamcache/v1alpha1"
	"k8s.io/client-go/informers"
	listerappsv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
)

type serverInfo struct {
	serverPort uint
	vegamPort  uint
}
type clusterInfo struct {
	info map[string]serverInfo
	sync.Mutex
}
type controller struct {
	vegamInformer       cache.SharedInformer
	vegamHasSynced      cache.InformerSynced
	deploymentInformer  cache.SharedInformer
	deploymentHasSynced cache.InformerSynced
	vegamLister         vegamlister.VegamCacheLister
	deploymentLister    listerappsv1.DeploymentLister
	clusterData         clusterInfo
}

func NewController(vegamcacheFactory vegaminformer.SharedInformerFactory, sharedInformer informers.SharedInformerFactory) *controller {
	vegamInformer := vegamcacheFactory.Vegamcacheoperator().V1alpha1().VegamCaches()
	deploymentInformer := sharedInformer.Apps().V1().Deployments()
	ctrl := &controller{
		vegamInformer:       vegamInformer.Informer(),
		vegamLister:         vegamInformer.Lister(),
		vegamHasSynced:      vegamInformer.Informer().HasSynced,
		deploymentInformer:  deploymentInformer.Informer(),
		deploymentLister:    deploymentInformer.Lister(),
		deploymentHasSynced: deploymentInformer.Informer().HasSynced,
	}
	return ctrl
}
