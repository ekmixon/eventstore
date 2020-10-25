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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

// PodLabels sets the set of labels of a PodSpecable's Pod template.
func PodLabels(ls labels.Set) ObjectOption {
	return func(object interface{}) {
		var metaObj metav1.Object

		switch o := object.(type) {
		case *appsv1.Deployment:
			metaObj = &o.Spec.Template
		case *servingv1.Service:
			metaObj = &o.Spec.Template
		}

		metaObj.SetLabels(ls)
	}
}

// Container adds a container to a PodSpecable's Pod template.
func Container(c *corev1.Container) ObjectOption {
	return func(object interface{}) {
		switch o := object.(type) {
		case *appsv1.Deployment:
			containers := &o.Spec.Template.Spec.Containers
			*containers = append(*containers, *c)
		case *servingv1.Service:
			containers := &o.Spec.Template.Spec.Containers
			*containers = []corev1.Container{*c}
		}
	}
}
