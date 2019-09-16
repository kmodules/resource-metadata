package v1alpha1

import (
	"fmt"
	"strings"

	kerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kmodules.xyz/resource-metadata/apis/meta"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"sigs.k8s.io/yaml"
)

func LoadByGVR(gvr schema.GroupVersionResource) (*v1alpha1.ResourceDescriptor, error) {
	var filename string
	if gvr.Group == "" && gvr.Version == "v1" {
		filename = fmt.Sprintf("core/v1/%s.yaml", gvr.Resource)
	} else {
		filename = fmt.Sprintf("%s/%s/%s.yaml", gvr.Group, gvr.Version, gvr.Resource)
	}
	return LoadByFile(filename)
}

func LoadByName(name string) (*v1alpha1.ResourceDescriptor, error) {
	filename := strings.Replace(name, "-", "/", 2) + ".yaml"
	return LoadByFile(filename)
}

func LoadByFile(filename string) (*v1alpha1.ResourceDescriptor, error) {
	data, err := Asset(filename)
	if err != nil {
		return nil, kerr.NewNotFound(schema.GroupResource{Group: meta.GroupName, Resource: v1alpha1.ResourceKindResourceDescriptor}, strings.TrimSuffix(filename, ".yaml"))
	}
	var obj v1alpha1.ResourceDescriptor
	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		return nil, kerr.NewInternalError(err)
	}
	return &obj, nil
}
