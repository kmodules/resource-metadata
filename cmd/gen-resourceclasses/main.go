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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/hub/resourceclasses"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/yaml"
)

func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	kc, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	err = createRegistry(kc, filepath.Join("$HOME/go/src/kmodules.xyz/resource-metadata", "hub", "resourceclasses"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println()
	fmt.Println("resource hub generated.")
}

func createRegistry(kc kubernetes.Interface, dir string) error {
	rsLists, err := kc.Discovery().ServerPreferredResources()
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return err
	}

	categories := make(map[string]*v1alpha1.ResourceClass)

	for _, rsList := range rsLists {
		gv, err := schema.ParseGroupVersion(rsList.GroupVersion)
		if err != nil {
			return err
		}

		name := resourceclasses.ResourceClassName(gv.Group)
		fmt.Println(name + " | " + gv.Version)

		rd, found := categories[name]
		if !found {
			rd = &v1alpha1.ResourceClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Kind:       v1alpha1.ResourceKindResourceClass,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
				Spec: v1alpha1.ResourceClassSpec{
					ResourceClassInfo: v1alpha1.ResourceClassInfo{
						APIGroup: gv.Group,
					},
				},
			}
			categories[name] = rd
		}

		for i := range rsList.APIResources {
			rs := rsList.APIResources[i]
			rd.Spec.Entries = append(rd.Spec.Entries, v1alpha1.MenuEntry{
				Type: &metav1.GroupKind{
					Group: gv.Group,
					Kind:  rs.Kind,
				},
				Required: false,
			})
		}
	}

	for _, rd := range categories {
		sort.Slice(rd.Spec.Entries, func(i, j int) bool { return rd.Spec.Entries[i].Type.Kind < rd.Spec.Entries[j].Type.Kind })

		data, err := yaml.Marshal(rd)
		if err != nil {
			return err
		}

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(dir, strings.ToLower(rd.Name)+".yaml"), data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
