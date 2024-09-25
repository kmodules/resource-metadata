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

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindManagedClusterProfileBinding = "ManagedClusterProfileBinding"
	ResourceManagedClusterProfileBinding     = "managedclusterprofilebinding"
	ResourceManagedClusterProfileBindings    = "managedclusterprofilebindings"
)

// ManagedClusterProfileBindingSpec defines the desired state of ManagedClusterProfileBinding
type ManagedClusterProfileBindingSpec struct {
	ProfileRef      core.LocalObjectReference `json:"profileRef"`
	ClusterMetadata ClusterMetadata           `json:"clusterMetadata"`
	Features        map[string]FeatureSpec    `json:"features,omitempty"`
}

type ClusterMetadata struct {
	Uid             string   `json:"uid"`
	Name            string   `json:"name"`
	ClusterManagers []string `json:"clusterManagers"`
	// +optional
	CAPI CapiMetadata `json:"capi"`
}

type CapiMetadata struct {
	// +optional
	Provider  CAPIProvider `json:"provider"`
	Namespace string       `json:"namespace"`
}

// +kubebuilder:validation:Enum=capa;capg;capz
type CAPIProvider string

const (
	CAPIProviderDisabled CAPIProvider = ""
	CAPIProviderCAPA     CAPIProvider = "capa"
	CAPIProviderCAPG     CAPIProvider = "capg"
	CAPIProviderCAPZ     CAPIProvider = "capz"
)

type BindingStatusPhase string

const (
	BindingStatusPhaseInProgress BindingStatusPhase = "InProgress"
	BindingStatusPhaseCurrent    BindingStatusPhase = "Current"
	BindingStatusPhaseFailed     BindingStatusPhase = "Failed"
)

// ManagedClusterProfileBindingStatus defines the observed state of ManagedClusterProfileBinding
type ManagedClusterProfileBindingStatus struct {
	// Phase indicates the state this Vault cluster jumps in.
	// +optional
	Phase BindingStatusPhase `json:"phase,omitempty"`
	// Represents the latest available observations of a GrafanaDashboard current state.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
	// ObservedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64             `json:"observedGeneration,omitempty"`
	ManifestWorks      map[string]string `json:"manifestWorks,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ManagedClusterProfileBinding is the Schema for the managedclusterprofilebindings API
type ManagedClusterProfileBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedClusterProfileBindingSpec   `json:"spec,omitempty"`
	Status ManagedClusterProfileBindingStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// ManagedClusterProfileBindingList contains a list of ManagedClusterProfileBinding
type ManagedClusterProfileBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedClusterProfileBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagedClusterProfileBinding{}, &ManagedClusterProfileBindingList{})
}
