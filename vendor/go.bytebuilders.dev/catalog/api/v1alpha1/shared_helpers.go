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
	"context"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	cutil "kmodules.xyz/client-go/conditions"
	"kmodules.xyz/client-go/meta"
	kvm_apis "kubevault.dev/apimachinery/apis"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func GetSecretEngineName(obj client.Object, dbName string) string {
	return meta.NameWithSuffix(dbName, getSuffix(obj)+"-engine")
}

func GetDatabaseRoleName(obj client.Object) string {
	return meta.NameWithSuffix(obj.GetName(), getSuffix(obj)+"-role")
}

func GetSecretAccessRequestName(obj client.Object) string {
	return meta.NameWithSuffix(obj.GetName(), getSuffix(obj)+"-req")
}

func getSuffix(obj client.Object) string {
	kind := obj.GetObjectKind().GroupVersionKind().Kind
	return strings.TrimSuffix(strings.ToLower(kind), "binding")
}

func GetSourceRefStatus(kc client.Client, sourceRef client.Object) (*runtime.RawExtension, error) {
	var src unstructured.Unstructured
	src.SetGroupVersionKind(sourceRef.GetObjectKind().GroupVersionKind())

	err := kc.Get(context.TODO(), types.NamespacedName{
		Namespace: sourceRef.GetNamespace(),
		Name:      sourceRef.GetName(),
	}, &src)
	if err != nil {
		return nil, err
	}

	if v, ok, err := unstructured.NestedFieldNoCopy(src.UnstructuredContent(), "status"); err == nil && ok {
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return &runtime.RawExtension{Raw: raw}, nil
	} else {
		return nil, err
	}
}

func IsAccessRequestExpired(conditions []kmapi.Condition) bool {
	idx, cond := cutil.GetCondition(conditions, string(BindingConditionTypeSecretAccessRequestReady))
	if idx == -1 {
		return false
	}
	return cond.Reason == BindingConditionReasonSecretAccessRequestExpired
}

func IsEngineAPIResourcesReady(conditions []kmapi.Condition) bool {
	se := cutil.IsConditionTrue(conditions, string(BindingConditionTypeSecretEngineReady))
	sar := cutil.IsConditionTrue(conditions, string(BindingConditionTypeSecretAccessRequestReady))
	role := cutil.IsConditionTrue(conditions, string(BindingConditionTypeRoleReady))

	return sar && role && se
}

func IsEngineAPIResourcesConditionDetermined(conditions []kmapi.Condition) bool {
	roleIndex, _ := cutil.GetCondition(conditions, string(BindingConditionTypeRoleReady))
	sarIndex, _ := cutil.GetCondition(conditions, string(BindingConditionTypeSecretAccessRequestReady))

	return (roleIndex != -1) && (sarIndex != -1)
}

func IsVaultReady(conditions []kmapi.Condition) bool {
	return cutil.IsConditionTrue(conditions, kvm_apis.AllReplicasAreReady) &&
		cutil.IsConditionTrue(conditions, kvm_apis.VaultServerAcceptingConnection) &&
		cutil.IsConditionTrue(conditions, kvm_apis.VaultServerInitialized) &&
		cutil.IsConditionTrue(conditions, kvm_apis.VaultServerUnsealed)
}

func IsSecretEngineReady(conditions []kmapi.Condition) bool {
	return cutil.IsConditionTrue(conditions, cutil.ConditionAvailable)
}

func IsRoleAvailable(conditions []kmapi.Condition) bool {
	return cutil.IsConditionTrue(conditions, cutil.ConditionAvailable)
}

func GetFinalizer() string {
	return GroupVersion.Group
}

func ConditionsOrder() []kmapi.ConditionType {
	return []kmapi.ConditionType{
		BindingConditionTypeDBReady,
		BindingConditionTypeVaultReady,
		BindingConditionTypeServiceAccountReady,
		BindingConditionTypeSecretEngineReady,
		BindingConditionTypeRoleReady,
		BindingConditionTypeSecretAccessRequestReady,
	}
}

func GetPhase(obj BindingInterface) BindingPhase {
	if !obj.GetDeletionTimestamp().IsZero() {
		return BindingPhaseTerminating
	}
	conditions := obj.GetConditions()
	var cond kmapi.Condition
	for i, _ := range conditions {
		c := conditions[i]
		if c.Type == kmapi.ReadyCondition {
			cond = c
			break
		}
	}
	if cond.Type != kmapi.ReadyCondition {
		klog.Errorf("no Ready condition in the status for %s/%s", obj.GetNamespace(), obj.GetName())
		return BindingPhasePending
	}

	if cond.Status == metav1.ConditionTrue {
		return BindingPhaseCurrent
	}
	if cond.Reason == BindingConditionReasonSecretAccessRequestExpired {
		return BindingPhaseExpired
	}

	// DB, Vault & SA should exists & Ready
	if cond.Reason == BindingConditionReasonDBNotCreated ||
		cond.Reason == BindingConditionReasonDBProvisioning ||
		cond.Reason == BindingConditionReasonVaultNotCreated ||
		cond.Reason == BindingConditionReasonVaultProvisioning ||
		cond.Reason == BindingConditionReasonServiceAccountNotCreated {
		return BindingPhasePending
	}

	// Engine, Role, AccessRequest
	if cond.Reason == BindingConditionReasonSecretEngineNotCreated ||
		cond.Reason == BindingConditionReasonSecretEngineNotReady ||
		cond.Reason == BindingConditionReasonRoleNotCreated ||
		cond.Reason == BindingConditionReasonRoleNotReady ||
		cond.Reason == BindingConditionReasonSecretAccessRequestNotCreated ||
		cond.Reason == BindingConditionReasonSecretAccessRequestNotReady {
		return BindingPhaseInProgress
	}

	if cond.Reason == BindingConditionReasonSecretAccessRequestDenied {
		return BindingPhaseFailed
	}

	return BindingPhaseCurrent
}
