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

package v1alpha1

import (
	"fmt"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/tools/clusterid"
)

func (v ResourceDescriptor) GetKey() string {
	return ResourceResourceDescriptor + "/" + v.Namespace + "/" + v.Name
}

func (v ResourceDescriptor) PolicyName() string {
	cluster := "-"
	if clusterid.ClusterName() != "" {
		cluster = clusterid.ClusterName()
	}
	return fmt.Sprintf("k8s.%s.%s.%s", cluster, v.Namespace, v.Name)
}

func (v ResourceDescriptor) OffshootSelectors() map[string]string {
	return map[string]string{
		"app":          "vault",
		"vault_policy": v.Name,
	}
}

func (v ResourceDescriptor) OffshootLabels() map[string]string {
	return meta_util.FilterKeys("kubevault.com", v.OffshootSelectors(), v.Labels)
}

func (v ResourceDescriptor) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceResourceDescriptors,
		Singular:      ResourceResourceDescriptor,
		Kind:          ResourceKindResourceDescriptor,
		ShortNames:    []string{"vp"},
		Categories:    []string{"vault", "policy", "appscode", "all"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "vault"},
		},
		SpecDefinitionName:      "kmodules.xyz/resource-metadata/apis/meta/v1alpha1.ResourceDescriptor",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: true,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
			{
				Name:     "Phase",
				Type:     "string",
				JSONPath: ".status.phase",
			},
			{
				Name:     "Age",
				Type:     "date",
				JSONPath: ".metadata.creationTimestamp",
			},
		},
	})
}

func (v ResourceDescriptor) IsValid() error {
	return nil
}
