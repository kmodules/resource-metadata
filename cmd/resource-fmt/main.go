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
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"

	metaapi "kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	uiapi "kmodules.xyz/resource-metadata/apis/ui/v1alpha1"

	"github.com/pkg/errors"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"gomodules.xyz/encoding/json"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

func handlePanic(filename string) {
	a := recover()
	if a != nil {
		fmt.Println("RECOVER", filename, a)
	}
}

func check(typ reflect.Type, filename string, fix bool) (string, error) {
	defer handlePanic(filename)

	data, err := os.ReadFile(filename)
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

	obj := reflect.New(typ).Interface()

	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		return "", err
	}
	parsed, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	// Then, Check them
	differ := diff.New()
	d, err := differ.Compare(sorted, parsed)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal file: %s", filename)
	}

	if fix {
		if d.Modified() {
			fmt.Println("formatted ", filename)
		}
		if f, ok := obj.(metaapi.YAMLFormatter); ok {
			data, err = f.ToYAML()
		} else {
			data, err = yaml.Marshal(obj)
		}
		if err != nil {
			return "", err
		}

		err = os.WriteFile(filename, data, 0o644)
		if err != nil {
			return "", err
		}
	} else if d.Modified() {
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
	}
	return "", nil
}

func checkType(t interface{}, plural string, fix bool) error {
	return filepath.Walk(fmt.Sprintf("./hub/%s/", plural), func(path string, info fs.FileInfo, err error) error {
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

		d, err := check(reflect.TypeOf(t), path, fix)
		if err != nil {
			return errors.Wrapf(err, "path=%s", path)
		}
		if d != "" {
			return fmt.Errorf("parsing diff found in file %s: %s", path, d)
		}
		return nil
	})
}

func MustCheckType(t interface{}, plural string, fix bool) {
	if err := checkType(t, plural, fix); err != nil {
		klog.ErrorS(err, "failed to check "+plural)
	}
}

func main() {
	fix := flag.Bool("fix", true, "Fix formatting")
	flag.Parse()

	MustCheckType(metaapi.ClusterProfile{}, "clusterprofiles", *fix)
	MustCheckType(metaapi.MenuOutline{}, "menuoutlines", *fix)
	MustCheckType(metaapi.ResourceBlockDefinition{}, "resourceblockdefinitions", *fix)
	MustCheckType(metaapi.ResourceDescriptor{}, "resourcedescriptors", *fix)
	MustCheckType(uiapi.ResourceEditor{}, "resourceeditors", *fix)
	MustCheckType(metaapi.ResourceOutline{}, "resourceoutlines", *fix)
	MustCheckType(metaapi.ResourceTableDefinition{}, "resourcetabledefinitions", *fix)
	MustCheckType(uiapi.ResourceDashboard{}, "resourcedashboards", *fix)
}
