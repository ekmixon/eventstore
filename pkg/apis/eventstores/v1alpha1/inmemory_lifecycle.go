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

package v1alpha1

import (
	"errors"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"knative.dev/eventing/pkg/apis/duck"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// GetGroupVersionKind implements kmeta.OwnerRefable.
func (s *InMemoryStore) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("InMemoryStore")
}

var inmemoryCondSet = apis.NewLivingConditionSet(
	ConditionDeploymentReady,
	ConditionServiceReady,
)

// GetConditionSet implements duckv1.KRShaped.
func (s *InMemoryStore) GetConditionSet() apis.ConditionSet {
	return inmemoryCondSet
}

// GetStatus implements duckv1.KRShaped.
func (s *InMemoryStore) GetStatus() *duckv1.Status {
	return &s.Status.Status
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *InMemoryStoreStatus) InitializeConditions() {
	inmemoryCondSet.Manage(s).InitializeConditions()
}

// PropagateDeploymentAvailability uses the readiness of the provided deployment to
// determine whether the Ready condition should be marked as true or false.
func (s *InMemoryStoreStatus) PropagateDeploymentAvailability(d *appsv1.Deployment) {
	if d == nil {
		inmemoryCondSet.Manage(s).MarkUnknown(ConditionDeploymentReady, ReasonNotFound,
			"The deployment can not be found")
		return
	}

	if duck.DeploymentIsAvailable(&d.Status, false) {
		inmemoryCondSet.Manage(s).MarkTrue(ConditionDeploymentReady)
		return
	}

	inmemoryCondSet.Manage(s).MarkFalse(ConditionDeploymentReady, ReasonUnavailable,
		"The deployment is not ready")
}

// PropagateServiceAvailability uses the presence of the provided service to
// determine whether the Ready condition should be marked as true or false.
func (s *InMemoryStoreStatus) PropagateServiceAvailability(svc *corev1.Service) {
	if svc == nil {
		inmemoryCondSet.Manage(s).MarkUnknown(ConditionServiceReady, ReasonNotFound,
			"The service can not be found")
		return
	}

	if svc.Spec.ClusterIP == "" {
		inmemoryCondSet.Manage(s).MarkFalse(ConditionServiceReady, ReasonUnavailable,
			"The service has no Cluster IP assigned")
	}

	if s.Address == nil {
		s.Address = &Addressable{}
	}

	url, err := serviceToAddress(svc)
	if err != nil {
		inmemoryCondSet.Manage(s).MarkFalse(ConditionServiceReady, ReasonUnavailable,
			err.Error())
	}

	s.Address.URI = &url
	inmemoryCondSet.Manage(s).MarkTrue(ConditionServiceReady)
}

func serviceToAddress(svc *corev1.Service) (string, error) {
	if len(svc.Spec.Ports) != 1 {
		return "", errors.New("service contains more than one port")
	}

	return "dns:///" + svc.Name + "." + svc.Namespace + ":" + strconv.Itoa(int(svc.Spec.Ports[0].Port)), nil
}
