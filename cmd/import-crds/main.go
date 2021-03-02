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
	goflag "flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	flag "github.com/spf13/pflag"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

/*
go run cmd/import-crds/main.go --input=$HOME/go/src/k8s.io/api/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/k8s.io/kube-aggregator/crds

go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/coreos/prometheus-operator/example/prometheus-operator-crd
go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/jetstack/cert-manager/deploy/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/appscode/voyager/api/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/stash.appscode.dev/apimachinery/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/kmodules.xyz/custom-resources/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/kubedb.dev/apimachinery/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/kubevault.dev/operator/api/crds
go run cmd/import-crds/main.go --input=$HOME/go/src/go.searchlight.dev/grafana-operator/crds

go run cmd/import-crds/main.go --input=$HOME/go/src/sigs.k8s.io/application/config/crd/bases

go run cmd/import-crds/main.go --input=$HOME/go/src/github.com/kubernetes-csi/external-snapshotter/client/config/crd

go run cmd/import-crds/main.go --input=$HOME/go/src/k8s.io/autoscaler/vertical-pod-autoscaler/deploy/vpa-v1-crd-gen.yaml
go run cmd/import-crds/main.go --input=$HOME/go/src/k8s.io/autoscaler/vertical-pod-autoscaler/deploy/vpa-v1-crd.yaml

*/
func main() {
	var input []string
	flag.StringSliceVar(&input, "input", nil, "List of crd urls or dir/files")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	dir := filepath.Join("hub", "resourcedescriptors")

	for _, location := range input {
		err := processLocation(location, dir)
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

func processLocation(location, dir string) error {
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
		crd, err := CustomResourceDefinition(buf.Bytes())
		if err != nil {
			return err
		}
		err = WriteDescriptor(crd, dir)
		if err != nil {
			return err
		}
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

				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				crd, err := CustomResourceDefinition(data)
				if err != nil {
					return err
				}
				err = WriteDescriptor(crd, dir)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			data, err := ioutil.ReadFile(location)
			if err != nil {
				return err
			}
			crd, err := CustomResourceDefinition(data)
			if err != nil {
				return err
			}
			err = WriteDescriptor(crd, dir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func crdFileExtension(dir string) string {
	entries, err := ioutil.ReadDir(dir)
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

		err := os.MkdirAll(baseDir, 0755)
		if err != nil {
			return err
		}

		filename := filepath.Join(baseDir, plural+".yaml")

		var rd v1alpha1.ResourceDescriptor
		if existing, err := ioutil.ReadFile(filename); os.IsNotExist(err) {
			rd = v1alpha1.ResourceDescriptor{
				TypeMeta: metav1.TypeMeta{
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Kind:       v1alpha1.ResourceKindResourceDescriptor,
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
				Spec: v1alpha1.ResourceDescriptorSpec{
					Resource: v1alpha1.ResourceID{
						Group:   crd.Spec.Group,
						Version: version,
						Name:    plural,
						Kind:    kind,
						Scope:   v1alpha1.ResourceScope(crd.Spec.Scope),
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
		addResourceRequirements(&rd)

		if rd.Spec.Validation != nil &&
			rd.Spec.Validation.OpenAPIV3Schema != nil {

			var mc crdv1.JSONSchemaProps
			err = yaml.Unmarshal([]byte(v1alpha1.ObjectMetaSchema), &mc)
			if err != nil {
				return err
			}
			if rd.Spec.Resource.Scope == v1alpha1.ClusterScoped {
				delete(mc.Properties, "namespace")
			}
			rd.Spec.Validation.OpenAPIV3Schema.Properties["metadata"] = mc
			delete(rd.Spec.Validation.OpenAPIV3Schema.Properties, "status")
		}

		data, err := yaml.Marshal(rd)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func addResourceRequirements(rd *v1alpha1.ResourceDescriptor) {
	rd.Spec.ResourceRequirements = []v1alpha1.ResourceRequirements{defaultResourcePath()}
	if rd.Spec.Resource.Group == "kubedb.com" {
		if rd.Spec.Resource.Kind == "Elasticsearch" {
			topologies := []string{"master", "data", "ingest"}
			for _, topology := range topologies {
				rd.Spec.ResourceRequirements = append(rd.Spec.ResourceRequirements, v1alpha1.ResourceRequirements{
					Units:     fmt.Sprintf("spec.topology.%s.replicas", topology),
					Resources: fmt.Sprintf("spec.topology.%s.resources", topology),
				})
			}
		} else if rd.Spec.Resource.Kind == "MongoDB" {
			topologies := []string{"shard", "configServer", "mongos"}
			for _, topology := range topologies {
				rd.Spec.ResourceRequirements = append(rd.Spec.ResourceRequirements, v1alpha1.ResourceRequirements{
					Units:     fmt.Sprintf("spec.shardTopology.%s.shards", topology),
					Shards:    fmt.Sprintf("spec.shardTopology.%s.replicas", topology),
					Resources: fmt.Sprintf("spec.shardTopology.%s.podTemplate.spec.resources", topology),
				})
			}
		}
	}
}

func defaultResourcePath() v1alpha1.ResourceRequirements {
	return v1alpha1.ResourceRequirements{
		Units:     "spec.replicas",
		Resources: "spec.podTemplate.spec.resources",
	}
}
