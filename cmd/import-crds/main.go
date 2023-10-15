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

package main

import (
	"bytes"
	"encoding/json"
	goflag "flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/tools/parser"
	rsapi "kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	uiapi "kmodules.xyz/resource-metadata/apis/ui/v1alpha1"

	flag "github.com/spf13/pflag"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

/*
# Add icons
go run cmd/icon-namer/main.go

# Import crds

go run cmd/import-crds/main.go --input=$HOME/go/src/kubeops.dev/supervisor/crds

go run cmd/import-crds/main.go --input=$HOME/go/src/kubeops.dev/ui-server/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/stash.appscode.dev/ui-server/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/go.openviz.dev/grafana-tools/crds

go run cmd/import-crds/main.go --input=$HOME/go/src/k8s.io/api/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/k8s.io/kube-aggregator/crds

go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/prometheus-operator/prometheus-operator/example/prometheus-operator-crd
go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/jetstack/cert-manager/deploy/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/voyagermesh.dev/apimachinery/crds

go run cmd/import-crds/main.go --input=$HOME/go/src/stash.appscode.dev/apimachinery/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/kmodules.xyz/custom-resources/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/kmodules.xyz/resource-metadata/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/kubedb.dev/apimachinery/crds

go run cmd/import-crds/main.go --input=$HOME/go/src/kubevault.dev/apimachinery/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/sigs.k8s.io/secrets-store-csi-driver/charts/secrets-store-csi-driver/crds

go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/kubernetes-csi/external-snapshotter/client/config/crd

go run cmd/import-crds/main.go --input=$HOME/go/src/k8s.io/autoscaler/vertical-pod-autoscaler/deploy/vpa-v1-crd-gen.yaml
go run cmd/import-crds/main.go --input=$HOME/go/src/k8s.io/autoscaler/vertical-pod-autoscaler/deploy/vpa-v1-crd.yaml

# FluxCD CRDs
helm template flux fluxcd-community/flux2 --output-dir=/tmp/fluxcd-manifests
go run cmd/import-crds/main.go --input=/tmp/fluxcd-manifests

# x-helm/kubepack
go run cmd/import-crds/main.go --input=$HOME/go/src/x-helm.dev/apimachinery/crds

# OPA
go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/open-policy-agent/gatekeeper/config/crd/bases

# External DNS
go run cmd/import-crds/main.go --input=$HOME/go/src/kubeops.dev/external-dns-operator/crds

# CAPI / Managed DB
go run cmd/import-crds/main.go --input=$HOME/go/src/sigs.k8s.io/cluster-api/config/crd/bases
go run cmd/import-crds/main.go --input=$HOME/go/src/sigs.k8s.io/cluster-api-provider-aws/config/crd/bases
go run cmd/import-crds/main.go --input=$HOME/go/src/sigs.k8s.io/cluster-api-provider-azure/config/crd/bases
go run cmd/import-crds/main.go --input=$HOME/go/src/sigs.k8s.io/cluster-api-provider-gcp/config/crd/bases

# crossplane
go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/crossplane/crossplane/cluster/crds

# kubeform
go run cmd/import-crds/main.go --input=$HOME/go/src/kubeform.dev/provider-aws/package/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/kubeform.dev/provider-azure/package/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/kubeform.dev/provider-gcp/package/crds

# upbound providers
go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/upbound/provider-aws/package/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/upbound/provider-azure/package/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/upbound/provider-gcp/package/crds

# ocm
go run cmd/import-crds/main.go --input=$HOME/go/src/open-cluster-management.io/api/cluster/v1
go run cmd/import-crds/main.go --input=$HOME/go/src/open-cluster-management.io/api/cluster/v1beta2
go run cmd/import-crds/main.go --input=$HOME/go/src/open-cluster-management.io/api/cluster/v1beta1
go run cmd/import-crds/main.go --input=$HOME/go/src/open-cluster-management.io/api/cluster/v1alpha1
go run cmd/import-crds/main.go --input=$HOME/go/src/open-cluster-management.io/api/work/v1
go run cmd/import-crds/main.go --input=$HOME/go/src/open-cluster-management.io/api/work/v1alpha1
go run cmd/import-crds/main.go --input=$HOME/go/src/open-cluster-management.io/api/crdsv1beta1
go run cmd/import-crds/main.go --input=$HOME/go/src/open-cluster-management.io/api/addon/v1alpha1
go run cmd/import-crds/main.go --input=$HOME/go/src/open-cluster-management.io/api/operator/v1

# Gateway
go run cmd/import-crds/main.go --input=$HOME/go/src/voyagermesh.dev/installer/charts/gateway-helm/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/voyagermesh.dev/installer/charts/gateway-helm/crds/generated
go run cmd/import-crds/main.go --input=$HOME/go/src/voyagermesh.dev/installer/charts/gateway-helm/crds/voyager

# kubeware
go run cmd/import-crds/main.go --input=$HOME/go/src/go.kubeware.dev/catalog/config/crd/bases

# falco
go run cmd/import-crds/main.go --input=$HOME/go/src/kubeops.dev/falco-ui-server/crds
*/
func main() {
	var input []string
	flag.StringSliceVar(&input, "input", nil, "List of crd urls or dir/files")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	for _, location := range input {
		err := processLocation(location)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func CustomResourceDefinition(data []byte) (*crdv1.CustomResourceDefinition, error) {
	var tm metav1.TypeMeta
	err := yaml.Unmarshal(data, &tm)
	if err != nil {
		return nil, err
	}

	if tm.APIVersion == crdv1.SchemeGroupVersion.String() {
		var out crdv1.CustomResourceDefinition
		err := yaml.Unmarshal(data, &out)
		if err != nil {
			return nil, err
		}
		return &out, nil
	}

	var defv1 crdv1beta1.CustomResourceDefinition
	err = yaml.Unmarshal(data, &defv1)
	if err != nil {
		return nil, err
	}

	var inner apiextensions.CustomResourceDefinition
	err = crdv1beta1.Convert_v1beta1_CustomResourceDefinition_To_apiextensions_CustomResourceDefinition(&defv1, &inner, nil)
	if err != nil {
		return nil, err
	}

	var out crdv1.CustomResourceDefinition
	err = crdv1.Convert_apiextensions_CustomResourceDefinition_To_v1_CustomResourceDefinition(&inner, &out, nil)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func processLocation(location string) error {
	u, err := url.Parse(location)
	if err != nil {
		return err
	}

	if u.Scheme != "" {
		resp, err := http.Get(u.String())
		if err != nil {
			return err
		}
		defer func() { _ = resp.Body.Close() }()
		var buf bytes.Buffer
		_, err = io.Copy(&buf, resp.Body)
		if err != nil {
			return err
		}
		return parser.ProcessResources(buf.Bytes(), processObject)
	} else {
		fi, err := os.Stat(location)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			ext := crdFileExtension(location)

			err = filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				if !strings.HasSuffix(info.Name(), ext) {
					return nil
				}

				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				return parser.ProcessResources(data, processObject)
			})
			if err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(location)
			if err != nil {
				return err
			}
			crd, err := CustomResourceDefinition(data)
			if err != nil {
				return err
			}
			err = WriteDescriptor(crd, filepath.Join("hub", rsapi.ResourceResourceDescriptors))
			if err != nil {
				return err
			}
			err = WriteEditor(crd, filepath.Join("hub", uiapi.ResourceResourceEditors))
			if err != nil {
				return err
			}
			err = WriteTableDefinition(crd, filepath.Join("hub", rsapi.ResourceResourceTableDefinitions))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func processObject(ri parser.ResourceInfo) error {
	data, err := json.Marshal(ri.Object)
	if err != nil {
		return err
	}
	crd, err := CustomResourceDefinition(data)
	if err != nil {
		return err
	}
	err = WriteDescriptor(crd, filepath.Join("hub", rsapi.ResourceResourceDescriptors))
	if err != nil {
		return err
	}
	err = WriteEditor(crd, filepath.Join("hub", uiapi.ResourceResourceEditors))
	if err != nil {
		return err
	}
	err = WriteTableDefinition(crd, filepath.Join("hub", rsapi.ResourceResourceTableDefinitions))
	if err != nil {
		return err
	}
	return nil
}

func crdFileExtension(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".v1.yaml") {
			return ".v1.yaml"
		}
	}
	return ".yaml"
}

func WriteDescriptor(crd *crdv1.CustomResourceDefinition, dir string) error {
	for _, v := range crd.Spec.Versions {
		version := v.Name

		kind := crd.Spec.Names.Kind
		plural := crd.Spec.Names.Plural

		name := fmt.Sprintf("%s-%s-%s", crd.Spec.Group, version, plural)
		baseDir := filepath.Join(dir, crd.Spec.Group, version)

		err := os.MkdirAll(baseDir, 0o755)
		if err != nil {
			return err
		}

		filename := filepath.Join(baseDir, plural+".yaml")

		var rd rsapi.ResourceDescriptor
		if existing, err := os.ReadFile(filename); os.IsNotExist(err) {
			rd = rsapi.ResourceDescriptor{
				TypeMeta: metav1.TypeMeta{
					APIVersion: rsapi.SchemeGroupVersion.String(),
					Kind:       rsapi.ResourceKindResourceDescriptor,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
					Labels: map[string]string{
						"k8s.io/group":    crd.Spec.Group,
						"k8s.io/version":  version,
						"k8s.io/resource": plural,
						"k8s.io/kind":     kind,
					},
				},
				Spec: rsapi.ResourceDescriptorSpec{
					Resource: kmapi.ResourceID{
						Group:   crd.Spec.Group,
						Version: version,
						Name:    plural,
						Kind:    kind,
						Scope:   kmapi.ResourceScope(crd.Spec.Scope),
					},
					Validation: v.Schema,
				},
			}
		} else {
			err = yaml.Unmarshal(existing, &rd)
			if err == nil {
				rd.Spec.Validation = v.Schema
			}
		}

		if rd.Spec.Validation != nil &&
			rd.Spec.Validation.OpenAPIV3Schema != nil {

			var mc crdv1.JSONSchemaProps
			err = yaml.Unmarshal([]byte(rsapi.ObjectMetaSchema), &mc)
			if err != nil {
				return err
			}
			if rd.Spec.Resource.Scope == kmapi.ClusterScoped {
				delete(mc.Properties, "namespace")
			}
			rd.Spec.Validation.OpenAPIV3Schema.Properties["metadata"] = mc
			delete(rd.Spec.Validation.OpenAPIV3Schema.Properties, "status")
		}

		data, err := yaml.Marshal(rd)
		if err != nil {
			return err
		}

		data, err = rsapi.FormatMetadata(data)
		if err != nil {
			return err
		}

		err = os.WriteFile(filename, data, 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteEditor(crd *crdv1.CustomResourceDefinition, dir string) error {
	for _, v := range crd.Spec.Versions {
		version := v.Name

		kind := crd.Spec.Names.Kind
		plural := crd.Spec.Names.Plural

		name := fmt.Sprintf("%s-%s-%s", crd.Spec.Group, version, plural)
		baseDir := filepath.Join(dir, crd.Spec.Group, version)

		err := os.MkdirAll(baseDir, 0o755)
		if err != nil {
			return err
		}

		filename := filepath.Join(baseDir, plural+".yaml")

		var ed uiapi.ResourceEditor
		if _, err := os.ReadFile(filename); os.IsNotExist(err) {
			ed = uiapi.ResourceEditor{
				TypeMeta: metav1.TypeMeta{
					APIVersion: uiapi.SchemeGroupVersion.String(),
					Kind:       uiapi.ResourceKindResourceEditor,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
					Labels: map[string]string{
						"k8s.io/group":    crd.Spec.Group,
						"k8s.io/version":  version,
						"k8s.io/resource": plural,
						"k8s.io/kind":     kind,
					},
				},
				Spec: uiapi.ResourceEditorSpec{
					Resource: kmapi.ResourceID{
						Group:   crd.Spec.Group,
						Version: version,
						Name:    plural,
						Kind:    kind,
						Scope:   kmapi.ResourceScope(crd.Spec.Scope),
					},
				},
			}

			data, err := yaml.Marshal(ed)
			if err != nil {
				return err
			}

			data, err = uiapi.FormatMetadata(data)
			if err != nil {
				return err
			}

			err = os.WriteFile(filename, data, 0o644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteTableDefinition(crd *crdv1.CustomResourceDefinition, dir string) error {
	for _, v := range crd.Spec.Versions {
		version := v.Name

		kind := crd.Spec.Names.Kind
		plural := crd.Spec.Names.Plural

		name := fmt.Sprintf("%s-%s-%s", crd.Spec.Group, version, plural)
		baseDir := filepath.Join(dir, crd.Spec.Group, version)

		err := os.MkdirAll(baseDir, 0o755)
		if err != nil {
			return err
		}

		filename := filepath.Join(baseDir, plural+".yaml")

		var td rsapi.ResourceTableDefinition
		if _, err := os.ReadFile(filename); os.IsNotExist(err) {
			td = rsapi.ResourceTableDefinition{
				TypeMeta: metav1.TypeMeta{
					APIVersion: rsapi.SchemeGroupVersion.String(),
					Kind:       rsapi.ResourceKindResourceTableDefinition,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
					Labels: map[string]string{
						"k8s.io/group":    crd.Spec.Group,
						"k8s.io/version":  version,
						"k8s.io/resource": plural,
						"k8s.io/kind":     kind,
					},
				},
				Spec: rsapi.ResourceTableDefinitionSpec{
					DefaultView: true,
					Resource: &kmapi.ResourceID{
						Group:   crd.Spec.Group,
						Version: version,
						Name:    plural,
						Kind:    kind,
						Scope:   kmapi.ResourceScope(crd.Spec.Scope),
					},
				},
			}

			data, err := yaml.Marshal(td)
			if err != nil {
				return err
			}

			data, err = rsapi.FormatMetadata(data)
			if err != nil {
				return err
			}

			err = os.WriteFile(filename, data, 0o644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
