package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	hub "kmodules.xyz/resource-metadata/hub/v1alpha1"
	"kmodules.xyz/resource-metadata/pkg/tableconvertor"
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

	gvr := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	rd, err := hub.LoadByGVR(gvr)
	if err != nil {
		log.Fatalln(err)
	}

	dep, err := dc.Resource(gvr).Namespace("default").Get("busy-dep", metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	c, err := tableconvertor.New(rd.Spec.DisplayColumns)
	if err != nil {
		log.Fatalln(err)
	}
	t, err := c.ConvertToTable(context.Background(), dep, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(t)
}
