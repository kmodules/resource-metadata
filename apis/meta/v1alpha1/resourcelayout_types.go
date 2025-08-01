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
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/resource-metadata/apis/shared"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	helmshared "x-helm.dev/apimachinery/apis/shared"
)

const (
	ResourceKindResourceLayout = "ResourceLayout"
	ResourceResourceLayout     = "resourcelayout"
	ResourceResourceLayouts    = "resourcelayouts"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=get
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=resourcelayouts,singular=resourcelayout,scope=Cluster
type ResourceLayout struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ResourceLayoutSpec `json:"spec,omitempty"`
}

type ResourceLayoutSpec struct {
	Resource      kmapi.ResourceID            `json:"resource"`
	DefaultLayout bool                        `json:"defaultLayout"`
	Header        *PageBlockLayout            `json:"header,omitempty"`
	TabBar        *PageBlockLayout            `json:"tabBar,omitempty"`
	Pages         []ResourcePageLayout        `json:"pages,omitempty"`
	UI            *shared.UIParameterTemplate `json:"ui,omitempty"`
}

type ResourcePageLayout struct {
	Name string `json:"name"`
	// +optional
	RequiredFeatureSets map[string]FeatureList `json:"requiredFeatureSets,omitempty"`
	Sections            []SectionLayout        `json:"sections,omitempty"`
}

type SectionLayout struct {
	Name    string                 `json:"name,omitempty"`
	Icons   []helmshared.ImageSpec `json:"icons,omitempty"`
	Info    *PageBlockLayout       `json:"info,omitempty"`
	Insight *PageBlockLayout       `json:"insight,omitempty"`
	Blocks  []PageBlockLayout      `json:"blocks,omitempty"`
	// +optional
	RequiredFeatureSets map[string]FeatureList `json:"requiredFeatureSets,omitempty"`
}

type PageBlockLayout struct {
	Kind      TableKind              `json:"kind"` // Connection | Subtable(Field)
	Name      string                 `json:"name,omitempty"`
	Width     int                    `json:"width,omitempty"`
	Icons     []helmshared.ImageSpec `json:"icons,omitempty"`
	FieldPath string                 `json:"fieldPath,omitempty"`

	*shared.ResourceLocator `json:",inline,omitempty"`
	DisplayMode             ResourceDisplayMode `json:"displayMode,omitempty"`
	Actions                 *ResourceActions    `json:"actions,omitempty"`

	View *PageBlockTableDefinition `json:"view,omitempty"`

	RequiredFeatureSets map[string]FeatureList `json:"requiredFeatureSets,omitempty"`
	Filters             map[string]bool        `json:"filters,omitempty"`
}

type FeatureList []string

type PageBlockTableDefinition struct {
	Columns []ResourceColumnDefinition `json:"columns,omitempty"`
	Sort    *TableSortOption           `json:"sort,omitempty"`
}

type TableSortOption struct {
	Order     TableSortOrder `json:"order,omitempty"`
	FieldName string         `json:"fieldName,omitempty"`
}

// +kubebuilder:validation:Enum=Ascending;Descending

type TableSortOrder string

const (
	TableSortOrderAscending  TableSortOrder = "Ascending"
	TableSortOrderDescending TableSortOrder = "Descending"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

type ResourceLayoutList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceLayout `json:"items,omitempty"`
}
