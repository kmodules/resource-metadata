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
	"path/filepath"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	uiapi "kmodules.xyz/resource-metadata/apis/ui/v1alpha1"
	blockdefs "kmodules.xyz/resource-metadata/hub/resourceblockdefinitions"
	dashboards "kmodules.xyz/resource-metadata/hub/resourcedashboards"
	"kmodules.xyz/resource-metadata/hub/resourcedescriptors"
	"kmodules.xyz/resource-metadata/hub/resourceoutlines"
	tabledefs "kmodules.xyz/resource-metadata/hub/resourcetabledefinitions"
	sc "kmodules.xyz/schema-checker"
)

func main() {
	if err := checkYAMLs(blockdefs.EmbeddedFS(), &v1alpha1.ResourceBlockDefinition{}); err != nil {
		panic(err)
	}
	if err := checkYAMLs(resourcedescriptors.EmbeddedFS(), &v1alpha1.ResourceDescriptor{}); err != nil {
		panic(err)
	}
	if err := checkYAMLs(resourceoutlines.EmbeddedFS(), &v1alpha1.ResourceOutline{}); err != nil {
		panic(err)
	}
	if err := checkYAMLs(tabledefs.EmbeddedFS(), &v1alpha1.ResourceTableDefinition{}); err != nil {
		panic(err)
	}
	if err := checkYAMLs(dashboards.EmbeddedFS(), &uiapi.ResourceDashboard{}); err != nil {
		panic(err)
	}
}

// hub packages embed a non-YAML "trigger" sentinel file for hot reload,
// which sc.CheckFS would try to parse as an object.
func checkYAMLs(fsys fs.FS, v interface{}) error {
	return fs.WalkDir(fsys, ".", func(path string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() || filepath.Ext(path) != ".yaml" {
			return err
		}
		d, err := sc.New(fsys).CheckObject(v, path)
		if err != nil {
			return err
		}
		if d != "" {
			return fmt.Errorf("%s: object does not match schema, diff: %s", path, d)
		}
		return nil
	})
}
