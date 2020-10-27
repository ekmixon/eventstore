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
	"strconv"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientgotesting "k8s.io/client-go/testing"

	"knative.dev/eventing/pkg/reconciler/source"
	fakekubeclient "knative.dev/pkg/client/injection/kube/client/fake"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	rt "knative.dev/pkg/reconciler/testing"

	"github.com/triggermesh/eventstore/pkg/apis/eventstores/v1alpha1"
	fakeinjectionclient "github.com/triggermesh/eventstore/pkg/generated/client/injection/client/fake"
	reconcilerv1alpha1 "github.com/triggermesh/eventstore/pkg/generated/client/injection/reconciler/eventstores/v1alpha1/inmemorystore"
	libreconciler "github.com/triggermesh/eventstore/pkg/reconciler"
	"github.com/triggermesh/eventstore/pkg/reconciler/resources"
	. "github.com/triggermesh/eventstore/pkg/reconciler/testing"
)

const (
	tNs   = "testns"
	tName = "test"
	tKey  = tNs + "/" + tName
	tUID  = types.UID("00000000-0000-0000-0000-000000000000")

	tImg = "registry/image:tag"
)

var (
	tGlobalTTL       = 1000
	tBridgeTTL       = 100
	tInstanceTTL     = 10
	tExpiredGCPeriod = 50
	tURI             = "dns:///my.service"
)

var tGenName = kmeta.ChildName(adapterName+"-", tName)

// Test the Reconcile method of the controller.Reconciler implemented by our controller.
//
// The environment for each test case is set up as follows:
//  1. MakeFactory initializes fake clients with the objects declared in the test case
//  2. MakeFactory injects those clients into a context along with fake event recorders, etc.
//  3. A Reconciler is constructed via a Ctor function using the values injected above
//  4. The Reconciler returned by MakeFactory is used to run the test case
func TestReconcile(t *testing.T) {
	testCases := rt.TableTest{
		// Creation

		{
			Name: "Store object creation",
			Key:  tKey,
			Objects: []runtime.Object{
				newStore(),
			},
			WantCreates: []runtime.Object{
				newDeployment(),
				newService(),
			},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newStore(
					withDeployment(newDeployment()),
					withService(newService()),
				),
			}},
			// WantEvents: []string{
			// 	createAdapterEvent(),
			// },
		},

		// Lifecycle

		// {
		// 	Name: "Adapter becomes Ready",
		// 	Key:  tKey,
		// 	Objects: []runtime.Object{
		// 		newStoreNotDeployed(),
		// 		newAdapterServiceReady(),
		// 	},
		// 	WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
		// 		Object: newStoreDeployed(),
		// 	}},
		// },
		// {
		// 	Name: "Adapter becomes NotReady",
		// 	Key:  tKey,
		// 	Objects: []runtime.Object{
		// 		newStoreDeployed(),
		// 		newAdapterServiceNotReady(),
		// 	},
		// 	WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
		// 		Object: newStoreNotDeployed(),
		// 	}},
		// },
		// {
		// 	Name: "Adapter is outdated",
		// 	Key:  tKey,
		// 	Objects: []runtime.Object{
		// 		newStoreDeployed(),
		// 		setAdapterImage(
		// 			newAdapterServiceReady(),
		// 			tImg+":old",
		// 		),
		// 	},
		// 	WantUpdates: []clientgotesting.UpdateActionImpl{{
		// 		Object: newAdapterServiceReady(),
		// 	}},
		// 	WantEvents: []string{
		// 		updateAdapterEvent(),
		// 	},
		// },

		// // Errors

		// {
		// 	Name: "Fail to create adapter service",
		// 	Key:  tKey,
		// 	WithReactors: []clientgotesting.ReactionFunc{
		// 		rt.InduceFailure("create", "services"),
		// 	},
		// 	Objects: []runtime.Object{
		// 		newStore(),
		// 	},
		// 	WantCreates: []runtime.Object{
		// 		newAdapterService(),
		// 	},
		// 	WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
		// 		Object: newStoreUnknownDeployed(false),
		// 	}},
		// 	WantEvents: []string{
		// 		failCreateAdapterEvent(),
		// 	},
		// 	WantErr: true,
		// },

		// {
		// 	Name: "Fail to update adapter service",
		// 	Key:  tKey,
		// 	WithReactors: []clientgotesting.ReactionFunc{
		// 		rt.InduceFailure("update", "services"),
		// 	},
		// 	Objects: []runtime.Object{
		// 		newStoreDeployed(),
		// 		setAdapterImage(
		// 			newAdapterServiceReady(),
		// 			tImg+":old",
		// 		),
		// 	},
		// 	WantUpdates: []clientgotesting.UpdateActionImpl{{
		// 		Object: newAdapterServiceReady(),
		// 	}},
		// 	WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
		// 		Object: newStoreUnknownDeployed(true),
		// 	}},
		// 	WantEvents: []string{
		// 		failUpdateAdapterEvent(),
		// 	},
		// 	WantErr: true,
		// },

		// // Edge cases

		// {
		// 	Name:    "Reconcile a non-existing object",
		// 	Key:     tKey,
		// 	Objects: nil,
		// 	WantErr: false,
		// },
	}

	testCases.Test(t, MakeFactory(reconcilerCtor))
}

