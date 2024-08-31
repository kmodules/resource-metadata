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
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
)

// GenericBindingSpec defines the desired state of GenericBinding
type GenericBindingSpec struct {
	// SourceRef refers to the source app instance.
	SourceRef kmapi.ObjectReference `json:"sourceRef"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GenericBinding is the Schema for the generic binding API
type GenericBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GenericBindingSpec `json:"spec,omitempty"`
	Status BindingStatus      `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GenericBindingList contains a list of GenericBinding
type GenericBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GenericBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GenericBinding{}, &GenericBindingList{})
}

var _ BindingInterface = &GenericBinding{}

func (in *GenericBinding) GetStatus() *BindingStatus {
	return &in.Status
}

func (in *GenericBinding) GetConditions() kmapi.Conditions {
	return in.Status.Conditions
}

func (in *GenericBinding) SetConditions(conditions kmapi.Conditions) {
	in.Status.Conditions = conditions
}

func (dst *GenericBinding) Duckify(srcRaw runtime.Object) error {
	switch src := srcRaw.(type) {
	case *DruidBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindDruidBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *ElasticsearchBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindElasticsearchBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *FerretDBBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindFerretDBBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *KafkaBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindKafkaBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *MariaDBBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindMariaDBBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *MemcachedBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindMemcachedBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *MSSQLServerBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindMSSQLServerBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *MongoDBBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindMongoDBBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *MySQLBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindMySQLBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *PerconaXtraDBBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindPerconaXtraDBBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *PgBouncerBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindPgBouncerBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *PgpoolBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindPgpoolBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *PostgresBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindPostgresBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *ProxySQLBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindProxySQLBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *RabbitMQBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindRabbitMQBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *RedisBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindRedisBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *SinglestoreBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindSinglestoreBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *SolrBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindSolrBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	case *ZooKeeperBinding:
		dst.TypeMeta = metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       ResourceKindZooKeeperBinding,
		}
		dst.ObjectMeta = src.ObjectMeta
		dst.Spec.SourceRef = src.Spec.SourceRef
		dst.Status = src.Status
		return nil
	}
	return fmt.Errorf("unknown src type %T", srcRaw)
}
