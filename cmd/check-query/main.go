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

	"kmodules.xyz/resource-metadata/apis/shared"
	"kmodules.xyz/resource-metadata/hub/resourceblockdefinitions"
)

func main() {
	for _, rd := range resourceblockdefinitions.List() {
		for i, p := range rd.Spec.Blocks {
			if p.ResourceLocator != nil {
				sec := p.ResourceLocator
				if sec.Query.Type != shared.RESTQuery && sec.Query.Type != shared.GraphQLQuery {
					panic(fmt.Errorf("rd.Spec.Blocks[%d]", i))
				}
			}
		}
	}
}
