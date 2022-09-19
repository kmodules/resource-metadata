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

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
)

func main() {
	filename := "/home/tamal/go/src/kmodules.xyz/resource-metadata/hub/resourcedescriptors/kubedb.com/v1alpha2/elasticsearches.yaml"
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	data, err = v1alpha1.FormatMetadata(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
