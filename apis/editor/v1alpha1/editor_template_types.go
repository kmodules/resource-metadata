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
	ResourceKindEditorTemplate = "EditorTemplate"
	ResourceEditorTemplate     = "editortemplate"
	ResourceEditorTemplates    = "editortemplates"
)

// EditorTemplate loads the editor model, manifest or resources for an existing
// installation identified by Metadata. It is the read-only aggregated-API
// equivalent of the b3
// /clusters/{owner}/{cluster}/helm/editor/{model,manifest,resources} endpoints.
//
// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=create
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=editortemplates,singular=editortemplate,scope=Cluster
type EditorTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// Request identifies the installation and what to load.
	// +optional
	Request *EditorTemplateRequest `json:"request,omitempty"`
	// Response holds the loaded output.
	// +optional
	Response *EditorTemplateResponse `json:"response,omitempty"`
}

type EditorTemplateRequest struct {
	// ChartRef optionally selects the chart source. When unset, the resource
	// editor chart resolved from Metadata is used.
	// +optional
	ChartRef *releasesapi.ChartSourceFlatRef `json:"chartRef,omitempty"`
	// Metadata identifies the installed release (resource id + release name/namespace).
	Metadata releasesapi.Metadata `json:"metadata"`
	// Output selects which artifact to load. Defaults to "resources".
	// +optional
	Output EditorOutput `json:"output,omitempty"`
}

type EditorTemplateResponse struct {
	// Values is set when Output is "model": the editor model (values) of the install.
	// +optional
	Values *runtime.RawExtension `json:"values,omitempty"`
	// Manifest is set when Output is "manifest": the rendered YAML manifest.
	// +optional
	Manifest string `json:"manifest,omitempty"`
	// Resources is set when Output is "resources": the CRDs + resources.
	// +optional
	Resources *releasesapi.ResourceOutput `json:"resources,omitempty"`
}
