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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/sch00lb0y/vegamcache-operator/pkg/apis/vegamcache/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// VegamCacheLister helps list VegamCaches.
type VegamCacheLister interface {
	// List lists all VegamCaches in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.VegamCache, err error)
	// VegamCaches returns an object that can list and get VegamCaches.
	VegamCaches(namespace string) VegamCacheNamespaceLister
	VegamCacheListerExpansion
}

// vegamCacheLister implements the VegamCacheLister interface.
type vegamCacheLister struct {
	indexer cache.Indexer
}

// NewVegamCacheLister returns a new VegamCacheLister.
func NewVegamCacheLister(indexer cache.Indexer) VegamCacheLister {
	return &vegamCacheLister{indexer: indexer}
}

// List lists all VegamCaches in the indexer.
func (s *vegamCacheLister) List(selector labels.Selector) (ret []*v1alpha1.VegamCache, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.VegamCache))
	})
	return ret, err
}

// VegamCaches returns an object that can list and get VegamCaches.
func (s *vegamCacheLister) VegamCaches(namespace string) VegamCacheNamespaceLister {
	return vegamCacheNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// VegamCacheNamespaceLister helps list and get VegamCaches.
type VegamCacheNamespaceLister interface {
	// List lists all VegamCaches in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.VegamCache, err error)
	// Get retrieves the VegamCache from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.VegamCache, error)
	VegamCacheNamespaceListerExpansion
}

// vegamCacheNamespaceLister implements the VegamCacheNamespaceLister
// interface.
type vegamCacheNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all VegamCaches in the indexer for a given namespace.
func (s vegamCacheNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.VegamCache, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.VegamCache))
	})
	return ret, err
}

// Get retrieves the VegamCache from the indexer for a given namespace and name.
func (s vegamCacheNamespaceLister) Get(name string) (*v1alpha1.VegamCache, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("vegamcache"), name)
	}
	return obj.(*v1alpha1.VegamCache), nil
}