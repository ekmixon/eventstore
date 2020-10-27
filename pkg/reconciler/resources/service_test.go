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
	"testing"

	"github.com/google/go-cmp/cmp"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestService(t *testing.T) {
	svc := NewService(tNs, tName,
		AddServiceSelector("test.selector/1", "val1"),
		ServicePort("h2c", 8080),
	)

	expectSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: tNs,
			Name:      tName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"test.selector/1": "val1",
			},
			Ports: []corev1.ServicePort{{
				Name:       "h2c",
				Port:       8080,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.IntOrString{IntVal: 8080},
			}},
		},
	}

	if d := cmp.Diff(expectSvc, svc); d != "" {
		t.Errorf("Unexpected diff: (-:expect, +:got) %s", d)
	}
}