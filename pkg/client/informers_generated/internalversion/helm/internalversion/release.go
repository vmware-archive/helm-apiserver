/*
Copyright 2017 The helm-apiserver Authors.

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

// This file was automatically generated by informer-gen

package internalversion

import (
	helm "github.com/bitnami-labs/helm-apiserver/pkg/apis/helm"
	internalclientset "github.com/bitnami-labs/helm-apiserver/pkg/client/clientset_generated/internalclientset"
	internalinterfaces "github.com/bitnami-labs/helm-apiserver/pkg/client/informers_generated/internalversion/internalinterfaces"
	internalversion "github.com/bitnami-labs/helm-apiserver/pkg/client/listers_generated/helm/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// ReleaseInformer provides access to a shared informer and lister for
// Releases.
type ReleaseInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.ReleaseLister
}

type releaseInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewReleaseInformer constructs a new informer for Release type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewReleaseInformer(client internalclientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredReleaseInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredReleaseInformer constructs a new informer for Release type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredReleaseInformer(client internalclientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Helm().Releases(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Helm().Releases(namespace).Watch(options)
			},
		},
		&helm.Release{},
		resyncPeriod,
		indexers,
	)
}

func (f *releaseInformer) defaultInformer(client internalclientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredReleaseInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *releaseInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&helm.Release{}, f.defaultInformer)
}

func (f *releaseInformer) Lister() internalversion.ReleaseLister {
	return internalversion.NewReleaseLister(f.Informer().GetIndexer())
}
