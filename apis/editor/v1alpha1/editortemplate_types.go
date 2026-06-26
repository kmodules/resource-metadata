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

// EditorTemplate loads the editor model, manifest and resources for an existing
// installation identified by its metadata. It is the aggregated-API equivalent
// of kubepack.dev/lib-app/pkg/editor.LoadResourceEditorModel (see
// chart_template.go loadEditorModel2).
//
// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=create
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
type EditorTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// Request identifies the installation to load.
	// +optional
	Request *EditorTemplateRequest `json:"request,omitempty"`
	// Response holds the loaded editor template.
	// +optional
	Response *EditorTemplateResponse `json:"response,omitempty"`
}

type EditorTemplateRequest struct {
	releasesapi.ModelMetadata `json:",inline"`
	// ChartRef optionally selects the chart source. When unset, the resource
	// editor chart resolved from the metadata resource is used.
	// +optional
	ChartRef *releasesapi.ChartSourceFlatRef `json:"chartRef,omitempty"`
}

type EditorTemplateResponse struct {
	// Manifest is the rendered manifest of the existing installation.
	// +optional
	Manifest string `json:"manifest,omitempty"`
	// Values is the editor model (values) of the existing installation.
	// +optional
	Values *runtime.RawExtension `json:"values,omitempty"`
	// Resources holds the individual resources of the existing installation.
	// +optional
	Resources []EditorTemplateResource `json:"resources,omitempty"`
}

type EditorTemplateResource struct {
	// +optional
	Filename string `json:"filename,omitempty"`
	// +optional
	Key string `json:"key,omitempty"`
	// +optional
	Data *runtime.RawExtension `json:"data,omitempty"`
}
