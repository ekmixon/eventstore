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
	knapis "knative.dev/pkg/apis"
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

	deploymentNotReady = "The deployment is not ready"
	deploymentNotFound = "The deployment can not be found"
	serviceNotfound    = "The service can not be found"
)

var (
	tGenName = kmeta.ChildName(adapterName+"-", tName)

	tGlobalTTL       = 1000
	tBridgeTTL       = 100
	tInstanceTTL     = 10
	tExpiredGCPeriod = 50
	tURI             = "dns:///" + tGenName + "." + tNs + ":8080"
)

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
					withDeploymentStatus(corev1.ConditionFalse, v1alpha1.ReasonUnavailable, deploymentNotReady),
					withServiceStatus(corev1.ConditionTrue, "", ""),
					withReadyStatus(corev1.ConditionFalse, v1alpha1.ReasonUnavailable, deploymentNotReady),
					withAddress(tURI),
				),
			}},
		},

		// Lifecycle

		{
			Name: "Deployment and Service Ready",
			Key:  tKey,
			Objects: []runtime.Object{
				newStore(),
				newDeployment(withDeploymentAvailable),
				newService(),
			},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newStore(
					withDeploymentStatus(corev1.ConditionTrue, "", ""),
					withServiceStatus(corev1.ConditionTrue, "", ""),
					withReadyStatus(corev1.ConditionTrue, "", ""),
					withAddress(tURI),
				),
			}},
		},
		{
			Name: "Deployment becomes NotReady",
			Key:  tKey,
			Objects: []runtime.Object{
				newStore(
					withDeploymentStatus(corev1.ConditionTrue, "", ""),
					withServiceStatus(corev1.ConditionTrue, "", ""),
					withReadyStatus(corev1.ConditionTrue, "", ""),
				),
				newDeployment(withDeploymentUnavailable),
				newService(),
			},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newStore(
					withDeploymentStatus(corev1.ConditionFalse, v1alpha1.ReasonUnavailable, deploymentNotReady),
					withServiceStatus(corev1.ConditionTrue, "", ""),
					withReadyStatus(corev1.ConditionFalse, v1alpha1.ReasonUnavailable, deploymentNotReady),
					withAddress(tURI),
				),
			}},
		},
		{
			Name: "Deployment is outdated",
			Key:  tKey,
			Objects: []runtime.Object{
				newStore(
					withDeploymentStatus(corev1.ConditionTrue, "", ""),
					withServiceStatus(corev1.ConditionTrue, "", ""),
					withReadyStatus(corev1.ConditionTrue, "", ""),
				),
				newDeployment(withDeploymentAvailable, withDeploymentImage("old-image")),
				newService(),
			},
			WantUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newDeployment(withDeploymentAvailable),
			}},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newStore(
					withDeploymentStatus(corev1.ConditionTrue, "", ""),
					withServiceStatus(corev1.ConditionTrue, "", ""),
					withReadyStatus(corev1.ConditionTrue, "", ""),
					withAddress(tURI),
				),
			}},
		},

		// Errors

		{
			Name: "Fail to create service",
			Key:  tKey,
			WithReactors: []clientgotesting.ReactionFunc{
				rt.InduceFailure("create", "services"),
			},
			Objects: []runtime.Object{
				newStore(),
			},
			WantCreates: []runtime.Object{
				newDeployment(),
				newService(),
			},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newStore(
					withDeploymentStatus(corev1.ConditionFalse, v1alpha1.ReasonUnavailable, deploymentNotReady),
					withServiceStatus(corev1.ConditionUnknown, v1alpha1.ReasonNotFound, serviceNotfound),
					withReadyStatus(corev1.ConditionFalse, v1alpha1.ReasonUnavailable, deploymentNotReady),
				),
			}},
			WantEvents: []string{
				failCreateServiceEvent(),
			},
			WantErr: true,
		},

		{
			Name: "Fail to update deployment",
			Key:  tKey,
			WithReactors: []clientgotesting.ReactionFunc{
				rt.InduceFailure("update", "deployments"),
			},
			Objects: []runtime.Object{
				newStore(
					withDeploymentStatus(corev1.ConditionTrue, "", ""),
					withServiceStatus(corev1.ConditionTrue, "", ""),
					withReadyStatus(corev1.ConditionTrue, "", ""),
				),
				newDeployment(withDeploymentAvailable, withDeploymentImage("old-image")),
				newService(),
			},
			WantUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newDeployment(withDeploymentAvailable),
			}},
			WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
				Object: newStore(
					withDeploymentStatus(corev1.ConditionUnknown, v1alpha1.ReasonNotFound, deploymentNotFound),
					withServiceStatus(corev1.ConditionTrue, "", ""),
					withReadyStatus(corev1.ConditionUnknown, v1alpha1.ReasonNotFound, deploymentNotFound),
				),
			}},
			WantEvents: []string{
				failUpdateDeploymentEvent(),
			},
			WantErr: true,
		},

		// Edge cases

		{
			Name:    "Reconcile a non-existing object",
			Key:     tKey,
			Objects: nil,
			WantErr: false,
		},
	}

	testCases.Test(t, MakeFactory(reconcilerCtor))
}

