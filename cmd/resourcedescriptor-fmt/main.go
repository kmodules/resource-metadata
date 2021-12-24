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
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	"github.com/pkg/errors"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"gomodules.xyz/encoding/json"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

func check(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var original map[string]interface{}
	err = yaml.Unmarshal(data, &original)
	if err != nil {
		return "", err
	}
	sorted, err := json.Marshal(&original)
	if err != nil {
		return "", err
	}

	var rd v1alpha1.ResourceDescriptor
	err = yaml.Unmarshal(data, &rd)
	if err != nil {
		return "", err
	}
	parsed, err := json.Marshal(rd)
	if err != nil {
		return "", err
	}

	// Then, Check them
	differ := diff.New()
	d, err := differ.Compare(sorted, parsed)
	if err != nil {
		fmt.Printf("Failed to unmarshal file: %s\n", err.Error())
		os.Exit(3)
	}

	if d.Modified() {
		config := formatter.AsciiFormatterConfig{
			ShowArrayIndex: true,
			Coloring:       true,
		}

		f := formatter.NewAsciiFormatter(original, config)
		result, err := f.Format(d)
		if err != nil {
			return "", err
		}
		return result, nil
	} else {
		if rd.Spec.Validation != nil &&
			rd.Spec.Validation.OpenAPIV3Schema != nil {

			var mc crdv1.JSONSchemaProps
			err = yaml.Unmarshal([]byte(v1alpha1.ObjectMetaSchema), &mc)
			if err != nil {
				return "", err
			}
			if rd.Spec.Resource.Scope == kmapi.ClusterScoped {
				delete(mc.Properties, "namespace")
			}
			rd.Spec.Validation.OpenAPIV3Schema.Properties["metadata"] = mc
			delete(rd.Spec.Validation.OpenAPIV3Schema.Properties, "status")
		}

		data, err := yaml.Marshal(rd)
		if err != nil {
			return "", err
		}

		data, err = v1alpha1.FormatMetadata(data)
		if err != nil {
			return "", err
		}

		err = ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func main() {
	err := filepath.Walk("./hub/resourcedescriptors/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(info.Name())
		if ext != ".yml" && ext != ".yaml" && ext != ".json" {
			return nil
		}

		d, err := check(path)
		if err != nil {
			return errors.Wrapf(err, "path=%s", path)
		}
		if d != "" {
			return fmt.Errorf("parsing diff found in file %s: %s", path, d)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
