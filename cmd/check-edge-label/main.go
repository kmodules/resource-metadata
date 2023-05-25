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
	"os"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub"
)

func main() {
	err := 0
	reg := hub.NewRegistryOfKnownResources()
	reg.Visit(func(key string, rdSrc *v1alpha1.ResourceDescriptor) {
		for _, c := range rdSrc.Spec.Connections {
			rdTarget, _ := reg.LoadByGVK(c.Target.GroupVersionKind())

			var offshoot bool
			for _, lbl := range c.Labels {
				if lbl == kmapi.EdgeLabelOffshoot {
					offshoot = true
					break
				}
			}
			if offshoot {
				if rdSrc.Spec.Resource.Scope != rdTarget.Spec.Resource.Scope {
					fmt.Printf("%+v has an offshoot label edge with %+v, but their scope does not match\n", rdSrc.Spec.Resource.GroupVersionKind(), rdTarget.Spec.Resource.GroupVersionKind())
					err++
				}
			}
		}
	})
	os.Exit(err)
}
