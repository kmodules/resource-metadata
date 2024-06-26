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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
	clientset "kmodules.xyz/resource-metadata/client/clientset/versioned"
	corev1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/core/v1alpha1"
	fakecorev1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/core/v1alpha1/fake"
	identityv1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/identity/v1alpha1"
	fakeidentityv1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/identity/v1alpha1/fake"
	managementv1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/management/v1alpha1"
	fakemanagementv1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/management/v1alpha1/fake"
	metav1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/meta/v1alpha1"
	fakemetav1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/meta/v1alpha1/fake"
	nodev1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/node/v1alpha1"
	fakenodev1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/node/v1alpha1/fake"
	uiv1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/ui/v1alpha1"
	fakeuiv1alpha1 "kmodules.xyz/resource-metadata/client/clientset/versioned/typed/ui/v1alpha1/fake"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	cs := &Clientset{tracker: o}
	cs.discovery = &fakediscovery.FakeDiscovery{Fake: &cs.Fake}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		return true, watch, nil
	})

	return cs
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
	tracker   testing.ObjectTracker
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *Clientset) Tracker() testing.ObjectTracker {
	return c.tracker
}

var (
	_ clientset.Interface = &Clientset{}
	_ testing.FakeClient  = &Clientset{}
)

// CoreV1alpha1 retrieves the CoreV1alpha1Client
func (c *Clientset) CoreV1alpha1() corev1alpha1.CoreV1alpha1Interface {
	return &fakecorev1alpha1.FakeCoreV1alpha1{Fake: &c.Fake}
}

// IdentityV1alpha1 retrieves the IdentityV1alpha1Client
func (c *Clientset) IdentityV1alpha1() identityv1alpha1.IdentityV1alpha1Interface {
	return &fakeidentityv1alpha1.FakeIdentityV1alpha1{Fake: &c.Fake}
}

// ManagementV1alpha1 retrieves the ManagementV1alpha1Client
func (c *Clientset) ManagementV1alpha1() managementv1alpha1.ManagementV1alpha1Interface {
	return &fakemanagementv1alpha1.FakeManagementV1alpha1{Fake: &c.Fake}
}

// MetaV1alpha1 retrieves the MetaV1alpha1Client
func (c *Clientset) MetaV1alpha1() metav1alpha1.MetaV1alpha1Interface {
	return &fakemetav1alpha1.FakeMetaV1alpha1{Fake: &c.Fake}
}

// NodeV1alpha1 retrieves the NodeV1alpha1Client
func (c *Clientset) NodeV1alpha1() nodev1alpha1.NodeV1alpha1Interface {
	return &fakenodev1alpha1.FakeNodeV1alpha1{Fake: &c.Fake}
}

// UiV1alpha1 retrieves the UiV1alpha1Client
func (c *Clientset) UiV1alpha1() uiv1alpha1.UiV1alpha1Interface {
	return &fakeuiv1alpha1.FakeUiV1alpha1{Fake: &c.Fake}
}
