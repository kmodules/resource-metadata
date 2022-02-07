/*
Copyright AppsCode Inc. and Contributors.

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

package v1alpha1

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetGenericResourceName(item client.Object) string {
	return fmt.Sprintf("%s~%s", item.GetName(), item.GetObjectKind().GroupVersionKind().GroupKind())
}

func ParseGenericResourceName(name string) (string, schema.GroupKind, error) {
	parts := strings.SplitN(name, "~", 2)
	if len(parts) != 2 {
		return "", schema.GroupKind{}, fmt.Errorf("expected resource name %s in format {.metadata.name}~Kind.Group", name)
	}
	return parts[0], schema.ParseGroupKind(parts[1]), nil
}
