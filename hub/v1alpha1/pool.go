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
	"sync"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	lru "github.com/hashicorp/golang-lru"
)

type Pool struct {
	cache *lru.Cache
	m     sync.Mutex
	f     func() KV
}

func NewPool(kvFactory func() KV) (*Pool, error) {
	cache, err := lru.New(1024)
	if err != nil {
		return nil, err
	}
	return &Pool{cache: cache, f: kvFactory}, nil
}

func NewLocalPool() (*Pool, error) {
	return NewPool(func() KV {
		return &KVLocal{
			shared: KnownResources,
			cache:  map[string]*v1alpha1.ResourceDescriptor{},
		}
	})
}

func (p *Pool) GetRegistry(uid string) *Registry {
	p.m.Lock()
	defer p.m.Unlock()

	val, found := p.cache.Get(uid)
	if found {
		return val.(*Registry)
	}
	r := NewRegistry(uid, p.f())
	p.cache.Add(uid, r)
	return r
}
