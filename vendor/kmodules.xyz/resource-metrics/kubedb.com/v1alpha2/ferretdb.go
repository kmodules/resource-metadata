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

package v1alpha2

import (
	"fmt"

	"kmodules.xyz/resource-metrics/api"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	api.Register(schema.GroupVersionKind{
		Group:   "kubedb.com",
		Version: "v1alpha2",
		Kind:    "FerretDB",
	}, FerretDB{}.ResourceCalculator())
}

type FerretDB struct{}

func (r FerretDB) ResourceCalculator() api.ResourceCalculator {
	return &api.ResourceCalculatorFuncs{
		AppRoles:               []api.PodRole{api.PodRolePrimary, api.PodRoleSecondary},
		RuntimeRoles:           []api.PodRole{api.PodRolePrimary, api.PodRoleSecondary, api.PodRoleExporter},
		RoleReplicasFn:         r.roleReplicasFn,
		ModeFn:                 r.modeFn,
		UsesTLSFn:              r.usesTLSFn,
		RoleResourceLimitsFn:   r.roleResourceFn(api.ResourceLimits),
		RoleResourceRequestsFn: r.roleResourceFn(api.ResourceRequests),
	}
}

func (r FerretDB) roleReplicasFn(obj map[string]interface{}) (api.ReplicaList, error) {
	replicas, found, err := unstructured.NestedInt64(obj, "spec", "server", "primary", "replicas")
	if err != nil {
		return nil, fmt.Errorf("failed to read spec.replicas %v: %w", obj, err)
	}
	if !found {
		return api.ReplicaList{api.PodRolePrimary: 1}, nil
	}

	ret := api.ReplicaList{api.PodRolePrimary: replicas}
	secRepplicas, found, err := unstructured.NestedInt64(obj, "spec", "server", "secondary", "replicas")
	if found && err == nil {
		ret[api.PodRoleSecondary] = secRepplicas
	}
	return ret, nil
}

func (r FerretDB) modeFn(obj map[string]interface{}) (string, error) {
	_, found, err := unstructured.NestedFieldNoCopy(obj, "spec", "server", "secondary")
	if !found || err != nil {
		return DBModePrimaryOnly, nil
	}
	return DBModeCluster, nil
}

func (r FerretDB) usesTLSFn(obj map[string]interface{}) (bool, error) {
	_, found, err := unstructured.NestedFieldNoCopy(obj, "spec", "tls")
	return found, err
}

func (r FerretDB) roleResourceFn(fn func(rr core.ResourceRequirements) core.ResourceList) func(obj map[string]interface{}) (map[api.PodRole]api.PodInfo, error) {
	return func(obj map[string]interface{}) (map[api.PodRole]api.PodInfo, error) {
		pc, pr, err := api.AppNodeResourcesV2(obj, fn, FerretDBContainerName, "spec", "server", "primary")
		if err != nil {
			return nil, err
		}

		exporter, err := api.ContainerResources(obj, fn, "spec", "monitor", "prometheus", "exporter")
		if err != nil {
			return nil, err
		}

		ret := map[api.PodRole]api.PodInfo{
			api.PodRolePrimary:  {Resource: pc, Replicas: pr},
			api.PodRoleExporter: {Resource: exporter, Replicas: pr},
		}

		_, found, err := unstructured.NestedFieldNoCopy(obj, "spec", "server", "secondary")
		if found && err == nil {
			sc, sr, err := api.AppNodeResourcesV2(obj, fn, FerretDBContainerName, "spec", "server", "secondary")
			if err != nil {
				return nil, err
			}
			sc[core.ResourceStorage] = pc[core.ResourceStorage]
			ret[api.PodRoleSecondary] = api.PodInfo{Resource: sc, Replicas: sr}
			ret[api.PodRoleExporter] = api.PodInfo{Resource: exporter, Replicas: pr + sr}
		}
		return ret, nil
	}
}
