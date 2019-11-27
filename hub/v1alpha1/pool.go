package v1alpha1

import (
	lru "github.com/hashicorp/golang-lru"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"sync"
)

type Pool struct {
	cache *lru.Cache
	m sync.Mutex
	f func() KV
}

func NewPool(kvFactory func() KV) (*Pool, error) {
	cache, err := lru.New(1024)
	if err != nil {
		return nil, err
	}
	return &Pool{cache: cache, f: kvFactory}, nil
}

func NewLocalPool() (*Pool, error) {
	return NewPool(func () KV {
	return	&kvLocal{
			shared: known,
			cache: map[string]*v1alpha1.ResourceDescriptor{},
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