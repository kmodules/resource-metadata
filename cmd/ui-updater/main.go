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
	"os"
	"path/filepath"

	kmapi "kmodules.xyz/client-go/api/v1"
	uiapi "kmodules.xyz/resource-metadata/apis/ui/v1alpha1"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	flag "github.com/spf13/pflag"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"gomodules.xyz/encoding/json"
	"sigs.k8s.io/yaml"
)

/*

```console
helm repo add appscode-charts-oci https://bundles.byte.builders/ui/
helm repo update
```

## Configure Development Helm Chart Repository

```console
helm repo add appscode-charts-oci-dev https://raw.githubusercontent.com/bytebuilders/ui-wizards/master/stable
helm repo update
```

*/

const (
	ociURL = "oci://ghcr.io/appscode-charts"
)

var (
	chartRegistryURL = flag.String("chart.registry-url", ociURL, "Chart registry url (prod/dev)")
	chartVersion     = flag.String("chart.version", "v0.12.0", "Chart version")
	useDigest        = flag.Bool("use-digest", true, "Use digest instead of tag")
)

var helmRepositories = map[string]string{
	"https://charts.appscode.com/stable/": "appscode-charts-legacy",
	ociURL:                                "appscode-charts-oci",
}

func check(filename string) (string, error) {
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

	var rd uiapi.ResourceEditor
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
			repoName := helmRepositories[*chartRegistryURL]
			if rd.Spec.UI.Options != nil {
				rd.Spec.UI.Options.SourceRef = kmapi.TypedObjectReference{
					APIGroup:  "source.toolkit.fluxcd.io",
					Kind:      "HelmRepository",
					Namespace: "",
					Name:      repoName,
				}
				rd.Spec.UI.Options.Version = getDigestOrVersion(repoName, rd.Spec.UI.Options.Name, *chartVersion)
			}
			if rd.Spec.UI.Editor != nil {
				rd.Spec.UI.Editor.SourceRef = kmapi.TypedObjectReference{
					APIGroup:  "source.toolkit.fluxcd.io",
					Kind:      "HelmRepository",
					Namespace: "",
					Name:      repoName,
				}
				rd.Spec.UI.Editor.Version = getDigestOrVersion(repoName, rd.Spec.UI.Editor.Name, *chartVersion)
			}
			for i, ag := range rd.Spec.UI.Actions {
				for j, a := range ag.Items {
					a.Editor.Version = getDigestOrVersion(a.Editor.SourceRef.Name, a.Editor.Name, *chartVersion)
					ag.Items[j] = a
				}
				rd.Spec.UI.Actions[i] = ag
			}
		}

		data, err := yaml.Marshal(rd)
		if err != nil {
			return "", err
		}

		err = os.WriteFile(filename, data, 0o644)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func main() {
	flag.Parse()

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

func getDigestOrVersion(repo, bin, ver string) string {
	if !*useDigest {
		return ver
	}
	if repo != "appscode-charts-oci" {
		return ver
	}
	digest, err := crane.Digest(fmt.Sprintf("ghcr.io/appscode-charts/%s:%s", bin, ver), crane.WithAuthFromKeychain(authn.DefaultKeychain))
	if err == nil {
		return digest
	}
	return ver
}
