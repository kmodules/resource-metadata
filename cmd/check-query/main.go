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

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub"
)

func main() {
	reg := hub.NewRegistryOfKnownResources()
	reg.Visit(func(key string, rd *v1alpha1.ResourceDescriptor) {
		for i, p := range rd.Spec.Pages {
			for j, sec := range p.Resources {
				if sec.Query.Type != v1alpha1.RESTQuery && sec.Query.Type != v1alpha1.GraphQLQuery {
					panic(fmt.Errorf("key=%s rd.Spec.Pages[%d].Resources[%d]", key, i, j))
				}
			}
		}
	})
}
