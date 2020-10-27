/*
Copyright (c) 2020 TriggerMesh Inc.

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

package testing

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakek8sclient "k8s.io/client-go/kubernetes/fake"
	applistersv1 "k8s.io/client-go/listers/apps/v1"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	rt "knative.dev/pkg/reconciler/testing"
	fakeservingclient "knative.dev/serving/pkg/client/clientset/versioned/fake"

	storesv1alpha1 "github.com/triggermesh/eventstore/pkg/apis/eventstores/v1alpha1"
	fakestoresclient "github.com/triggermesh/eventstore/pkg/generated/client/clientset/internalclientset/fake"
	storeslisters "github.com/triggermesh/eventstore/pkg/generated/client/listers/eventstores/v1alpha1"
)

var clientSetSchemes = []func(*runtime.Scheme) error{
	fakestoresclient.AddToScheme,
	fakek8sclient.AddToScheme,
	fakeservingclient.AddToScheme,
}

// NewScheme returns a new scheme populated with the types defined in clientSetSchemes.
func NewScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()

	sb := runtime.NewSchemeBuilder(clientSetSchemes...)
	if err := sb.AddToScheme(scheme); err != nil {
		panic(fmt.Errorf("error building Scheme: %s", err))
	}

	return scheme
}

// Listers returns listers and objects filtered from those listers.
type Listers struct {
	sorter rt.ObjectSorter
}

// NewListers returns a new instance of Listers initialized with the given objects.
func NewListers(scheme *runtime.Scheme, objs []runtime.Object) Listers {
	ls := Listers{
		sorter: rt.NewObjectSorter(scheme),
	}

	ls.sorter.AddObjects(objs...)

	return ls
}

// IndexerFor returns the indexer for the given object.
func (l *Listers) IndexerFor(obj runtime.Object) cache.Indexer {
	return l.sorter.IndexerForObjectType(obj)
}

// GetInMemoryStoreObjects returns objects from the stores API.
func (l *Listers) GetInMemoryStoreObjects() []runtime.Object {
	return l.sorter.ObjectsForSchemeFunc(fakestoresclient.AddToScheme)
}

// GetKubeObjects returns objects from Kubernetes APIs.
func (l *Listers) GetKubeObjects() []runtime.Object {
	return l.sorter.ObjectsForSchemeFunc(fakek8sclient.AddToScheme)
}

// GetServingObjects returns objects from the serving API.
func (l *Listers) GetServingObjects() []runtime.Object {
	return l.sorter.ObjectsForSchemeFunc(fakeservingclient.AddToScheme)
}

// GetInMemoryStoreLister returns a Lister for InMemoryStore objects.
func (l *Listers) GetInMemoryStoreLister() storeslisters.InMemoryStoreLister {
	return storeslisters.NewInMemoryStoreLister(l.IndexerFor(&storesv1alpha1.InMemoryStore{}))
}

// GetDeploymentLister returns a lister for Deployment objects.
func (l *Listers) GetDeploymentLister() applistersv1.DeploymentLister {
	return applistersv1.NewDeploymentLister(l.IndexerFor(&appsv1.Deployment{}))
}

// GetServiceLister returns a lister for Service objects.
func (l *Listers) GetServiceLister() corelistersv1.ServiceLister {
	return corelistersv1.NewServiceLister(l.IndexerFor(&corev1.Service{}))
}
