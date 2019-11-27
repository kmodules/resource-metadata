package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	hub "kmodules.xyz/resource-metadata/hub/v1alpha1"
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

	r := hub.NewRegistryOfKnownTypes()

	{
		list, err := dc.Resource(gvr).Namespace("default").List(metav1.ListOptions{})
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
		dep, err := dc.Resource(gvr).Namespace("default").Get("busy-dep", metav1.GetOptions{})
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
