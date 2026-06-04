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
)

const (
	ResourceKindNatsCredentialRequest = "NatsCredentialRequest"
	ResourceNatsCredentialRequest     = "natscredentialrequest"
	ResourceNatsCredentialRequests    = "natscredentialrequests"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=create
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=natscredentialrequests,singular=natscredentialrequest,scope=Cluster
type NatsCredentialRequest struct {
	metav1.TypeMeta `json:",inline"`
	// Request describes the attributes for the nats credential request.
	// +optional
	Request *NatsCredentialRequestRequest `json:"request,omitempty"`
	// Response describes the attributes for the nats credential response.
	// +optional
	Response *NatsCredentialRequestResponse `json:"response,omitempty"`
}

type NatsCredentialRequestRequest struct {
	Features string `json:"features,omitempty"`
	License  []byte `json:"license,omitempty"`
}

type NatsCredentialRequestResponse struct {
	NatsConfig `json:",inline"`
	Credential []byte `json:"credential,omitempty"`
}

type NatsConfig struct {
	Subject string `json:"natsSubject"`
	Server  string `json:"natsServer"`
}

func init() {
	SchemeBuilder.Register(&NatsCredentialRequest{})
}