// reconcilerCtor returns a Ctor for a HTTPTarget Reconciler.
var reconcilerCtor Ctor = func(t *testing.T, ctx context.Context, ls *Listers) controller.Reconciler {
	adapterCfg := &adapterConfig{
		Image:   tImg,
		configs: &source.EmptyVarsGenerator{},
	}

	r := &reconciler{
		adapterCfg: adapterCfg,
		dpr:        libreconciler.NewDeploymentReconciler(fakekubeclient.Get(ctx).AppsV1(), ls.GetDeploymentLister()),
		svcr:       libreconciler.NewServiceReconciler(fakekubeclient.Get(ctx).CoreV1(), ls.GetServiceLister()),
	}

	return reconcilerv1alpha1.NewReconciler(ctx, logging.FromContext(ctx),
		fakeinjectionclient.Get(ctx), ls.GetInMemoryStoreLister(),
		controller.GetEventRecorder(ctx), r)
}

/* Event targets */

type inMemoryStoreOptions func(*v1alpha1.InMemoryStore)

// newStore returns a test store object with pre-filled attributes.
func newStore(opts ...inMemoryStoreOptions) *v1alpha1.InMemoryStore {
	o := &v1alpha1.InMemoryStore{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: tNs,
			Name:      tName,
			UID:       tUID,
		},
		Spec: v1alpha1.InMemoryStoreSpec{
			DefaultGlobalTTL:       &tGlobalTTL,
			DefaultBridgeTTL:       &tBridgeTTL,
			DefaultInstanceTTL:     &tInstanceTTL,
			DefaultExpiredGCPeriod: &tExpiredGCPeriod,
		},
	}

	o.Status.InitializeConditions()

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func withDeployment(d *appsv1.Deployment) inMemoryStoreOptions {
	return func(o *v1alpha1.InMemoryStore) {
		o.Status.PropagateDeploymentAvailability(d)
	}
}

func withService(s *corev1.Service) inMemoryStoreOptions {
	return func(o *v1alpha1.InMemoryStore) {
		o.Status.PropagateServiceAvailability(s)
	}
}

type deploymentOptions func(*appsv1.Deployment)

