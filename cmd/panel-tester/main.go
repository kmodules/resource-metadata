package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"kmodules.xyz/resource-metadata/hub"

	"k8s.io/client-go/kubernetes"
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

	kc, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	info, err := kc.Discovery().ServerVersion()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v", info)

	reg := hub.NewRegistryOfKnownResources()
	panel, err := reg.DefaultResourcePanel()
	if err != nil {
		log.Fatalln(err)
	}
	data, err := json.MarshalIndent(panel, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(data))
}
