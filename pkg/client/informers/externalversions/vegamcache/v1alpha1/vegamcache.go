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
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	vegamcache_v1alpha1 "github.com/sch00lb0y/vegamcache-operator/pkg/apis/vegamcache/v1alpha1"
	versioned "github.com/sch00lb0y/vegamcache-operator/pkg/client/clientset/versioned"
	internalinterfaces "github.com/sch00lb0y/vegamcache-operator/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/sch00lb0y/vegamcache-operator/pkg/client/listers/vegamcache/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// VegamCacheInformer provides access to a shared informer and lister for
// VegamCaches.
type VegamCacheInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.VegamCacheLister
}

type vegamCacheInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewVegamCacheInformer constructs a new informer for VegamCache type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewVegamCacheInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredVegamCacheInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredVegamCacheInformer constructs a new informer for VegamCache type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredVegamCacheInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.VegamcacheoperatorV1alpha1().VegamCaches(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.VegamcacheoperatorV1alpha1().VegamCaches(namespace).Watch(options)
			},
		},
		&vegamcache_v1alpha1.VegamCache{},
		resyncPeriod,
		indexers,
	)
}

func (f *vegamCacheInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredVegamCacheInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *vegamCacheInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&vegamcache_v1alpha1.VegamCache{}, f.defaultInformer)
}

func (f *vegamCacheInformer) Lister() v1alpha1.VegamCacheLister {
	return v1alpha1.NewVegamCacheLister(f.Informer().GetIndexer())
}