// reconcilerCtor returns a Ctor for a HTTPTarget Reconciler.
var reconcilerCtor Ctor = func(t *testing.T, ctx context.Context, ls *Listers) controller.Reconciler {
	adapterCfg := &adapterConfig{
		Image:   tImg,
		configs: &source.EmptyVarsGenerator{},
	}

	r := &Reconciler{
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

func withAddress(address string) inMemoryStoreOptions {
	return func(o *v1alpha1.InMemoryStore) {
		o.Status.Address = &v1alpha1.Addressable{URI: &address}
	}
}

func withDeploymentStatus(status corev1.ConditionStatus, reason, message string) inMemoryStoreOptions {
	return withStatus(v1alpha1.ConditionDeploymentReady, status, reason, message)
}

func withReadyStatus(status corev1.ConditionStatus, reason, message string) inMemoryStoreOptions {
	return withStatus(v1alpha1.ConditionReady, status, reason, message)
}

func withServiceStatus(status corev1.ConditionStatus, reason, message string) inMemoryStoreOptions {
	return withStatus(v1alpha1.ConditionServiceReady, status, reason, message)
}

func withStatus(cType knapis.ConditionType, status corev1.ConditionStatus, reason, message string) inMemoryStoreOptions {
	return func(o *v1alpha1.InMemoryStore) {
		var c *knapis.Condition

		conds := o.Status.Conditions
		for i := range conds {
			if conds[i].Type == cType {
				c = &conds[i]
				break
			}
		}

		if c == nil {
			c = &knapis.Condition{Type: cType}
			o.Status.Conditions = append(conds, *c)
		}

		c.Reason = reason
		c.Message = message
		c.Status = status
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

func withDeploymentAvailable(d *appsv1.Deployment) {
	d.Status = appsv1.DeploymentStatus{
		Conditions: []appsv1.DeploymentCondition{{
			Type:   appsv1.DeploymentAvailable,
			Status: "True",
		}},
	}
}

func withDeploymentUnavailable(d *appsv1.Deployment) {
	d.Status = appsv1.DeploymentStatus{
		Conditions: []appsv1.DeploymentCondition{{
			Type:   appsv1.DeploymentAvailable,
			Status: "False",
		}},
	}
}

func withDeploymentImage(image string) deploymentOptions {
	return func(d *appsv1.Deployment) {
		d.Spec.Template.Spec.Containers[0].Image = image
	}
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

func failUpdateDeploymentEvent() string {
	return Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for update deployments")
}

func failCreateServiceEvent() string {
	return Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for create services")
}
