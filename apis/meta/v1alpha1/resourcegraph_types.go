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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindResourceGraph = "ResourceGraph"
	ResourceResourceGraph     = "resourcegraph"
	ResourceResourceGraphs    = "resourcegraphs"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ResourceGraph struct {
	metav1.TypeMeta `json:",inline"`
	// Request describes the attributes for the graph request.
	// +optional
	Request *ResourceGraphRequest `json:"request,omitempty"`
	// Response describes the attributes for the graph response.
	// +optional
	Response *ResourceGraphResponse `json:"response,omitempty"`
}

type ResourceGraphRequest struct {
	Source kmapi.ObjectInfo `json:"source"`
}

type ResourceGraphResponse struct {
	Resources   []kmapi.ResourceID `json:"resources"`
	Connections []ObjectConnection `json:"connections"`
}

type ObjectConnection struct {
	Source ObjectPointer `json:"source"`
	Target ObjectPointer `json:"target"`
	Labels []string      `json:"labels"`
}

type ObjectPointer struct {
	ResourceID int    `json:"resourceID"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
}
