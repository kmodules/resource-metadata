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
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	blockdefs "kmodules.xyz/resource-metadata/hub/resourceblockdefinitions"
	dashboards "kmodules.xyz/resource-metadata/hub/resourcedashboards"
	"kmodules.xyz/resource-metadata/hub/resourcedescriptors"
	"kmodules.xyz/resource-metadata/hub/resourceoutlines"
	tabledefs "kmodules.xyz/resource-metadata/hub/resourcetabledefinitions"
	sc "kmodules.xyz/schema-checker"
)

func main() {
	if err := sc.CheckFS(blockdefs.EmbeddedFS(), &v1alpha1.ResourceBlockDefinition{}); err != nil {
		panic(err)
	}
	if err := sc.CheckFS(resourcedescriptors.EmbeddedFS(), &v1alpha1.ResourceDescriptor{}); err != nil {
		panic(err)
	}
	if err := sc.CheckFS(resourceoutlines.EmbeddedFS(), &v1alpha1.ResourceOutline{}); err != nil {
		panic(err)
	}
	if err := sc.CheckFS(tabledefs.EmbeddedFS(), &v1alpha1.ResourceTableDefinition{}); err != nil {
		panic(err)
	}
	if err := sc.CheckFS(dashboards.EmbeddedFS(), &v1alpha1.ResourceDashboard{}); err != nil {
		panic(err)
	}
}
