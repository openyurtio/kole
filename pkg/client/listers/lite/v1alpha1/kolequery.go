/*
Copyright 2022 The OpenYurt Authors.

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
	v1alpha1 "github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// KoleQueryLister helps list KoleQueries.
// All objects returned here must be treated as read-only.
type KoleQueryLister interface {
	// List lists all KoleQueries in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.KoleQuery, err error)
	// KoleQueries returns an object that can list and get KoleQueries.
	KoleQueries(namespace string) KoleQueryNamespaceLister
	KoleQueryListerExpansion
}

// koleQueryLister implements the KoleQueryLister interface.
type koleQueryLister struct {
	indexer cache.Indexer
}

// NewKoleQueryLister returns a new KoleQueryLister.
func NewKoleQueryLister(indexer cache.Indexer) KoleQueryLister {
	return &koleQueryLister{indexer: indexer}
}

// List lists all KoleQueries in the indexer.
func (s *koleQueryLister) List(selector labels.Selector) (ret []*v1alpha1.KoleQuery, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.KoleQuery))
	})
	return ret, err
}

// KoleQueries returns an object that can list and get KoleQueries.
func (s *koleQueryLister) KoleQueries(namespace string) KoleQueryNamespaceLister {
	return koleQueryNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// KoleQueryNamespaceLister helps list and get KoleQueries.
// All objects returned here must be treated as read-only.
type KoleQueryNamespaceLister interface {
	// List lists all KoleQueries in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.KoleQuery, err error)
	// Get retrieves the KoleQuery from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.KoleQuery, error)
	KoleQueryNamespaceListerExpansion
}

// koleQueryNamespaceLister implements the KoleQueryNamespaceLister
// interface.
type koleQueryNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all KoleQueries in the indexer for a given namespace.
func (s koleQueryNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.KoleQuery, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.KoleQuery))
	})
	return ret, err
}

// Get retrieves the KoleQuery from the indexer for a given namespace and name.
func (s koleQueryNamespaceLister) Get(name string) (*v1alpha1.KoleQuery, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("kolequery"), name)
	}
	return obj.(*v1alpha1.KoleQuery), nil
}