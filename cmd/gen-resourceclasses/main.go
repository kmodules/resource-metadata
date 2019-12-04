/*
Copyright The Kmodules Authors.

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

	err = createRegistry(kc, filepath.Join("/home/tamal/go/src/kmodules.xyz/resource-metadata", "hub", "resourceclasses"))
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
					APIGroup: gv.Group,
				},
			}
			categories[name] = rd
		}

		for i := range rsList.APIResources {
			rs := rsList.APIResources[i]
			rd.Spec.Resources = append(rd.Spec.Resources, v1alpha1.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: rs.Name,
			})
		}
	}

	for _, rd := range categories {
		sort.Slice(rd.Spec.Resources, func(i, j int) bool { return rd.Spec.Resources[i].Resource < rd.Spec.Resources[j].Resource })

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
