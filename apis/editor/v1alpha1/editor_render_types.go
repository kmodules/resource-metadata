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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	releasesapi "x-helm.dev/apimachinery/apis/releases/v1alpha1"
)

const (
	ResourceKindEditorRender = "EditorRender"
	ResourceEditorRender     = "editorrender"
	ResourceEditorRenders    = "editorrenders"
)

// EditorOutput selects which artifact an editor render/template request returns.
type EditorOutput string

const (
	// EditorOutputModel returns the editor model (values).
	EditorOutputModel EditorOutput = "model"
	// EditorOutputManifest returns the rendered YAML manifest.
	EditorOutputManifest EditorOutput = "manifest"
	// EditorOutputResources returns the parsed CRDs + resources.
	EditorOutputResources EditorOutput = "resources"
)

// EditorRender renders an editor model, manifest or resources from a set of
// options, without touching any existing installation. It is the read-only
// aggregated-API equivalent of the b3
// /clusters/{owner}/{cluster}/helm/options/{model,manifest,resources} endpoints.
//
// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=create
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=editorrenders,singular=editorrender,scope=Cluster
type EditorRender struct {
	metav1.TypeMeta `json:",inline"`
	// Request describes what to render.
	// +optional
	Request *EditorRenderRequest `json:"request,omitempty"`
	// Response holds the rendered output.
	// +optional
	Response *EditorRenderResponse `json:"response,omitempty"`
}

type EditorRenderRequest struct {
	// ChartRef optionally selects the chart source. When unset, the resource
	// editor chart resolved from the options metadata is used.
	// +optional
	ChartRef *releasesapi.ChartSourceFlatRef `json:"chartRef,omitempty"`
	// Options is the editor options model (values) to render from.
	// +optional
	Options *runtime.RawExtension `json:"options,omitempty"`
	// Output selects which artifact to render. Defaults to "resources".
	// +optional
	Output EditorOutput `json:"output,omitempty"`
	// SkipCRDs omits CRDs from the resources output.
	// +optional
	SkipCRDs bool `json:"skipCRDs,omitempty"`
}

type EditorRenderResponse struct {
	// Model is set when Output is "model": the generated editor model (values).
	// +optional
	Model *runtime.RawExtension `json:"model,omitempty"`
	// Manifest is set when Output is "manifest": the rendered YAML manifest.
	// +optional
	Manifest string `json:"manifest,omitempty"`
	// Resources is set when Output is "resources": the rendered CRDs + resources.
	// +optional
	Resources *releasesapi.ResourceOutput `json:"resources,omitempty"`
}
