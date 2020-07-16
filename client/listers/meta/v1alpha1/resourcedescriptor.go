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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "kmodules.xyz/resource-metadata/apis/meta/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ResourceDescriptorLister helps list ResourceDescriptors.
type ResourceDescriptorLister interface {
	// List lists all ResourceDescriptors in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.ResourceDescriptor, err error)
	// Get retrieves the ResourceDescriptor from the index for a given name.
	Get(name string) (*v1alpha1.ResourceDescriptor, error)
	ResourceDescriptorListerExpansion
}

// resourceDescriptorLister implements the ResourceDescriptorLister interface.
type resourceDescriptorLister struct {
	indexer cache.Indexer
}

// NewResourceDescriptorLister returns a new ResourceDescriptorLister.
func NewResourceDescriptorLister(indexer cache.Indexer) ResourceDescriptorLister {
	return &resourceDescriptorLister{indexer: indexer}
}

// List lists all ResourceDescriptors in the indexer.
func (s *resourceDescriptorLister) List(selector labels.Selector) (ret []*v1alpha1.ResourceDescriptor, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ResourceDescriptor))
	})
	return ret, err
}

// Get retrieves the ResourceDescriptor from the index for a given name.
func (s *resourceDescriptorLister) Get(name string) (*v1alpha1.ResourceDescriptor, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("resourcedescriptor"), name)
	}
	return obj.(*v1alpha1.ResourceDescriptor), nil
}
