// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package meta

import (
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContactData) DeepCopyInto(out *ContactData) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContactData.
func (in *ContactData) DeepCopy() *ContactData {
	if in == nil {
		return nil
	}
	out := new(ContactData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Edge) DeepCopyInto(out *Edge) {
	*out = *in
	out.Src = in.Src
	out.Dst = in.Dst
	in.Connection.DeepCopyInto(&out.Connection)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Edge.
func (in *Edge) DeepCopy() *Edge {
	if in == nil {
		return nil
	}
	out := new(Edge)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Entry) DeepCopyInto(out *Entry) {
	*out = *in
	out.Type = in.Type
	if in.Icons != nil {
		in, out := &in.Icons, &out.Icons
		*out = make([]ImageSpec, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Entry.
func (in *Entry) DeepCopy() *Entry {
	if in == nil {
		return nil
	}
	out := new(Entry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GraphFinder) DeepCopyInto(out *GraphFinder) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Request != nil {
		in, out := &in.Request, &out.Request
		*out = new(GraphRequest)
		**out = **in
	}
	if in.Response != nil {
		in, out := &in.Response, &out.Response
		*out = new(GraphResponse)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GraphFinder.
func (in *GraphFinder) DeepCopy() *GraphFinder {
	if in == nil {
		return nil
	}
	out := new(GraphFinder)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GraphFinder) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GraphRequest) DeepCopyInto(out *GraphRequest) {
	*out = *in
	out.Source = in.Source
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GraphRequest.
func (in *GraphRequest) DeepCopy() *GraphRequest {
	if in == nil {
		return nil
	}
	out := new(GraphRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GraphResponse) DeepCopyInto(out *GraphResponse) {
	*out = *in
	out.Source = in.Source
	if in.Connections != nil {
		in, out := &in.Connections, &out.Connections
		*out = make([]Edge, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GraphResponse.
func (in *GraphResponse) DeepCopy() *GraphResponse {
	if in == nil {
		return nil
	}
	out := new(GraphResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GroupVersionResource) DeepCopyInto(out *GroupVersionResource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GroupVersionResource.
func (in *GroupVersionResource) DeepCopy() *GroupVersionResource {
	if in == nil {
		return nil
	}
	out := new(GroupVersionResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImageSpec) DeepCopyInto(out *ImageSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImageSpec.
func (in *ImageSpec) DeepCopy() *ImageSpec {
	if in == nil {
		return nil
	}
	out := new(ImageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Link) DeepCopyInto(out *Link) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Link.
func (in *Link) DeepCopy() *Link {
	if in == nil {
		return nil
	}
	out := new(Link)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Path) DeepCopyInto(out *Path) {
	*out = *in
	out.Source = in.Source
	out.Target = in.Target
	if in.Edges != nil {
		in, out := &in.Edges, &out.Edges
		*out = make([]Edge, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Path.
func (in *Path) DeepCopy() *Path {
	if in == nil {
		return nil
	}
	out := new(Path)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PathFinder) DeepCopyInto(out *PathFinder) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Request != nil {
		in, out := &in.Request, &out.Request
		*out = new(PathRequest)
		(*in).DeepCopyInto(*out)
	}
	if in.Response != nil {
		in, out := &in.Response, &out.Response
		*out = new(PathResponse)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PathFinder.
func (in *PathFinder) DeepCopy() *PathFinder {
	if in == nil {
		return nil
	}
	out := new(PathFinder)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PathFinder) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PathRequest) DeepCopyInto(out *PathRequest) {
	*out = *in
	out.Source = in.Source
	if in.Target != nil {
		in, out := &in.Target, &out.Target
		*out = new(GroupVersionResource)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PathRequest.
func (in *PathRequest) DeepCopy() *PathRequest {
	if in == nil {
		return nil
	}
	out := new(PathRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PathResponse) DeepCopyInto(out *PathResponse) {
	*out = *in
	if in.Paths != nil {
		in, out := &in.Paths, &out.Paths
		*out = make([]Path, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PathResponse.
func (in *PathResponse) DeepCopy() *PathResponse {
	if in == nil {
		return nil
	}
	out := new(PathResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClass) DeepCopyInto(out *ResourceClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClass.
func (in *ResourceClass) DeepCopy() *ResourceClass {
	if in == nil {
		return nil
	}
	out := new(ResourceClass)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClassList) DeepCopyInto(out *ResourceClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceClass, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClassList.
func (in *ResourceClassList) DeepCopy() *ResourceClassList {
	if in == nil {
		return nil
	}
	out := new(ResourceClassList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceClassSpec) DeepCopyInto(out *ResourceClassSpec) {
	*out = *in
	if in.Entries != nil {
		in, out := &in.Entries, &out.Entries
		*out = make([]Entry, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Icons != nil {
		in, out := &in.Icons, &out.Icons
		*out = make([]ImageSpec, len(*in))
		copy(*out, *in)
	}
	if in.Maintainers != nil {
		in, out := &in.Maintainers, &out.Maintainers
		*out = make([]ContactData, len(*in))
		copy(*out, *in)
	}
	if in.Links != nil {
		in, out := &in.Links, &out.Links
		*out = make([]Link, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceClassSpec.
func (in *ResourceClassSpec) DeepCopy() *ResourceClassSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceClassSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceColumnDefinition) DeepCopyInto(out *ResourceColumnDefinition) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceColumnDefinition.
func (in *ResourceColumnDefinition) DeepCopy() *ResourceColumnDefinition {
	if in == nil {
		return nil
	}
	out := new(ResourceColumnDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceConnection) DeepCopyInto(out *ResourceConnection) {
	*out = *in
	out.Target = in.Target
	in.ResourceConnectionSpec.DeepCopyInto(&out.ResourceConnectionSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceConnection.
func (in *ResourceConnection) DeepCopy() *ResourceConnection {
	if in == nil {
		return nil
	}
	out := new(ResourceConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceConnectionSpec) DeepCopyInto(out *ResourceConnectionSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.References != nil {
		in, out := &in.References, &out.References
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceConnectionSpec.
func (in *ResourceConnectionSpec) DeepCopy() *ResourceConnectionSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceConnectionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceDescriptor) DeepCopyInto(out *ResourceDescriptor) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceDescriptor.
func (in *ResourceDescriptor) DeepCopy() *ResourceDescriptor {
	if in == nil {
		return nil
	}
	out := new(ResourceDescriptor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceDescriptor) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceDescriptorList) DeepCopyInto(out *ResourceDescriptorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceDescriptor, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceDescriptorList.
func (in *ResourceDescriptorList) DeepCopy() *ResourceDescriptorList {
	if in == nil {
		return nil
	}
	out := new(ResourceDescriptorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceDescriptorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceDescriptorSpec) DeepCopyInto(out *ResourceDescriptorSpec) {
	*out = *in
	out.Resource = in.Resource
	if in.Columns != nil {
		in, out := &in.Columns, &out.Columns
		*out = make([]ResourceColumnDefinition, len(*in))
		copy(*out, *in)
	}
	if in.SubTables != nil {
		in, out := &in.SubTables, &out.SubTables
		*out = make([]ResourceSubTableDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Connections != nil {
		in, out := &in.Connections, &out.Connections
		*out = make([]ResourceConnection, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.KeyTargets != nil {
		in, out := &in.KeyTargets, &out.KeyTargets
		*out = make([]v1.TypeMeta, len(*in))
		copy(*out, *in)
	}
	if in.Validation != nil {
		in, out := &in.Validation, &out.Validation
		*out = new(apiextensions.CustomResourceValidation)
		(*in).DeepCopyInto(*out)
	}
	if in.Icons != nil {
		in, out := &in.Icons, &out.Icons
		*out = make([]ImageSpec, len(*in))
		copy(*out, *in)
	}
	if in.Maintainers != nil {
		in, out := &in.Maintainers, &out.Maintainers
		*out = make([]ContactData, len(*in))
		copy(*out, *in)
	}
	if in.Links != nil {
		in, out := &in.Links, &out.Links
		*out = make([]Link, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceDescriptorSpec.
func (in *ResourceDescriptorSpec) DeepCopy() *ResourceDescriptorSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceDescriptorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceID) DeepCopyInto(out *ResourceID) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceID.
func (in *ResourceID) DeepCopy() *ResourceID {
	if in == nil {
		return nil
	}
	out := new(ResourceID)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSubTableDefinition) DeepCopyInto(out *ResourceSubTableDefinition) {
	*out = *in
	if in.Columns != nil {
		in, out := &in.Columns, &out.Columns
		*out = make([]ResourceColumnDefinition, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSubTableDefinition.
func (in *ResourceSubTableDefinition) DeepCopy() *ResourceSubTableDefinition {
	if in == nil {
		return nil
	}
	out := new(ResourceSubTableDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubTable) DeepCopyInto(out *SubTable) {
	*out = *in
	if in.ColumnDefinitions != nil {
		in, out := &in.ColumnDefinitions, &out.ColumnDefinitions
		*out = make([]ResourceColumnDefinition, len(*in))
		copy(*out, *in)
	}
	if in.Rows != nil {
		in, out := &in.Rows, &out.Rows
		*out = make([]TableRow, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubTable.
func (in *SubTable) DeepCopy() *SubTable {
	if in == nil {
		return nil
	}
	out := new(SubTable)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Table) DeepCopyInto(out *Table) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.ColumnDefinitions != nil {
		in, out := &in.ColumnDefinitions, &out.ColumnDefinitions
		*out = make([]ResourceColumnDefinition, len(*in))
		copy(*out, *in)
	}
	if in.Rows != nil {
		in, out := &in.Rows, &out.Rows
		*out = make([]TableRow, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SubTables != nil {
		in, out := &in.SubTables, &out.SubTables
		*out = make([]SubTable, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Table.
func (in *Table) DeepCopy() *Table {
	if in == nil {
		return nil
	}
	out := new(Table)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TableOptions) DeepCopyInto(out *TableOptions) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TableOptions.
func (in *TableOptions) DeepCopy() *TableOptions {
	if in == nil {
		return nil
	}
	out := new(TableOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TableRow) DeepCopyInto(out *TableRow) {
	clone := in.DeepCopy()
	*out = *clone
	return
}
