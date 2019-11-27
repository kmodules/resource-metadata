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

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

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

	err = createRegistry(kc, filepath.Join("hub", "v1alpha1"))
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

	for _, rsList := range rsLists {
		fmt.Println(rsList.GroupVersion)
		for i := range rsList.APIResources {
			rs := rsList.APIResources[i]

			gv, err := schema.ParseGroupVersion(rsList.GroupVersion)
			if err != nil {
				return err
			}
			rs.Group = gv.Group
			rs.Version = gv.Version

			scope := v1alpha1.ClusterScoped
			if rs.Namespaced {
				scope = v1alpha1.NamespaceScoped
			}

			name := fmt.Sprintf("%s-%s-%s", rs.Group, rs.Version, rs.Name)
			baseDir := filepath.Join(dir, rs.Group, rs.Version)
			if rs.Group == "" {
				name = fmt.Sprintf("core-%s-%s", rs.Version, rs.Name)
				baseDir = filepath.Join(dir, "core", rs.Version)
			}

			rd := v1alpha1.ResourceDescriptor{
				TypeMeta: metav1.TypeMeta{
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Kind:       v1alpha1.ResourceKindResourceDescriptor,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
					Labels: map[string]string{
						"k8s.io/group":    rs.Group,
						"k8s.io/version":  rs.Version,
						"k8s.io/resource": rs.Name,
						"k8s.io/kind":     rs.Kind,
					},
				},
				Spec: v1alpha1.ResourceDescriptorSpec{
					Resource: v1alpha1.ResourceID{
						Group:   rs.Group,
						Version: rs.Version,
						Name:    rs.Name,
						Kind:    rs.Kind,
						Scope:   scope,
					},
				},
			}
			data, err := yaml.Marshal(rd)
			if err != nil {
				return err
			}

			err = os.MkdirAll(baseDir, 0755)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(filepath.Join(baseDir, rs.Name+".yaml"), data, 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
