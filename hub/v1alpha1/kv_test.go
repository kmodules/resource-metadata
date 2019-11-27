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

package v1alpha1

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func TestRegister(t *testing.T) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	assert.NoError(t, err)
	var dc discovery.DiscoveryInterface
	dc, err = discovery.NewDiscoveryClientForConfig(config)
	fmt.Println("config.Host: ", config.Host)
	assert.NoError(t, err)
	gvr := schema.GroupVersionResource{
		Group:    "monitoring.coreos.com",
		Version:  "v1",
		Resource: "alertmanagers",
	}
	assert.NoError(t, Register(gvr, dc, config))
	rd1, err := LoadByGVR(gvr)
	assert.NoError(t, err)
	rd2, err := LoadByFile("monitoring.coreos.com/v1/alertmanagers.yaml")
	assert.NoError(t, err)
	assert.Equal(t, rd1, rd2)
}

func TestDiscovery(t *testing.T) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	assert.NoError(t, err)
	var dc discovery.DiscoveryInterface
	dc, err = discovery.NewDiscoveryClientForConfig(config)
	fmt.Println("config.Host: ", config.Host)
	assert.NoError(t, err)
	list, err := dc.ServerPreferredResources()
	assert.NoError(t, err)
	for _, ls := range list {
		for _, rs := range ls.APIResources {
			fmt.Println(ls.GroupVersion+"/"+rs.Name, ": ", rs.Verbs)
		}
	}
	lst, err := dc.ServerResourcesForGroupVersion("v1")
	assert.NoError(t, err)
	for _, rs := range lst.APIResources {
		fmt.Println(lst.GroupVersion+"/"+rs.Name, ": ", rs.Verbs)
	}
}
