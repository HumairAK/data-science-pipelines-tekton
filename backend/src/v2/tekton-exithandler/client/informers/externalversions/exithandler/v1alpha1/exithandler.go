// Copyright 2023 kubeflow.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	exithandlerv1alpha1 "github.com/kubeflow/pipelines/backend/src/v2/tekton-exithandler/apis/exithandler/v1alpha1"
	versioned "github.com/kubeflow/pipelines/backend/src/v2/tekton-exithandler/client/clientset/versioned"
	internalinterfaces "github.com/kubeflow/pipelines/backend/src/v2/tekton-exithandler/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kubeflow/pipelines/backend/src/v2/tekton-exithandler/client/listers/exithandler/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ExitHandlerInformer provides access to a shared informer and lister for
// ExitHandlers.
type ExitHandlerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ExitHandlerLister
}

type exitHandlerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewExitHandlerInformer constructs a new informer for ExitHandler type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewExitHandlerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredExitHandlerInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredExitHandlerInformer constructs a new informer for ExitHandler type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredExitHandlerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CustomV1alpha1().ExitHandlers(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CustomV1alpha1().ExitHandlers(namespace).Watch(context.TODO(), options)
			},
		},
		&exithandlerv1alpha1.ExitHandler{},
		resyncPeriod,
		indexers,
	)
}

func (f *exitHandlerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredExitHandlerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *exitHandlerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&exithandlerv1alpha1.ExitHandler{}, f.defaultInformer)
}

func (f *exitHandlerInformer) Lister() v1alpha1.ExitHandlerLister {
	return v1alpha1.NewExitHandlerLister(f.Informer().GetIndexer())
}
