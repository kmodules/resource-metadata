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

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestParseGVR(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    schema.GroupVersionResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "core-v1-services",
			args: args{
				name: "core-v1-services",
			},
			want: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "services",
			},
			wantErr: assert.NoError,
		},
		{
			name: "apps-v1-deployments",
			args: args{
				name: "apps-v1-deployments",
			},
			want: schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			},
			wantErr: assert.NoError,
		},
		{
			name: "charts.x-helm.dev-v1alpha1-clusterchartpresets",
			args: args{
				name: "charts.x-helm.dev-v1alpha1-clusterchartpresets",
			},
			want: schema.GroupVersionResource{
				Group:    "charts.x-helm.dev",
				Version:  "v1alpha1",
				Resource: "clusterchartpresets",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGVR(tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseGVR(%v)", tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, *got, "ParseGVR(%v)", tt.args.name)
		})
	}
}

func Test_toFilename(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "core-v1-services",
			args: args{
				name: "core-v1-services",
			},
			want: "core/v1/services.yaml",
		},
		{
			name: "apps-v1-deployments",
			args: args{
				name: "apps-v1-deployments",
			},
			want: "apps/v1/deployments.yaml",
		},
		{
			name: "charts.x-helm.dev-v1alpha1-clusterchartpresets",
			args: args{
				name: "charts.x-helm.dev-v1alpha1-clusterchartpresets",
			},
			want: "charts.x-helm.dev/v1alpha1/clusterchartpresets.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, toFilename(tt.args.name), "toFilename(%v)", tt.args.name)
		})
	}
}
