package v1alpha1

import (
	"fmt"
	"strings"

	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kmodules.xyz/resource-metadata/apis/meta"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"sigs.k8s.io/yaml"
)

var (
	regGVK = make(map[schema.GroupVersionKind]*v1alpha1.ResourceID, len(_bindata))
	regGVR = make(map[schema.GroupVersionResource]*v1alpha1.ResourceID, len(_bindata))
)

func init() {
	for _, filename := range AssetNames() {
		rd, err := LoadByFile(filename)
		if err != nil {
			panic(err)
		}
		v := rd.Spec.Resource // copy
		regGVK[v.GroupVersionKind()] = &v
		regGVR[v.GroupVersionResource()] = &v
	}
}

func GVR(t metav1.TypeMeta) schema.GroupVersionResource {
	// handle not found
	return regGVK[t.GroupVersionKind()].GroupVersionResource()
}

func TypeMeta(gvr schema.GroupVersionResource) metav1.TypeMeta {
	// handle not found
	return regGVR[gvr].TypeMeta()
}

func GVK(gvr schema.GroupVersionResource) schema.GroupVersionKind {
	// handle not found
	return regGVR[gvr].GroupVersionKind()
}

func IsNamespaced(t metav1.TypeMeta) bool {
	return regGVK[t.GroupVersionKind()].Scope == v1alpha1.NamespaceScoped
}

func Types() []metav1.TypeMeta {
	types := make([]metav1.TypeMeta, 0, len(regGVK))
	for _, v := range regGVK {
		types = append(types, v.TypeMeta())
	}
	return types
}

func LoadByGVK(gvk schema.GroupVersionKind) (*v1alpha1.ResourceDescriptor, error) {
	id, ok := regGVK[gvk]
	if !ok {
		return nil, fmt.Errorf("unknown GVK:%v", gvk.String())
	}
	return LoadByGVR(id.GroupVersionResource())
}

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