func newDeployment(opts ...deploymentOptions) *appsv1.Deployment {
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: tNs,
			Name:      tGenName,
			Labels: labels.Set{
				resources.AppNameLabel:      adapterName,
				resources.AppInstanceLabel:  tName,
				resources.AppComponentLabel: resources.AdapterComponent,
				resources.AppPartOfLabel:    resources.PartOf,
				resources.AppManagedByLabel: resources.ManagedController,
			},
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(NewOwnerRefable(
					tName,
					(&v1alpha1.InMemoryStore{}).GetGroupVersionKind(),
					tUID,
				)),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					resources.AppNameLabel:     adapterName,
					resources.AppInstanceLabel: tName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels.Set{
						resources.AppNameLabel:      adapterName,
						resources.AppInstanceLabel:  tName,
						resources.AppComponentLabel: resources.AdapterComponent,
						resources.AppPartOfLabel:    resources.PartOf,
						resources.AppManagedByLabel: resources.ManagedController,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "adapter",
						Image: tImg,
						Ports: []corev1.ContainerPort{{
							Name:          "h2c",
							ContainerPort: 8080,
						}},
						Env: []corev1.EnvVar{
							{Name: "NAMESPACE", Value: tNs},
							{Name: "NAME", Value: tName},
							{Name: "EVENTSTORE_DEFAULT_GLOBAL_TTL", Value: strconv.Itoa(tGlobalTTL)},
							{Name: "EVENTSTORE_DEFAULT_BRIDGE_TTL", Value: strconv.Itoa(tBridgeTTL)},
							{Name: "EVENTSTORE_DEFAULT_INSTANCE_TTL", Value: strconv.Itoa(tInstanceTTL)},
							{Name: "EVENTSTORE_DEFAULT_EXPIRED_GC_PERIOD", Value: strconv.Itoa(tExpiredGCPeriod)},
							{Name: "K_LOGGING_CONFIG", Value: ""},
							{Name: "K_METRICS_CONFIG", Value: ""},
							{Name: "K_TRACING_CONFIG", Value: ""},
						},
					}},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}

	for _, opt := range opts {
		opt(d)
	}
	return d
}

type serviceOptions func(*corev1.Service)

func newService(opts ...serviceOptions) *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: tNs,
			Name:      tGenName,
			Labels: labels.Set{
				resources.AppNameLabel:      adapterName,
				resources.AppInstanceLabel:  tName,
				resources.AppComponentLabel: resources.AdapterComponent,
				resources.AppPartOfLabel:    resources.PartOf,
				resources.AppManagedByLabel: resources.ManagedController,
			},
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(NewOwnerRefable(
					tName,
					(&v1alpha1.InMemoryStore{}).GetGroupVersionKind(),
					tUID,
				)),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				resources.AppNameLabel:     adapterName,
				resources.AppInstanceLabel: tName,
			},
			Ports: []corev1.ServicePort{{
				Name:       "h2c",
				Port:       8080,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.IntOrString{IntVal: 8080},
			}},
		},
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Deployed: Unknown
// func newStoreUnknownDeployed(adapterExists bool) *v1alpha1.InMemoryStore {
// 	o := newStore()
// 	o.Status.PropagateDeploymentAvailability(nil)

// 	// cover the case where the URL was already set because an adapter was successfully created at an earlier time,
// 	// but the new adapter status can't be propagated, e.g. due to an update error
// 	if adapterExists {
// 		o.Status.Address = &v1alpha1.Addressable{
// 			URI: &tURI,
// 		}
// 	}

// 	return o
// }

// // Deployed: True
// func newStoreDeployed() *v1alpha1.InMemoryStore {
// 	o := newStore()
// 	o.Status.PropagateDeploymentAvailability(newAdapterServiceReady())
// 	return o
// }

// // Deployed: False
// func newStoreNotDeployed() *v1alpha1.InMemoryStore {
// 	o := newStore()
// 	o.Status.PropagateAvailability(newAdapterServiceNotReady())
// 	return o
// }

// /* Adapter service */

