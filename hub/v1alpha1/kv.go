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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"sigs.k8s.io/yaml"
	"strings"
	"sync"
)

type KV interface {
	Set(key string, val *v1alpha1.ResourceDescriptor)
	Get(key string) (*v1alpha1.ResourceDescriptor, bool)
	Visit (func(key string, val *v1alpha1.ResourceDescriptor))
}

type kvMap struct {
	cache map[string]*v1alpha1.ResourceDescriptor
	m     sync.RWMutex
}

var _ KV = &kvMap{}

func (s *kvMap) Set(key string, val *v1alpha1.ResourceDescriptor) {
	s.m.Lock()
	s.cache[key] = val
	s.m.Unlock()
}

func (s *kvMap) Get(key string) (*v1alpha1.ResourceDescriptor, bool) {
	s.m.RLock()
	val, found := s.cache[key]
	s.m.RUnlock()
	return val, found
}

func (s *kvMap) Visit (f func(key string, val *v1alpha1.ResourceDescriptor)) {
	s.m.RLock()
	for k, v := range s.cache {
		f(k, v)
	}
	s.m.RUnlock()
}

type kvLocal struct {
	shared KV
	cache  map[string]*v1alpha1.ResourceDescriptor
}

var _ KV = &kvLocal{}

func (s *kvLocal) Set(key string, val *v1alpha1.ResourceDescriptor) {
	s.cache[key] = val
}

func (s *kvLocal) Get(key string) (*v1alpha1.ResourceDescriptor, bool) {
	val, found := s.shared.Get(key)
	if found {
		return val, found
	}
	val, found = s.cache[key]
	return val, found
}

func (s *kvLocal) Visit (f func(key string, val *v1alpha1.ResourceDescriptor)) {
	s.shared.Visit(f)
	for k, v := range s.cache {
		f(k, v)
	}
}

var (
	known KV = &kvMap{
		cache: make(map[string]*v1alpha1.ResourceDescriptor),
	}
)

func init() {
	for _, filename := range AssetNames() {
		rd, err := LoadByFile(filename)
		if err != nil {
			panic(err)
		}
		known.Set(filename, rd)
	}
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
		return nil, err
	}
	var obj v1alpha1.ResourceDescriptor
	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
