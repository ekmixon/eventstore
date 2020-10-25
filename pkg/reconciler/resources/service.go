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

package resources

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// NewService creates a Service object.
func NewService(ns, name string, opts ...ObjectOption) *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      name,
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// AddServiceSelector adds a label selector to the Service spec.
func AddServiceSelector(key, val string) ObjectOption {
	return func(object interface{}) {
		s := object.(*corev1.Service)

		selector := &s.Spec.Selector
		if *selector == nil {
			*selector = make(map[string]string, 1)
		}
		(*selector)[key] = val
	}
}

// ServicePort adds a port to a Service.
func ServicePort(name string, port int32) ObjectOption {
	return func(object interface{}) {
		s := object.(*corev1.Service)
		s.Spec.Ports = []corev1.ServicePort{{
			Name:       name,
			Port:       port,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.IntOrString{IntVal: port},
		}}
	}
}