// // newAdapterService returns a test Service object with pre-filled attributes.
// func newAdapterService() *servingv1.Service {
// 	return &servingv1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Namespace: tNs,
// 			Name:      tGenName,
// 			Labels: labels.Set{
// 				resources.AppNameLabel:      adapterName,
// 				resources.AppInstanceLabel:  tName,
// 				resources.AppComponentLabel: resources.AdapterComponent,
// 				resources.AppPartOfLabel:    resources.PartOf,
// 				resources.AppManagedByLabel: resources.ManagedController,
// 				serving.VisibilityLabelKey:  serving.VisibilityClusterLocal,
// 			},
// 			OwnerReferences: []metav1.OwnerReference{
// 				*kmeta.NewControllerRef(NewOwnerRefable(
// 					tName,
// 					(&v1alpha1.HTTPTarget{}).GetGroupVersionKind(),
// 					tUID,
// 				)),
// 			},
// 		},
// 		Spec: servingv1.ServiceSpec{
// 			ConfigurationSpec: servingv1.ConfigurationSpec{
// 				Template: servingv1.RevisionTemplateSpec{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Labels: labels.Set{
// 							resources.AppNameLabel:      adapterName,
// 							resources.AppInstanceLabel:  tName,
// 							resources.AppComponentLabel: resources.AdapterComponent,
// 							resources.AppPartOfLabel:    resources.PartOf,
// 							resources.AppManagedByLabel: resources.ManagedController,
// 						},
// 					},
// 					Spec: servingv1.RevisionSpec{
// 						PodSpec: corev1.PodSpec{
// 							Containers: []corev1.Container{{
// 								Name:  resources.AdapterComponent,
// 								Image: tImg,
// 								Env: []corev1.EnvVar{
// 									{
// 										Name:  resources.EnvNamespace,
// 										Value: tNs,
// 									}, {
// 										Name:  resources.EnvName,
// 										Value: tName,
// 									}, {
// 										Name:  envHTTPEventType,
// 										Value: tResponseType,
// 									}, {
// 										Name:  envHTTPEventSource,
// 										Value: tResponseSource,
// 									}, {
// 										Name:  envHTTPURL,
// 										Value: tEndpointURL.String(),
// 									}, {
// 										Name:  envHTTPMethod,
// 										Value: tMethod,
// 									}, {
// 										Name:  envHTTPSkipVerify,
// 										Value: tSkipVerify,
// 									}, {
// 										Name:  envHTTPHeaders,
// 										Value: "key1:value1,key2:value2",
// 									}, {
// 										Name: source.EnvLoggingCfg,
// 									}, {
// 										Name: source.EnvMetricsCfg,
// 									}, {
// 										Name: source.EnvTracingCfg,
// 									},
// 								},
// 							}},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }

// // Ready: True
// func newAdapterServiceReady() *servingv1.Service {
// 	svc := newAdapterService()
// 	svc.Status.SetConditions(apis.Conditions{{
// 		Type:   v1alpha1.ConditionReady,
// 		Status: corev1.ConditionTrue,
// 	}})
// 	svc.Status.URL = &tAdapterURL
// 	return svc
// }

// // Ready: False
// func newAdapterServiceNotReady() *servingv1.Service {
// 	svc := newAdapterService()
// 	svc.Status.SetConditions(apis.Conditions{{
// 		Type:   v1alpha1.ConditionReady,
// 		Status: corev1.ConditionFalse,
// 	}})
// 	return svc
// }

// func setAdapterImage(o *servingv1.Service, img string) *servingv1.Service {
// 	o.Spec.Template.Spec.Containers[0].Image = img
// 	return o
// }

// /* Events */

// // TODO(antoineco): make event generators public inside pkg/reconciler for
// // easy reuse in tests

// func createAdapterEvent() string {
// 	return Eventf(corev1.EventTypeNormal, "KServiceCreated", "created kservice: \"%s/%s\"",
// 		tNs, tGenName)
// }

// func updateAdapterEvent() string {
// 	return Eventf(corev1.EventTypeNormal, "KServiceUpdated", "updated kservice: \"%s/%s\"",
// 		tNs, tGenName)
// }
// func failCreateAdapterEvent() string {
// 	return Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for create services")
// }

// func failUpdateAdapterEvent() string {
// 	return Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for update services")
// }
