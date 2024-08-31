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
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindRabbitMQBinding = "RabbitMQBinding"
	ResourceRabbitMQBinding     = "rabbitmqbinding"
	ResourceRabbitMQBindings    = "rabbitmqbindings"
)

// RabbitMQBindingSpec defines the desired state of RabbitMQBinding
type RabbitMQBindingSpec struct {
	// SourceRef refers to the source app instance.
	SourceRef kmapi.ObjectReference `json:"sourceRef"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=rmbinding,categories={binding,kubedb,appscode}
// +kubebuilder:printcolumn:name="Src_NS",type="string",JSONPath=".spec.sourceRef.namespace"
// +kubebuilder:printcolumn:name="Src_Name",type="string",JSONPath=".spec.sourceRef.name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// RabbitMQBinding is the Schema for the rabbitmqbindings API
type RabbitMQBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RabbitMQBindingSpec `json:"spec,omitempty"`
	Status BindingStatus       `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RabbitMQBindingList contains a list of RabbitMQBinding
type RabbitMQBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RabbitMQBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RabbitMQBinding{}, &RabbitMQBindingList{})
}
