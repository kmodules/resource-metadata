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
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"kmodules.xyz/resource-metadata/hub"

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

	reg := hub.NewRegistryOfKnownResources()
	err = reg.DiscoverResources(config)
	if err != nil {
		log.Fatalln(err)
	}

	{
		panel, err := reg.DefaultResourcePanel()
		if err != nil {
			log.Fatalln(err)
		}
		data, err := json.MarshalIndent(panel, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		err = ioutil.WriteFile("hub/defaultpanel.json", data, 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
	{
		panel, err := reg.CompleteResourcePanel()
		if err != nil {
			log.Fatalln(err)
		}
		data, err := json.MarshalIndent(panel, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		err = ioutil.WriteFile("hub/completepanel.json", data, 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
