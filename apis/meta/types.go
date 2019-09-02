/*
Copyright 2019 The ResourceMetadata Project Authors.

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

package meta

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ReferenceType string

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ResourceDescriptor struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec ResourceDescriptorSpec
}

type ResourceDescriptorSpec struct {
	Resource       ResourceID
	DisplayColumns []ResourceColumnDefinition
	Connections    []ResourceConnection
	KeyTargets     []GroupVersionResource
}

type ResourceID struct {
	Group   string
	Version string
	Name    string
	Kind    string
	Scope   ResourceScope
}

type ResourceScope string

const (
	ClusterScoped   ResourceScope = "Cluster"
	NamespaceScoped ResourceScope = "Namespaced"
)

type ConnectionType string

const (
	MatchSelector ConnectionType = "MatchSelector"
	MatchName     ConnectionType = "MatchName"
	MatchRef      ConnectionType = "MatchRef"
	OwnedBy       ConnectionType = "OwnedBy"
)

type GroupVersionResource struct {
	Group    string
	Version  string
	Resource string
}

type ResourceConnection struct {
	Target GroupVersionResource
	ResourceConnectionSpec
}

type ResourceConnectionSpec struct {
	Type            ConnectionType
	NamespacePath   string
	TargetLabelPath string
	SelectorPath    string
	Selector        *metav1.LabelSelector
	NameTemplate    string
	References      []string
	Level           OwnershipLevel
}

type OwnershipLevel string

const (
	Reference  OwnershipLevel = ""
	Owner      OwnershipLevel = "Owner"
	Controller OwnershipLevel = "Controller"
)

// ResourceColumnDefinition specifies a column for server side printing.
type ResourceColumnDefinition struct {
	// name is a human readable name for the column.
	Name string
	// type is an OpenAPI type definition for this column.
	// See https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#data-types for more.
	Type string
	// format is an optional OpenAPI type definition for this column. The 'name' format is applied
	// to the primary identifier column to assist in clients identifying column is the resource name.
	// See https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#data-types for more.
	Format string
	// description is a human readable description of this column.
	Description string
	// priority is an integer defining the relative importance of this column compared to others. Lower
	// numbers are considered higher priority. Columns that may be omitted in limited space scenarios
	// should be given a higher priority.
	Priority int32

	// JSONPath is a simple JSON path, i.e. without array notation.
	JSONPath string
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceDescriptorList is a list of ResourceDescriptor objects.
type ResourceDescriptorList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []ResourceDescriptor
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PathFinder struct {
	metav1.TypeMeta
	Request  *PathRequest
	Response *PathResponse
}

type PathRequest struct {
	Source GroupVersionResource
	Target *GroupVersionResource
}

type PathResponse struct {
	Paths []Path
}

type Path struct {
	Source   GroupVersionResource
	Target   GroupVersionResource
	Distance uint64
	Edges    []Edge
}

type Edge struct {
	Src        GroupVersionResource
	Dst        GroupVersionResource
	W          uint64
	Connection ResourceConnectionSpec
	Forward    bool
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type GraphFinder struct {
	metav1.TypeMeta
	Request  *GraphRequest
	Response *GraphResponse
}

type GraphRequest struct {
	Source GroupVersionResource
}

type GraphResponse struct {
	Source      GroupVersionResource
	Connections []Edge
}
