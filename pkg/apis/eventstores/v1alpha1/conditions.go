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
	"knative.dev/pkg/apis"
)

// status conditions
const (
	// ConditionReady has status True when the store is ready to receive requests.
	ConditionReady = apis.ConditionReady
	// ConditionDeploymentReady has status True when the store's adapter is up and running.
	ConditionDeploymentReady apis.ConditionType = "DeploymentReady"
	// ConditionServiceReady has status True when the store's adapter is up and running.
	ConditionServiceReady apis.ConditionType = "ServiceReady"
)

// reasons for conditions
const (
	// ReasonUnavailable is set on an object Ready condition when the resource in unavailable.
	ReasonUnavailable = "ResourceUnavailable"

	// ReasonNotFound is set on an object Ready condition when the resource is not found.
	ReasonNotFound = "NotFound"
)
