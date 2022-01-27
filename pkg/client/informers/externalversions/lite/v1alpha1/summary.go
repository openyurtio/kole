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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	litev1alpha1 "github.com/openyurtio/kole/pkg/apis/lite/v1alpha1"
	versioned "github.com/openyurtio/kole/pkg/client/clientset/versioned"
	internalinterfaces "github.com/openyurtio/kole/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/openyurtio/kole/pkg/client/listers/lite/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SummaryInformer provides access to a shared informer and lister for
// Summaries.
type SummaryInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.SummaryLister
}

type summaryInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSummaryInformer constructs a new informer for Summary type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSummaryInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSummaryInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSummaryInformer constructs a new informer for Summary type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSummaryInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.LiteV1alpha1().Summaries(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.LiteV1alpha1().Summaries(namespace).Watch(context.TODO(), options)
			},
		},
		&litev1alpha1.Summary{},
		resyncPeriod,
		indexers,
	)
}

func (f *summaryInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSummaryInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *summaryInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&litev1alpha1.Summary{}, f.defaultInformer)
}

func (f *summaryInformer) Lister() v1alpha1.SummaryLister {
	return v1alpha1.NewSummaryLister(f.Informer().GetIndexer())
}