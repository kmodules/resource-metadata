package v1alpha1

import (
	"fmt"
	"strings"
	"sync"

	"kmodules.xyz/resource-metadata/apis/meta"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

var (
	regGVK = make(map[schema.GroupVersionKind]*v1alpha1.ResourceID, len(_bindata))
	regGVR = make(map[schema.GroupVersionResource]*v1alpha1.ResourceID, len(_bindata))

	cache       = make(map[string]*v1alpha1.ResourceDescriptor)
	m           sync.RWMutex
	clusterHost string
)

func GetCachedResourceDescriptor() map[string]*v1alpha1.ResourceDescriptor {
	return cache
}

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

func Register(gvr schema.GroupVersionResource, dc discovery.ServerResourcesInterface, config *rest.Config) error {
	m.RLock()
	if _, found := regGVR[gvr]; found {
		m.RUnlock()
		return nil
	}
	m.RUnlock()

	reg, err := createRegistry(dc)
	if err != nil {
		return err
	}
	clusterHost = config.Host

	m.Lock()
	for filename, rd := range reg {
		if _, found := cache[cachedFileName(filename)]; !found {
			regGVK[rd.Spec.Resource.GroupVersionKind()] = &rd.Spec.Resource
			regGVR[rd.Spec.Resource.GroupVersionResource()] = &rd.Spec.Resource
			cache[cachedFileName(filename)] = rd
		}
	}
	m.Unlock()

	return nil
}

func createRegistry(dc discovery.ServerResourcesInterface) (map[string]*v1alpha1.ResourceDescriptor, error) {
	rsLists, err := dc.ServerPreferredResources()
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return nil, err
	}

	reg := make(map[string]*v1alpha1.ResourceDescriptor)
	for _, rsList := range rsLists {
		for i := range rsList.APIResources {
			rs := rsList.APIResources[i]

			gv, err := schema.ParseGroupVersion(rsList.GroupVersion)
			if err != nil {
				return nil, err
			}
			rs.Group = gv.Group
			rs.Version = gv.Version

			scope := v1alpha1.ClusterScoped
			if rs.Namespaced {
				scope = v1alpha1.NamespaceScoped
			}

			filename := fmt.Sprintf("%s/%s/%s.yaml", rs.Group, rs.Version, rs.Name)
			reg[filename] = &v1alpha1.ResourceDescriptor{
				TypeMeta: metav1.TypeMeta{
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Kind:       v1alpha1.ResourceKindResourceDescriptor,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s-%s", rs.Group, rs.Version, rs.Name),
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
		}
	}

	for _, name := range AssetNames() {
		delete(reg, name)
	}
	return reg, nil
}

func GVR(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	m.RLock()
	defer m.RUnlock()
	rid, exist := regGVK[gvk]
	if !exist {
		return schema.GroupVersionResource{}, fmt.Errorf("gvk %v isn't registered", gvk)
	}
	return rid.GroupVersionResource(), nil
}

func TypeMeta(gvr schema.GroupVersionResource) (metav1.TypeMeta, error) {
	m.RLock()
	defer m.RUnlock()
	rid, exist := regGVR[gvr]
	if !exist {
		return metav1.TypeMeta{}, fmt.Errorf("gvr %v isn't registered", gvr)
	}
	return rid.TypeMeta(), nil
}

func GVK(gvr schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	m.RLock()
	defer m.RUnlock()
	rid, exist := regGVR[gvr]
	if !exist {
		return schema.GroupVersionKind{}, fmt.Errorf("gvr %v isn't registered", gvr)
	}
	return rid.GroupVersionKind(), nil
}

func IsNamespaced(gvr schema.GroupVersionResource) (bool, error) {
	m.RLock()
	defer m.RUnlock()
	rid, exist := regGVR[gvr]
	if !exist {
		return false, fmt.Errorf("gvr %v isn't registered", gvr)
	}
	return rid.Scope == v1alpha1.NamespaceScoped, nil
}

func Types() []metav1.TypeMeta {
	m.RLock()
	defer m.RUnlock()

	types := make([]metav1.TypeMeta, 0, len(regGVK))
	for _, v := range regGVK {
		types = append(types, v.TypeMeta())
	}
	return types
}

func Resources() []schema.GroupVersionResource {
	m.RLock()
	defer m.RUnlock()

	resources := make([]schema.GroupVersionResource, 0, len(regGVR))
	for k := range regGVR {
		resources = append(resources, k)
	}
	return resources
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
	if data, err := Asset(filename); err == nil {
		var obj v1alpha1.ResourceDescriptor
		err = yaml.Unmarshal(data, &obj)
		if err != nil {
			return nil, kerr.NewInternalError(err)
		}
		return &obj, nil
	}

	m.RLock()
	defer m.RUnlock()

	obj, ok := cache[cachedFileName(filename)]
	if !ok {
		return nil, kerr.NewNotFound(schema.GroupResource{Group: meta.GroupName, Resource: v1alpha1.ResourceKindResourceDescriptor}, strings.TrimSuffix(filename, ".yaml"))
	}
	return obj, nil
}

func cachedFileName(filename string) string {
	return clusterHost + "/" + filename
}
