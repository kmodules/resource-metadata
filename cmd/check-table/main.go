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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"kmodules.xyz/resource-metadata/hub"
	"kmodules.xyz/resource-metadata/pkg/tableconvertor"

	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := crd_cs.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	r := hub.NewRegistryOfKnownResources()

	{
		list, err := dc.Resource(gvr).Namespace("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalln(err)
		}

		t, err := tableconvertor.TableForList(r, client.CustomResourceDefinitions(), gvr, list.Items)
		if err != nil {
			log.Fatalln(err)
		}

		data, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(data))
	}

	{
		dep, err := dc.Resource(gvr).Namespace("default").Get(context.TODO(), "busy-dep", metav1.GetOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(dep.GroupVersionKind().String())

		t, err := tableconvertor.TableForObject(r, client.CustomResourceDefinitions(), dep)
		if err != nil {
			log.Fatalln(err)
		}

		data, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(data))
	}
}
