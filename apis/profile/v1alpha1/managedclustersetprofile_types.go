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
	uiapi "kmodules.xyz/resource-metadata/apis/ui/v1alpha1"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindManagedClusterSetProfile = "ManagedClusterSetProfile"
	ResourceManagedClusterSetProfile     = "managedclustersetprofile"
	ResourceManagedClusterSetProfiles    = "managedclustersetprofiles"
)

// ManagedClusterSetProfileSpec defines the desired state of ManagedClusterSetProfile
type ManagedClusterSetProfileSpec struct {
	Namespaces []string               `json:"namespaces,omitempty"`
	Features   map[string]FeatureSpec `json:"features,omitempty"`
}

type FeatureSpec struct {
	// Chart specifies the chart information that will be used by the FluxCD to install the respective feature
	// +optional
	Chart uiapi.ChartInfo `json:"chart,omitempty"`
	// ValuesFrom holds references to resources containing Helm values for this HelmRelease,
	// and information about how they should be merged.
	ValuesFrom []uiapi.ValuesReference `json:"valuesFrom,omitempty"`
	// Values holds the values for this Helm release.
	// +optional
	Values *apiextensionsv1.JSON `json:"values,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ManagedClusterSetProfile is the Schema for the managedclustersetprofiles API
type ManagedClusterSetProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ManagedClusterSetProfileSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// ManagedClusterSetProfileList contains a list of ManagedClusterSetProfile
type ManagedClusterSetProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedClusterSetProfile `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagedClusterSetProfile{}, &ManagedClusterSetProfileList{})
}
