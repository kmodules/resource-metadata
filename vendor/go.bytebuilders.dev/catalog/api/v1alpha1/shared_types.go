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
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	BindingConditionTypeDBReady                  kmapi.ConditionType = "DBReady"
	BindingConditionTypeVaultReady               kmapi.ConditionType = "VaultReady"
	BindingConditionTypeServiceAccountReady      kmapi.ConditionType = "ServiceAccountReady"
	BindingConditionTypeSecretEngineReady        kmapi.ConditionType = "SecretEngineReady"
	BindingConditionTypeRoleReady                kmapi.ConditionType = "RoleReady"
	BindingConditionTypeSecretAccessRequestReady kmapi.ConditionType = "SecretAccessRequestReady"
)

const (
	BindingConditionReasonDBNotCreated   = "DBNotCreated"
	BindingConditionReasonDBProvisioning = "DBProvisioning"

	BindingConditionReasonVaultNotCreated   = "VaultNotCreated"
	BindingConditionReasonVaultProvisioning = "VaultProvisioning"

	BindingConditionReasonServiceAccountNotCreated = "ServiceAccountNotCreated"

	BindingConditionReasonSecretEngineNotCreated = "SecretEngineNotCreated"
	BindingConditionReasonSecretEngineNotReady   = "SecretEngineNotReady"

	BindingConditionReasonRoleNotCreated = "RoleNotCreated"
	BindingConditionReasonRoleNotReady   = "RoleNotReady"

	BindingConditionReasonSecretAccessRequestNotCreated = "SecretAccessRequestNotCreated"
	BindingConditionReasonSecretAccessRequestNotReady   = "SecretAccessRequestNotReady"
	BindingConditionReasonSecretAccessRequestExpired    = "SecretAccessRequestExpired"
	BindingConditionReasonSecretAccessRequestApproved   = "SecretAccessRequestApproved"
	BindingConditionReasonSecretAccessRequestDenied     = "SecretAccessRequestDenied"
)

// +kubebuilder:validation:Enum=Pending;InProgress;Terminating;Current;Failed;Expired
type BindingPhase string

const (
	BindingPhasePending     BindingPhase = "Pending"
	BindingPhaseInProgress  BindingPhase = "InProgress"
	BindingPhaseTerminating BindingPhase = "Terminating"
	BindingPhaseCurrent     BindingPhase = "Current"
	BindingPhaseFailed      BindingPhase = "Failed"
	BindingPhaseExpired     BindingPhase = "Expired"
)

/*Pending : If DB / VaultServer not found
InProgress: If role or accessReq is not ensured yet. Or their phase is not determined yet
Current: all ok, secret is valid
Expired: all ok, secret is expired
Failed: role or accessReq failed for some reason*/

// BindingStatus defines the observed state of App
type BindingStatus struct {
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	Conditions []kmapi.Condition `json:"conditions,omitempty"`

	// Specifies the current phase of the App
	// +optional
	Phase BindingPhase `json:"phase,omitempty"`

	// +optional
	SecretRef *core.LocalObjectReference `json:"secretRef,omitempty"`

	// +optional
	Source *runtime.RawExtension `json:"source,omitempty"`

	// +optional
	Gateway *Gateway `json:"gateway,omitempty"`
}

type Gateway struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	// +optional
	IP string `json:"ip,omitempty"`
	// +optional
	Hostname string `json:"hostname,omitempty"`
	// Services is an optional configuration for services used to expose database
	// +optional
	Services []NamedServiceStatus `json:"services,omitempty"`
	// UI is an optional list of database web uis
	// +optional
	UI []NamedURL `json:"ui,omitempty"`
}

type NamedServiceStatus struct {
	// Alias represents the identifier of the service.
	Alias string `json:"alias"`

	Ports []ofst.GatewayPort `json:"ports"`
}

type NamedURL struct {
	// Alias represents the identifier of the service.
	// This should match the db ui chart name
	Alias string `json:"alias"`

	// URL of the database ui
	URL string `json:"url"`

	Port ofst.GatewayPort `json:"port"`

	// HelmRelease is the name of the helm release used to deploy this ui
	// The name format is typically <alias>-<db-name>
	// +optional
	HelmRelease *core.LocalObjectReference `json:"helmRelease,omitempty"`
}

// +k8s:deepcopy-gen=false
type BindingInterface interface {
	client.Object
	GetStatus() *BindingStatus
	GetConditions() kmapi.Conditions
	SetConditions(conditions kmapi.Conditions)
}
