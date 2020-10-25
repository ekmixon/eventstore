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
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"knative.dev/eventing/pkg/reconciler/source"
	"knative.dev/pkg/kmeta"

	"github.com/triggermesh/eventstore/pkg/apis/eventstores/v1alpha1"
	pkgreconciler "github.com/triggermesh/eventstore/pkg/reconciler"
	"github.com/triggermesh/eventstore/pkg/reconciler/resources"
)

const (
	adapterName = "inmemorystorage"

	envInMemoryDefaultGlobalTTL       = "EVENTSTORE_DEFAULT_GLOBAL_TTL"
	envInMemoryDefaultBridgeTTL       = "EVENTSTORE_DEFAULT_BRIDGE_TTL"
	envInMemoryDefaultInstanceTTL     = "EVENTSTORE_DEFAULT_INSTANCE_TTL"
	envInMemoryDefaultExpiredGCPeriod = "EVENTSTORE_DEFAULT_EXPIRED_GC_PERIOD"

	listenPort = 8080
)

// adapterConfig contains properties used to configure the eventstore's adapter.
// Public fields are automatically populated by envconfig.
type adapterConfig struct {
	// Configuration accessor for logging/metrics/tracing
	configs source.ConfigAccessor
	// Container image
	Image string `default:"gcr.io/triggermesh-private/eventstore-inmemory"`
}

// makeAdapterKnService returns a Knative Service object for the store's adapter.
func makeAdapterDeployment(o *v1alpha1.InMemoryStore, cfg *adapterConfig) *appsv1.Deployment {
	envApp := makeCommonAppEnv(o)

	deploymentLabels := pkgreconciler.MakeAdapterLabels(adapterName, o.Name)
	podLabels := pkgreconciler.MakeAdapterLabels(adapterName, o.Name)
	name := kmeta.ChildName(adapterName+"-", o.Name)
	envSvc := pkgreconciler.MakeServiceEnv(o.Name, o.Namespace)
	envObs := pkgreconciler.MakeObsEnv(cfg.configs)
	envs := append(envSvc, append(envApp, envObs...)...)

	return resources.NewDeployment(o.Namespace, name,
		// Deployment
		resources.Labels(deploymentLabels),
		resources.Selector(resources.AppNameLabel, adapterName),
		resources.Controller(o),

		// Pod
		resources.Image(cfg.Image),
		resources.Port("h2c", listenPort),
		resources.EnvVars(envs...),
		resources.PodLabels(podLabels),
	)
}

func makeCommonAppEnv(o *v1alpha1.InMemoryStore) []corev1.EnvVar {
	env := []corev1.EnvVar{}

	if o.Spec.DefaultGlobalTTL != nil {
		env = append(env, corev1.EnvVar{
			Name:  envInMemoryDefaultGlobalTTL,
			Value: strconv.Itoa(*o.Spec.DefaultGlobalTTL),
		})
	}

	if o.Spec.DefaultBridgeTTL != nil {
		env = append(env, corev1.EnvVar{
			Name:  envInMemoryDefaultBridgeTTL,
			Value: strconv.Itoa(*o.Spec.DefaultBridgeTTL),
		})
	}

	if o.Spec.DefaultInstanceTTL != nil {
		env = append(env, corev1.EnvVar{
			Name:  envInMemoryDefaultInstanceTTL,
			Value: strconv.Itoa(*o.Spec.DefaultInstanceTTL),
		})
	}

	if o.Spec.DefaultExpiredGCPeriod != nil {
		env = append(env, corev1.EnvVar{
			Name:  envInMemoryDefaultExpiredGCPeriod,
			Value: strconv.Itoa(*o.Spec.DefaultExpiredGCPeriod),
		})
	}

	return env
}
