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
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

/*
go run cmd/import-crds/main.go --input=/home/tamal/go/src/github.com/coreos/prometheus-operator/example/prometheus-operator-crd
go run cmd/import-crds/main.go --input=/home/tamal/go/src/github.com/jetstack/cert-manager/deploy/charts/cert-manager/crds
go run cmd/import-crds/main.go --input=/home/tamal/go/src/github.com/appscode/voyager/api/crds
go run cmd/import-crds/main.go --input=/home/tamal/go/src/stash.appscode.dev/apimachinery/crds
go run cmd/import-crds/main.go --input=/home/tamal/go/src/kmodules.xyz/custom-resources/crds
go run cmd/import-crds/main.go --input=/home/tamal/go/src/kubedb.dev/apimachinery/crds
go run cmd/import-crds/main.go --input=/home/tamal/go/src/kubevault.dev/operator/api/crds
go run cmd/import-crds/main.go --input=/home/tamal/go/src/go.searchlight.dev/grafana-operator/crds
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

func CustomResourceDefinition(data []byte) (*apiextensions.CustomResourceDefinition, error) {
	var out apiextensions.CustomResourceDefinition
	err := yaml.Unmarshal(data, &out)
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
		err = WriteDesciptor(crd, dir)
		if err != nil {
			return err
		}
	} else {
		fi, err := os.Stat(location)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			err = filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				if !strings.HasSuffix(info.Name(), "yaml") {
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
				err = WriteDesciptor(crd, dir)
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
			err = WriteDesciptor(crd, dir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func WriteDesciptor(crd *apiextensions.CustomResourceDefinition, dir string) error {
	version := crd.Spec.Version
	if len(crd.Spec.Versions) > 0 {
		version = crd.Spec.Versions[0].Name
	}

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
					Scope:   v1alpha1.ResourceScope(string(crd.Spec.Scope)),
				},
				Validation: crd.Spec.Validation,
			},
		}
	} else {
		err = yaml.Unmarshal(existing, &rd)
		if err == nil {
			rd.Spec.Validation = crd.Spec.Validation
		}
	}

	data, err := yaml.Marshal(rd)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
