/*
Copyright AppsCode Inc. and Contributors.

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
	"testing"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestConvertToStringQuantity(t *testing.T) {
	tests := []struct {
		name    string
		input   core.ResourceList
		wantStr map[core.ResourceName]string
	}{
		{
			name: "memory in Ki converted to Gi",
			input: core.ResourceList{
				core.ResourceMemory: resource.MustParse("437643756Ki"),
			},
			wantStr: map[core.ResourceName]string{
				core.ResourceMemory: "417.37Gi",
			},
		},
		{
			name: "memory in Gi formatted with two decimals",
			input: core.ResourceList{
				core.ResourceMemory: resource.MustParse("16Gi"),
			},
			wantStr: map[core.ResourceName]string{
				core.ResourceMemory: "16Gi",
			},
		},
		{
			name: "memory in Mi formatted with two decimals",
			input: core.ResourceList{
				core.ResourceMemory: resource.MustParse("512Mi"),
			},
			wantStr: map[core.ResourceName]string{
				core.ResourceMemory: "512Mi",
			},
		},
		{
			name: "cpu is unchanged",
			input: core.ResourceList{
				core.ResourceCPU: resource.MustParse("2000m"),
			},
			wantStr: map[core.ResourceName]string{
				core.ResourceCPU: "2",
			},
		},
		{
			name: "memory in bytes converted to Gi",
			input: core.ResourceList{
				core.ResourceMemory: resource.MustParse("23432543534"),
			},
			wantStr: map[core.ResourceName]string{
				core.ResourceMemory: "21.82Gi",
			}, // 4187593113600 / (1024 * 1024 * 1024 * 1024) - 3.80859375 Gi
		},
		{
			name: "memory in bytes converted to Gi",
			input: core.ResourceList{
				core.ResourceMemory: resource.MustParse("4187593113600m"),
			},
			wantStr: map[core.ResourceName]string{
				core.ResourceMemory: "3.9Gi",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := ConvertToStringQuantity(tt.input)
			for name, want := range tt.wantStr {
				got, ok := out[name]
				if !ok {
					t.Errorf("resource %s missing from output", name)
					continue
				}
				if got != want {
					t.Errorf("resource %s: got %q, want %q", name, got, want)
				}
			}
		})
	}
}
