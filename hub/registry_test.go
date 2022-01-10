/*
Copyright AppsCode Inc. and Contributors

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

package hub

import (
	"fmt"
	"testing"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub/resourceclasses"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRegistry_LoadDefaultResourceClass(t *testing.T) {
	reg := NewRegistry("some-uid", NewKVLocal())
	panel, err := reg.DefaultResourcePanel(resourceclasses.ClusterUI)
	assert.NoError(t, err)

	for _, rc := range panel.Sections {
		for _, entry := range rc.Entries {
			println(rc.Name)
			fmt.Println(entry)
		}
	}
}

func TestResourcePanel_Minus(t *testing.T) {
	reg := NewRegistryOfKnownResources()
	panel, err := reg.CompleteResourcePanel(resourceclasses.ClusterUI)
	assert.NoError(t, err)

	type fields struct {
		TypeMeta v1.TypeMeta
		Sections []*v1alpha1.PanelSection
	}
	type args struct {
		b *v1alpha1.ResourcePanel
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "default",
			fields: fields{
				TypeMeta: panel.TypeMeta,
				Sections: panel.Sections,
			},
			args: args{
				b: &v1alpha1.ResourcePanel{
					TypeMeta: panel.TypeMeta,
					Sections: []*v1alpha1.PanelSection{
						{
							Name:   "Admissionregistration",
							Weight: 2,
							Entries: []v1alpha1.PanelEntry{
								{
									Name: "MutatingWebhookConfiguration",
									Resource: &kmapi.ResourceID{
										Group:   "admissionregistration.k8s.io",
										Version: "v1beta1",
										Name:    "mutatingwebhookconfigurations",
									},
									Namespaced: false,
								},
								{
									Name: "ValidatingWebhookConfiguration",
									Resource: &kmapi.ResourceID{
										Group:   "admissionregistration.k8s.io",
										Version: "v1beta1",
										Name:    "validatingwebhookconfigurations",
									},
									Namespaced: false,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &v1alpha1.ResourcePanel{
				TypeMeta: tt.fields.TypeMeta,
				Sections: tt.fields.Sections,
			}
			a.Minus(tt.args.b)
		})
	}
}
