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

// +genduck
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventStoreClientSpec defines the fields that clients need to configure
// to access EventStorage servers.
type EventStoreClientSpec struct {
	// EventStoreConnection to the store instance
	// +optional
	EventStoreConnection *EventStoreConnection `json:"eventStore,omitempty"`
}

// EventStoreConnection contains the data to connect to
// an EventStore instance
type EventStoreConnection struct {
	// URI is the gRPC location to the EventStore
	URI string `json:"uri"`
}
