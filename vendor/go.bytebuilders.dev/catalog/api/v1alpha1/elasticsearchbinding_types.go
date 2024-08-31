/*
Copyright 2023.

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
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindElasticsearchBinding = "ElasticsearchBinding"
	ResourceElasticsearchBinding     = "elasticsearchbinding"
	ResourceElasticsearchBindings    = "elasticsearchbindings"
)

// ElasticsearchBindingSpec defines the desired state of ElasticsearchBinding
type ElasticsearchBindingSpec struct {
	// SourceRef refers to the source app instance.
	SourceRef kmapi.ObjectReference `json:"sourceRef"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=esbinding,categories={binding,kubedb,appscode}
// +kubebuilder:printcolumn:name="Src_NS",type="string",JSONPath=".spec.sourceRef.namespace"
// +kubebuilder:printcolumn:name="Src_Name",type="string",JSONPath=".spec.sourceRef.name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ElasticsearchBinding is the Schema for the elasticsearchbindings API
type ElasticsearchBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticsearchBindingSpec `json:"spec,omitempty"`
	Status BindingStatus            `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ElasticsearchBindingList contains a list of ElasticsearchBinding
type ElasticsearchBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticsearchBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticsearchBinding{}, &ElasticsearchBindingList{})
}

var _ BindingInterface = &ElasticsearchBinding{}

func (in *ElasticsearchBinding) GetStatus() *BindingStatus {
	return &in.Status
}

func (in *ElasticsearchBinding) GetConditions() kmapi.Conditions {
	return in.Status.Conditions
}

func (in *ElasticsearchBinding) SetConditions(conditions kmapi.Conditions) {
	in.Status.Conditions = conditions
}
