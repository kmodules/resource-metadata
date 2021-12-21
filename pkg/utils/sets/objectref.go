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

package sets

import (
	"reflect"
	"sort"

	apiv1 "kmodules.xyz/client-go/api/v1"
)

// sets.ObjectReference is a set of apiv1.ObjectReferences, implemented via map[apiv1.ObjectReference]struct{} for minimal memory consumption.
type ObjectReference map[apiv1.ObjectReference]Empty

// NewObjectReference creates a ObjectReference from a list of values.
func NewObjectReference(items ...apiv1.ObjectReference) ObjectReference {
	ss := make(ObjectReference, len(items))
	ss.Insert(items...)
	return ss
}

// ObjectReferenceKeySet creates a ObjectReference from a keys of a map[apiv1.ObjectReference](? extends interface{}).
// If the value passed in is not actually a map, this will panic.
func ObjectReferenceKeySet(theMap interface{}) ObjectReference {
	v := reflect.ValueOf(theMap)
	ret := ObjectReference{}

	for _, keyValue := range v.MapKeys() {
		ret.Insert(keyValue.Interface().(apiv1.ObjectReference))
	}
	return ret
}

// Insert adds items to the set.
func (s ObjectReference) Insert(items ...apiv1.ObjectReference) ObjectReference {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

// Delete removes all items from the set.
func (s ObjectReference) Delete(items ...apiv1.ObjectReference) ObjectReference {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has returns true if and only if item is contained in the set.
func (s ObjectReference) Has(item apiv1.ObjectReference) bool {
	_, contained := s[item]
	return contained
}

// HasAll returns true if and only if all items are contained in the set.
func (s ObjectReference) HasAll(items ...apiv1.ObjectReference) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any items are contained in the set.
func (s ObjectReference) HasAny(items ...apiv1.ObjectReference) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

// Difference returns a set of objects that are not in s2
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.Difference(s2) = {a3}
// s2.Difference(s1) = {a4, a5}
func (s ObjectReference) Difference(s2 ObjectReference) ObjectReference {
	result := NewObjectReference()
	for key := range s {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// Union returns a new set which includes items in either s1 or s2.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Union(s2) = {a1, a2, a3, a4}
// s2.Union(s1) = {a1, a2, a3, a4}
func (s1 ObjectReference) Union(s2 ObjectReference) ObjectReference {
	result := NewObjectReference()
	for key := range s1 {
		result.Insert(key)
	}
	for key := range s2 {
		result.Insert(key)
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}
func (s1 ObjectReference) Intersection(s2 ObjectReference) ObjectReference {
	var walk, other ObjectReference
	result := NewObjectReference()
	if s1.Len() < s2.Len() {
		walk = s1
		other = s2
	} else {
		walk = s2
		other = s1
	}
	for key := range walk {
		if other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s1 ObjectReference) IsSuperset(s2 ObjectReference) bool {
	for item := range s2 {
		if !s1.Has(item) {
			return false
		}
	}
	return true
}

// Equal returns true if and only if s1 is equal (as a set) to s2.
// Two sets are equal if their membership is identical.
// (In practice, this means same elements, order doesn't matter)
func (s1 ObjectReference) Equal(s2 ObjectReference) bool {
	return len(s1) == len(s2) && s1.IsSuperset(s2)
}

type sortableSliceOfObjectReference []apiv1.ObjectReference

func (s sortableSliceOfObjectReference) Len() int           { return len(s) }
func (s sortableSliceOfObjectReference) Less(i, j int) bool { return lessObjectReference(s[i], s[j]) }
func (s sortableSliceOfObjectReference) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// List returns the contents as a sorted apiv1.ObjectReference slice.
func (s ObjectReference) List() []apiv1.ObjectReference {
	res := make(sortableSliceOfObjectReference, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	sort.Sort(res)
	return []apiv1.ObjectReference(res)
}

// UnsortedList returns the slice with contents in random order.
func (s ObjectReference) UnsortedList() []apiv1.ObjectReference {
	res := make([]apiv1.ObjectReference, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

// Returns a single element from the set.
func (s ObjectReference) PopAny() (apiv1.ObjectReference, bool) {
	for key := range s {
		s.Delete(key)
		return key, true
	}
	var zeroValue apiv1.ObjectReference
	return zeroValue, false
}

// Len returns the size of the set.
func (s ObjectReference) Len() int {
	return len(s)
}

func lessObjectReference(lhs, rhs apiv1.ObjectReference) bool {
	if lhs.Namespace != rhs.Namespace {
		return lhs.Namespace < rhs.Namespace
	}
	return lhs.Name < rhs.Name
}
