package graph

import (
	"os"
	"path/filepath"
	"testing"

	"kmodules.xyz/resource-metadata/hub/v1alpha1"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func TestLoadGraph(t *testing.T) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	assert.NoError(t, err)
	var dc discovery.DiscoveryInterface
	dc, err = discovery.NewDiscoveryClientForConfig(config)
	assert.NoError(t, err)
	gvr := schema.GroupVersionResource{
		Group:    "monitoring.coreos.com",
		Version:  "v1",
		Resource: "alertmanagers",
	}
	assert.NoError(t, v1alpha1.Register(gvr, dc, config))
	graph, err := LoadGraph()
	assert.NoError(t, err)
	assert.NotNil(t, graph)
}
