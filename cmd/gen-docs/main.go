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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"kmodules.xyz/client-go/tools/parser"
	"kmodules.xyz/resource-metadata/hub/resourcedescriptors"
	"kmodules.xyz/resource-metadata/hub/resourceeditors"

	"github.com/pkg/errors"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func main() {
	dir := "/Users/tamal/go/src/go.crdhub.dev/crdhub/public"
	os.MkdirAll(dir, 0o755)
	GenDescriptor(dir)
	GenResourceEditor(dir)
	ImportCRD(dir)
}

func GenDescriptor(dir string) error {
	for _, rd := range resourcedescriptors.List() {
		g := rd.Spec.Resource.Group
		if g == "" {
			g = "core"
		}
		r := rd.Spec.Resource.Name

		rdir := filepath.Join(dir, g, r)
		os.MkdirAll(rdir, 0o755)

		vdir := filepath.Join(rdir, rd.Spec.Resource.Version)
		os.MkdirAll(vdir, 0o755)

		data, err := yaml.Marshal(rd)
		if err != nil {
			return err
		}
		os.WriteFile(filepath.Join(vdir, "resourcedescriptor.yaml"), data, 0o644)

		if rd.Spec.Validation != nil && rd.Spec.Validation.OpenAPIV3Schema != nil {
			schema := rd.Spec.Validation.OpenAPIV3Schema
			if prop, ok := schema.Properties["apiVersion"]; ok {
				prop.Enum = []crdv1.JSON{
					{[]byte(fmt.Sprintf("%q", rd.Spec.Resource.GroupVersion()))},
				}
				schema.Properties["apiVersion"] = prop
			}
			if prop, ok := schema.Properties["kind"]; ok {
				prop.Enum = []crdv1.JSON{
					{[]byte(fmt.Sprintf("%q", rd.Spec.Resource.Kind))},
				}
				schema.Properties["kind"] = prop
			}

			data, err := yaml.Marshal(schema)
			if err != nil {
				return err
			}
			os.WriteFile(filepath.Join(vdir, "openapiv3_schema.yaml"), data, 0o644)
		}
	}
	return nil
}

func GenResourceEditor(dir string) error {
	for _, rd := range resourceeditors.List() {
		g := rd.Spec.Resource.Group
		if g == "" {
			g = "core"
		}
		r := rd.Spec.Resource.Name

		rdir := filepath.Join(dir, g, r)
		os.MkdirAll(rdir, 0o755)

		// os.WriteFile(filepath.Join(rdir, "crd.yaml"), []byte("crd.yaml"), 0644)

		vdir := filepath.Join(rdir, rd.Spec.Resource.Version)
		os.MkdirAll(vdir, 0o755)

		data, err := yaml.Marshal(rd)
		if err != nil {
			return err
		}
		os.WriteFile(filepath.Join(vdir, "resourceeditor.yaml"), data, 0o644)
	}
	return nil
}

func ImportUIWizards() {
}

/*
# Add icons
go run cmd/icon-namer/main.go

# Import crds
# FluxCD CRDs
# helm template flux fluxcd-community/flux2 --output-dir=/tmp/fluxcd-manifests
/tmp/fluxcd-manifests
*/
func ImportCRD(dir string) error {
	_, f2, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("can't detect crd_list.txt file location")
	}
	fmt.Println(f2)

	data, err := os.ReadFile(filepath.Join(filepath.Dir(f2), "crd_list.txt"))
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = os.ExpandEnv(line)
		processLocation(dir, line)
	}
	return nil
}

func processLocation(dir, location string) error {
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
		return parser.ProcessResources(buf.Bytes(), processObject(dir))
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
				return parser.ProcessResources(data, processObject(dir))
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
			err = WriteDescriptor(dir, crd)
			if err != nil {
				return err
			}
		}
	}
	return nil
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

func processObject(dir string) func(ri parser.ResourceInfo) error {
	return func(ri parser.ResourceInfo) error {
		data, err := json.Marshal(ri.Object)
		if err != nil {
			return err
		}
		crd, err := CustomResourceDefinition(data)
		if err != nil {
			return err
		}
		err = WriteDescriptor(dir, crd)
		if err != nil {
			return err
		}
		return nil
	}
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

func WriteDescriptor(dir string, crd *crdv1.CustomResourceDefinition) error {
	rdir := filepath.Join(dir, crd.Spec.Group, crd.Spec.Names.Plural)
	os.MkdirAll(rdir, 0o755)
	data, err := yaml.Marshal(crd)
	if err != nil {
		return err
	}
	os.WriteFile(filepath.Join(rdir, "crd.yaml"), data, 0o644)
	return nil
}
