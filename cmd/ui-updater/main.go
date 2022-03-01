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

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	flag "github.com/spf13/pflag"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"gomodules.xyz/encoding/json"
	"sigs.k8s.io/yaml"
)

/*

```console
helm repo add bytebuilders-ui https://bundles.byte.builders/ui/
helm repo update
```

## Configure Development Helm Chart Repository

```console
helm repo add bytebuilders-ui-dev https://raw.githubusercontent.com/bytebuilders/ui-wizards/master/stable
helm repo update
```

*/

const (
	prodURL = "https://bundles.byte.builders/ui/"
	devURL  = "https://raw.githubusercontent.com/bytebuilders/ui-wizards/master/stable"
)

var (
	chartRegistryURL = flag.String("chart.registry-url", prodURL, "Chart registry url (prod/dev)")
	chartVersion     = flag.String("chart.version", "v0.3.0", "Chart version")
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

	var rd v1alpha1.ResourceEditor
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
		if rd.Spec.UI != nil {
			if rd.Spec.UI.Options != nil {
				rd.Spec.UI.Options.URL = *chartRegistryURL
				rd.Spec.UI.Options.Version = *chartVersion
			}
			if rd.Spec.UI.Editor != nil {
				rd.Spec.UI.Editor.URL = *chartRegistryURL
				rd.Spec.UI.Editor.Version = *chartVersion
			}
		}

		data, err := yaml.Marshal(rd)
		if err != nil {
			return "", err
		}

		err = ioutil.WriteFile(filename, data, 0o644)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func main() {
	flag.Parse()
	switch *chartRegistryURL {
	case "prod":
		*chartRegistryURL = prodURL
	case "dev", "qa":
		*chartRegistryURL = devURL
	}

	err := filepath.Walk("./hub/resourceeditors/", func(path string, info fs.FileInfo, err error) error {
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
			return err
		}
		if d != "" {
			return fmt.Errorf("parsing error in file %s: %s", path, d)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
