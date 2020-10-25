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

	pkgreconciler "knative.dev/pkg/reconciler"

	"github.com/triggermesh/eventstore/pkg/apis/eventstores/v1alpha1"
	reconcilerv1alpha1 "github.com/triggermesh/eventstore/pkg/generated/client/injection/reconciler/eventstores/v1alpha1/inmemorystore"
	libreconciler "github.com/triggermesh/eventstore/pkg/reconciler"
)

// Reconciler implements controller.Reconciler for the event store type.
type reconciler struct {
	// adapter properties
	adapterCfg *adapterConfig

	// Knative Service reconciler
	ksvcr libreconciler.KServiceReconciler
	dpr   libreconciler.DeploymentReconciler
}

// Check that our Reconciler implements Interface
var _ reconcilerv1alpha1.Interface = (*reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *reconciler) ReconcileKind(ctx context.Context, o *v1alpha1.InMemoryStore) pkgreconciler.Event {
	o.Status.InitializeConditions()
	o.Status.ObservedGeneration = o.Generation

	d, event := r.dpr.ReconcileDeployment(ctx, o, makeAdapterDeployment(o, r.adapterCfg))
	// adapter, event := r.ksvcr.ReconcileKService(ctx, o, makeAdapterKnService(o, r.adapterCfg))

	// deployment, event := r.

	o.Status.PropagateDeploymentAvailability(d)

	return event
}
