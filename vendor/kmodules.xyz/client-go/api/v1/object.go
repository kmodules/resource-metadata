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

package v1

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ObjectReference contains enough information to let you inspect or modify the referred object.
type ObjectReference struct {
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

type ObjectID struct {
	Group     string `json:"group,omitempty" protobuf:"bytes,1,opt,name=group"`
	Kind      string `json:"kind,omitempty" protobuf:"bytes,2,opt,name=kind"`
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
	Name      string `json:"name,omitempty" protobuf:"bytes,4,opt,name=name"`
}

func (oid *ObjectID) Key() string {
	return fmt.Sprintf("G=%s,K=%s,NS=%s,N=%s", oid.Group, oid.Kind, oid.Namespace, oid.Name)
}

func NewObjectID(obj client.Object) *ObjectID {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return &ObjectID{
		Group:     gvk.Group,
		Kind:      gvk.Kind,
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
}

func ParseObjectID(key string) (*ObjectID, error) {
	parts := strings.FieldsFunc(key, func(r rune) bool {
		return r == ',' || r == '='
	})

	var id ObjectID
	for i := 0; i < len(parts); i += 2 {
		switch parts[i] {
		case "G":
			id.Group = parts[i+1]
		case "K":
			id.Kind = parts[i+1]
		case "NS":
			id.Namespace = parts[i+1]
		case "N":
			id.Name = parts[i+1]
		default:
			return nil, fmt.Errorf("unknown key %key", parts[i])
		}
	}
	if id.Group == "" {
		return nil, fmt.Errorf("group not set")
	}
	if id.Kind == "" {
		return nil, fmt.Errorf("kind not set")
	}
	if id.Name == "" {
		return nil, fmt.Errorf("name not set")
	}
	return &id, nil
}

func (oid *ObjectID) GroupKind() schema.GroupKind {
	return schema.GroupKind{Group: oid.Group, Kind: oid.Kind}
}

func (oid *ObjectID) MetaGroupKind() metav1.GroupKind {
	return metav1.GroupKind{Group: oid.Group, Kind: oid.Kind}
}

func (oid *ObjectID) ObjectReference() ObjectReference {
	return ObjectReference{Namespace: oid.Namespace, Name: oid.Name}
}

func (oid *ObjectID) ObjectKey() client.ObjectKey {
	return client.ObjectKey{Namespace: oid.Namespace, Name: oid.Name}
}
