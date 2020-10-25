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

package inmemory

import (
	"context"

	"github.com/kelseyhightower/envconfig"

	"k8s.io/client-go/tools/cache"

	"knative.dev/eventing/pkg/reconciler/source"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
	serviceinformer "knative.dev/pkg/client/injection/kube/informers/core/v1/service"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	"github.com/triggermesh/eventstore/pkg/apis/eventstores/v1alpha1"
	informerv1alpha1 "github.com/triggermesh/eventstore/pkg/generated/client/injection/informers/eventstores/v1alpha1/inmemorystore"
	reconcilerv1alpha1 "github.com/triggermesh/eventstore/pkg/generated/client/injection/reconciler/eventstores/v1alpha1/inmemorystore"
	libreconciler "github.com/triggermesh/eventstore/pkg/reconciler"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	adapterCfg := &adapterConfig{
		configs: source.WatchConfigurations(ctx, adapterName, cmw, source.WithLogging, source.WithMetrics),
	}
	envconfig.MustProcess(adapterName, adapterCfg)

	storeInformer := informerv1alpha1.Get(ctx)
	deploymentInformer := deploymentinformer.Get(ctx)
	serviceInformer := serviceinformer.Get(ctx)

	r := &reconciler{
		dpr:  libreconciler.NewDeploymentReconciler(kubeclient.Get(ctx).AppsV1(), deploymentInformer.Lister()),
		svcr: libreconciler.NewServiceReconciler(kubeclient.Get(ctx).CoreV1(), serviceInformer.Lister()),

		adapterCfg: adapterCfg,
	}

	impl := reconcilerv1alpha1.NewImpl(ctx, r)
	logging.FromContext(ctx).Info("Setting up event handlers")

	storeInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	deploymentInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterControllerGVK((&v1alpha1.InMemoryStore{}).GetGroupVersionKind()),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	serviceInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterControllerGVK((&v1alpha1.InMemoryStore{}).GetGroupVersionKind()),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
