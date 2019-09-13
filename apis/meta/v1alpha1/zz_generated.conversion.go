// +build !ignore_autogenerated

/*
Copyright 2019 The ResourceMetadata Project Authors.

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

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	meta "kmodules.xyz/resource-metadata/apis/meta"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Edge)(nil), (*meta.Edge)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Edge_To_meta_Edge(a.(*Edge), b.(*meta.Edge), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.Edge)(nil), (*Edge)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_Edge_To_v1alpha1_Edge(a.(*meta.Edge), b.(*Edge), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*GraphFinder)(nil), (*meta.GraphFinder)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_GraphFinder_To_meta_GraphFinder(a.(*GraphFinder), b.(*meta.GraphFinder), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.GraphFinder)(nil), (*GraphFinder)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_GraphFinder_To_v1alpha1_GraphFinder(a.(*meta.GraphFinder), b.(*GraphFinder), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*GraphRequest)(nil), (*meta.GraphRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_GraphRequest_To_meta_GraphRequest(a.(*GraphRequest), b.(*meta.GraphRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.GraphRequest)(nil), (*GraphRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_GraphRequest_To_v1alpha1_GraphRequest(a.(*meta.GraphRequest), b.(*GraphRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*GraphResponse)(nil), (*meta.GraphResponse)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_GraphResponse_To_meta_GraphResponse(a.(*GraphResponse), b.(*meta.GraphResponse), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.GraphResponse)(nil), (*GraphResponse)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_GraphResponse_To_v1alpha1_GraphResponse(a.(*meta.GraphResponse), b.(*GraphResponse), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Path)(nil), (*meta.Path)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Path_To_meta_Path(a.(*Path), b.(*meta.Path), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.Path)(nil), (*Path)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_Path_To_v1alpha1_Path(a.(*meta.Path), b.(*Path), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*PathFinder)(nil), (*meta.PathFinder)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PathFinder_To_meta_PathFinder(a.(*PathFinder), b.(*meta.PathFinder), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.PathFinder)(nil), (*PathFinder)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_PathFinder_To_v1alpha1_PathFinder(a.(*meta.PathFinder), b.(*PathFinder), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*PathRequest)(nil), (*meta.PathRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PathRequest_To_meta_PathRequest(a.(*PathRequest), b.(*meta.PathRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.PathRequest)(nil), (*PathRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_PathRequest_To_v1alpha1_PathRequest(a.(*meta.PathRequest), b.(*PathRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*PathResponse)(nil), (*meta.PathResponse)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PathResponse_To_meta_PathResponse(a.(*PathResponse), b.(*meta.PathResponse), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.PathResponse)(nil), (*PathResponse)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_PathResponse_To_v1alpha1_PathResponse(a.(*meta.PathResponse), b.(*PathResponse), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ResourceColumnDefinition)(nil), (*meta.ResourceColumnDefinition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ResourceColumnDefinition_To_meta_ResourceColumnDefinition(a.(*ResourceColumnDefinition), b.(*meta.ResourceColumnDefinition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.ResourceColumnDefinition)(nil), (*ResourceColumnDefinition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_ResourceColumnDefinition_To_v1alpha1_ResourceColumnDefinition(a.(*meta.ResourceColumnDefinition), b.(*ResourceColumnDefinition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ResourceConnection)(nil), (*meta.ResourceConnection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ResourceConnection_To_meta_ResourceConnection(a.(*ResourceConnection), b.(*meta.ResourceConnection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.ResourceConnection)(nil), (*ResourceConnection)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_ResourceConnection_To_v1alpha1_ResourceConnection(a.(*meta.ResourceConnection), b.(*ResourceConnection), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ResourceConnectionSpec)(nil), (*meta.ResourceConnectionSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ResourceConnectionSpec_To_meta_ResourceConnectionSpec(a.(*ResourceConnectionSpec), b.(*meta.ResourceConnectionSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.ResourceConnectionSpec)(nil), (*ResourceConnectionSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_ResourceConnectionSpec_To_v1alpha1_ResourceConnectionSpec(a.(*meta.ResourceConnectionSpec), b.(*ResourceConnectionSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ResourceDescriptor)(nil), (*meta.ResourceDescriptor)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ResourceDescriptor_To_meta_ResourceDescriptor(a.(*ResourceDescriptor), b.(*meta.ResourceDescriptor), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.ResourceDescriptor)(nil), (*ResourceDescriptor)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_ResourceDescriptor_To_v1alpha1_ResourceDescriptor(a.(*meta.ResourceDescriptor), b.(*ResourceDescriptor), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ResourceDescriptorList)(nil), (*meta.ResourceDescriptorList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ResourceDescriptorList_To_meta_ResourceDescriptorList(a.(*ResourceDescriptorList), b.(*meta.ResourceDescriptorList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.ResourceDescriptorList)(nil), (*ResourceDescriptorList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_ResourceDescriptorList_To_v1alpha1_ResourceDescriptorList(a.(*meta.ResourceDescriptorList), b.(*ResourceDescriptorList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ResourceDescriptorSpec)(nil), (*meta.ResourceDescriptorSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ResourceDescriptorSpec_To_meta_ResourceDescriptorSpec(a.(*ResourceDescriptorSpec), b.(*meta.ResourceDescriptorSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.ResourceDescriptorSpec)(nil), (*ResourceDescriptorSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_ResourceDescriptorSpec_To_v1alpha1_ResourceDescriptorSpec(a.(*meta.ResourceDescriptorSpec), b.(*ResourceDescriptorSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ResourceID)(nil), (*meta.ResourceID)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ResourceID_To_meta_ResourceID(a.(*ResourceID), b.(*meta.ResourceID), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*meta.ResourceID)(nil), (*ResourceID)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_meta_ResourceID_To_v1alpha1_ResourceID(a.(*meta.ResourceID), b.(*ResourceID), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_Edge_To_meta_Edge(in *Edge, out *meta.Edge, s conversion.Scope) error {
	out.W = in.W
	if err := Convert_v1alpha1_ResourceConnectionSpec_To_meta_ResourceConnectionSpec(&in.Connection, &out.Connection, s); err != nil {
		return err
	}
	out.Forward = in.Forward
	return nil
}

// Convert_v1alpha1_Edge_To_meta_Edge is an autogenerated conversion function.
func Convert_v1alpha1_Edge_To_meta_Edge(in *Edge, out *meta.Edge, s conversion.Scope) error {
	return autoConvert_v1alpha1_Edge_To_meta_Edge(in, out, s)
}

func autoConvert_meta_Edge_To_v1alpha1_Edge(in *meta.Edge, out *Edge, s conversion.Scope) error {
	out.W = in.W
	if err := Convert_meta_ResourceConnectionSpec_To_v1alpha1_ResourceConnectionSpec(&in.Connection, &out.Connection, s); err != nil {
		return err
	}
	out.Forward = in.Forward
	return nil
}

// Convert_meta_Edge_To_v1alpha1_Edge is an autogenerated conversion function.
func Convert_meta_Edge_To_v1alpha1_Edge(in *meta.Edge, out *Edge, s conversion.Scope) error {
	return autoConvert_meta_Edge_To_v1alpha1_Edge(in, out, s)
}

func autoConvert_v1alpha1_GraphFinder_To_meta_GraphFinder(in *GraphFinder, out *meta.GraphFinder, s conversion.Scope) error {
	out.Request = (*meta.GraphRequest)(unsafe.Pointer(in.Request))
	out.Response = (*meta.GraphResponse)(unsafe.Pointer(in.Response))
	return nil
}

// Convert_v1alpha1_GraphFinder_To_meta_GraphFinder is an autogenerated conversion function.
func Convert_v1alpha1_GraphFinder_To_meta_GraphFinder(in *GraphFinder, out *meta.GraphFinder, s conversion.Scope) error {
	return autoConvert_v1alpha1_GraphFinder_To_meta_GraphFinder(in, out, s)
}

func autoConvert_meta_GraphFinder_To_v1alpha1_GraphFinder(in *meta.GraphFinder, out *GraphFinder, s conversion.Scope) error {
	out.Request = (*GraphRequest)(unsafe.Pointer(in.Request))
	out.Response = (*GraphResponse)(unsafe.Pointer(in.Response))
	return nil
}

// Convert_meta_GraphFinder_To_v1alpha1_GraphFinder is an autogenerated conversion function.
func Convert_meta_GraphFinder_To_v1alpha1_GraphFinder(in *meta.GraphFinder, out *GraphFinder, s conversion.Scope) error {
	return autoConvert_meta_GraphFinder_To_v1alpha1_GraphFinder(in, out, s)
}

func autoConvert_v1alpha1_GraphRequest_To_meta_GraphRequest(in *GraphRequest, out *meta.GraphRequest, s conversion.Scope) error {
	return nil
}

// Convert_v1alpha1_GraphRequest_To_meta_GraphRequest is an autogenerated conversion function.
func Convert_v1alpha1_GraphRequest_To_meta_GraphRequest(in *GraphRequest, out *meta.GraphRequest, s conversion.Scope) error {
	return autoConvert_v1alpha1_GraphRequest_To_meta_GraphRequest(in, out, s)
}

func autoConvert_meta_GraphRequest_To_v1alpha1_GraphRequest(in *meta.GraphRequest, out *GraphRequest, s conversion.Scope) error {
	return nil
}

// Convert_meta_GraphRequest_To_v1alpha1_GraphRequest is an autogenerated conversion function.
func Convert_meta_GraphRequest_To_v1alpha1_GraphRequest(in *meta.GraphRequest, out *GraphRequest, s conversion.Scope) error {
	return autoConvert_meta_GraphRequest_To_v1alpha1_GraphRequest(in, out, s)
}

func autoConvert_v1alpha1_GraphResponse_To_meta_GraphResponse(in *GraphResponse, out *meta.GraphResponse, s conversion.Scope) error {
	out.Connections = *(*[]meta.Edge)(unsafe.Pointer(&in.Connections))
	return nil
}

// Convert_v1alpha1_GraphResponse_To_meta_GraphResponse is an autogenerated conversion function.
func Convert_v1alpha1_GraphResponse_To_meta_GraphResponse(in *GraphResponse, out *meta.GraphResponse, s conversion.Scope) error {
	return autoConvert_v1alpha1_GraphResponse_To_meta_GraphResponse(in, out, s)
}

func autoConvert_meta_GraphResponse_To_v1alpha1_GraphResponse(in *meta.GraphResponse, out *GraphResponse, s conversion.Scope) error {
	out.Connections = *(*[]Edge)(unsafe.Pointer(&in.Connections))
	return nil
}

// Convert_meta_GraphResponse_To_v1alpha1_GraphResponse is an autogenerated conversion function.
func Convert_meta_GraphResponse_To_v1alpha1_GraphResponse(in *meta.GraphResponse, out *GraphResponse, s conversion.Scope) error {
	return autoConvert_meta_GraphResponse_To_v1alpha1_GraphResponse(in, out, s)
}

func autoConvert_v1alpha1_Path_To_meta_Path(in *Path, out *meta.Path, s conversion.Scope) error {
	out.Distance = in.Distance
	out.Edges = *(*[]meta.Edge)(unsafe.Pointer(&in.Edges))
	return nil
}

// Convert_v1alpha1_Path_To_meta_Path is an autogenerated conversion function.
func Convert_v1alpha1_Path_To_meta_Path(in *Path, out *meta.Path, s conversion.Scope) error {
	return autoConvert_v1alpha1_Path_To_meta_Path(in, out, s)
}

func autoConvert_meta_Path_To_v1alpha1_Path(in *meta.Path, out *Path, s conversion.Scope) error {
	out.Distance = in.Distance
	out.Edges = *(*[]Edge)(unsafe.Pointer(&in.Edges))
	return nil
}

// Convert_meta_Path_To_v1alpha1_Path is an autogenerated conversion function.
func Convert_meta_Path_To_v1alpha1_Path(in *meta.Path, out *Path, s conversion.Scope) error {
	return autoConvert_meta_Path_To_v1alpha1_Path(in, out, s)
}

func autoConvert_v1alpha1_PathFinder_To_meta_PathFinder(in *PathFinder, out *meta.PathFinder, s conversion.Scope) error {
	out.Request = (*meta.PathRequest)(unsafe.Pointer(in.Request))
	out.Response = (*meta.PathResponse)(unsafe.Pointer(in.Response))
	return nil
}

// Convert_v1alpha1_PathFinder_To_meta_PathFinder is an autogenerated conversion function.
func Convert_v1alpha1_PathFinder_To_meta_PathFinder(in *PathFinder, out *meta.PathFinder, s conversion.Scope) error {
	return autoConvert_v1alpha1_PathFinder_To_meta_PathFinder(in, out, s)
}

func autoConvert_meta_PathFinder_To_v1alpha1_PathFinder(in *meta.PathFinder, out *PathFinder, s conversion.Scope) error {
	out.Request = (*PathRequest)(unsafe.Pointer(in.Request))
	out.Response = (*PathResponse)(unsafe.Pointer(in.Response))
	return nil
}

// Convert_meta_PathFinder_To_v1alpha1_PathFinder is an autogenerated conversion function.
func Convert_meta_PathFinder_To_v1alpha1_PathFinder(in *meta.PathFinder, out *PathFinder, s conversion.Scope) error {
	return autoConvert_meta_PathFinder_To_v1alpha1_PathFinder(in, out, s)
}

func autoConvert_v1alpha1_PathRequest_To_meta_PathRequest(in *PathRequest, out *meta.PathRequest, s conversion.Scope) error {
	out.Target = (*v1.TypeMeta)(unsafe.Pointer(in.Target))
	return nil
}

// Convert_v1alpha1_PathRequest_To_meta_PathRequest is an autogenerated conversion function.
func Convert_v1alpha1_PathRequest_To_meta_PathRequest(in *PathRequest, out *meta.PathRequest, s conversion.Scope) error {
	return autoConvert_v1alpha1_PathRequest_To_meta_PathRequest(in, out, s)
}

func autoConvert_meta_PathRequest_To_v1alpha1_PathRequest(in *meta.PathRequest, out *PathRequest, s conversion.Scope) error {
	out.Target = (*v1.TypeMeta)(unsafe.Pointer(in.Target))
	return nil
}

// Convert_meta_PathRequest_To_v1alpha1_PathRequest is an autogenerated conversion function.
func Convert_meta_PathRequest_To_v1alpha1_PathRequest(in *meta.PathRequest, out *PathRequest, s conversion.Scope) error {
	return autoConvert_meta_PathRequest_To_v1alpha1_PathRequest(in, out, s)
}

func autoConvert_v1alpha1_PathResponse_To_meta_PathResponse(in *PathResponse, out *meta.PathResponse, s conversion.Scope) error {
	out.Paths = *(*[]meta.Path)(unsafe.Pointer(&in.Paths))
	return nil
}

// Convert_v1alpha1_PathResponse_To_meta_PathResponse is an autogenerated conversion function.
func Convert_v1alpha1_PathResponse_To_meta_PathResponse(in *PathResponse, out *meta.PathResponse, s conversion.Scope) error {
	return autoConvert_v1alpha1_PathResponse_To_meta_PathResponse(in, out, s)
}

func autoConvert_meta_PathResponse_To_v1alpha1_PathResponse(in *meta.PathResponse, out *PathResponse, s conversion.Scope) error {
	out.Paths = *(*[]Path)(unsafe.Pointer(&in.Paths))
	return nil
}

// Convert_meta_PathResponse_To_v1alpha1_PathResponse is an autogenerated conversion function.
func Convert_meta_PathResponse_To_v1alpha1_PathResponse(in *meta.PathResponse, out *PathResponse, s conversion.Scope) error {
	return autoConvert_meta_PathResponse_To_v1alpha1_PathResponse(in, out, s)
}

func autoConvert_v1alpha1_ResourceColumnDefinition_To_meta_ResourceColumnDefinition(in *ResourceColumnDefinition, out *meta.ResourceColumnDefinition, s conversion.Scope) error {
	out.Name = in.Name
	out.Type = in.Type
	out.Format = in.Format
	out.Description = in.Description
	out.Priority = in.Priority
	out.JSONPath = in.JSONPath
	return nil
}

// Convert_v1alpha1_ResourceColumnDefinition_To_meta_ResourceColumnDefinition is an autogenerated conversion function.
func Convert_v1alpha1_ResourceColumnDefinition_To_meta_ResourceColumnDefinition(in *ResourceColumnDefinition, out *meta.ResourceColumnDefinition, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResourceColumnDefinition_To_meta_ResourceColumnDefinition(in, out, s)
}

func autoConvert_meta_ResourceColumnDefinition_To_v1alpha1_ResourceColumnDefinition(in *meta.ResourceColumnDefinition, out *ResourceColumnDefinition, s conversion.Scope) error {
	out.Name = in.Name
	out.Type = in.Type
	out.Format = in.Format
	out.Description = in.Description
	out.Priority = in.Priority
	out.JSONPath = in.JSONPath
	return nil
}

// Convert_meta_ResourceColumnDefinition_To_v1alpha1_ResourceColumnDefinition is an autogenerated conversion function.
func Convert_meta_ResourceColumnDefinition_To_v1alpha1_ResourceColumnDefinition(in *meta.ResourceColumnDefinition, out *ResourceColumnDefinition, s conversion.Scope) error {
	return autoConvert_meta_ResourceColumnDefinition_To_v1alpha1_ResourceColumnDefinition(in, out, s)
}

func autoConvert_v1alpha1_ResourceConnection_To_meta_ResourceConnection(in *ResourceConnection, out *meta.ResourceConnection, s conversion.Scope) error {
	if err := Convert_v1alpha1_ResourceConnectionSpec_To_meta_ResourceConnectionSpec(&in.ResourceConnectionSpec, &out.ResourceConnectionSpec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_ResourceConnection_To_meta_ResourceConnection is an autogenerated conversion function.
func Convert_v1alpha1_ResourceConnection_To_meta_ResourceConnection(in *ResourceConnection, out *meta.ResourceConnection, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResourceConnection_To_meta_ResourceConnection(in, out, s)
}

func autoConvert_meta_ResourceConnection_To_v1alpha1_ResourceConnection(in *meta.ResourceConnection, out *ResourceConnection, s conversion.Scope) error {
	if err := Convert_meta_ResourceConnectionSpec_To_v1alpha1_ResourceConnectionSpec(&in.ResourceConnectionSpec, &out.ResourceConnectionSpec, s); err != nil {
		return err
	}
	return nil
}

// Convert_meta_ResourceConnection_To_v1alpha1_ResourceConnection is an autogenerated conversion function.
func Convert_meta_ResourceConnection_To_v1alpha1_ResourceConnection(in *meta.ResourceConnection, out *ResourceConnection, s conversion.Scope) error {
	return autoConvert_meta_ResourceConnection_To_v1alpha1_ResourceConnection(in, out, s)
}

func autoConvert_v1alpha1_ResourceConnectionSpec_To_meta_ResourceConnectionSpec(in *ResourceConnectionSpec, out *meta.ResourceConnectionSpec, s conversion.Scope) error {
	out.Type = meta.ConnectionType(in.Type)
	out.NamespacePath = in.NamespacePath
	out.TargetLabelPath = in.TargetLabelPath
	out.SelectorPath = in.SelectorPath
	out.Selector = (*v1.LabelSelector)(unsafe.Pointer(in.Selector))
	out.NameTemplate = in.NameTemplate
	out.References = *(*[]string)(unsafe.Pointer(&in.References))
	out.Level = meta.OwnershipLevel(in.Level)
	return nil
}

// Convert_v1alpha1_ResourceConnectionSpec_To_meta_ResourceConnectionSpec is an autogenerated conversion function.
func Convert_v1alpha1_ResourceConnectionSpec_To_meta_ResourceConnectionSpec(in *ResourceConnectionSpec, out *meta.ResourceConnectionSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResourceConnectionSpec_To_meta_ResourceConnectionSpec(in, out, s)
}

func autoConvert_meta_ResourceConnectionSpec_To_v1alpha1_ResourceConnectionSpec(in *meta.ResourceConnectionSpec, out *ResourceConnectionSpec, s conversion.Scope) error {
	out.Type = ConnectionType(in.Type)
	out.NamespacePath = in.NamespacePath
	out.TargetLabelPath = in.TargetLabelPath
	out.SelectorPath = in.SelectorPath
	out.Selector = (*v1.LabelSelector)(unsafe.Pointer(in.Selector))
	out.NameTemplate = in.NameTemplate
	out.References = *(*[]string)(unsafe.Pointer(&in.References))
	out.Level = OwnershipLevel(in.Level)
	return nil
}

// Convert_meta_ResourceConnectionSpec_To_v1alpha1_ResourceConnectionSpec is an autogenerated conversion function.
func Convert_meta_ResourceConnectionSpec_To_v1alpha1_ResourceConnectionSpec(in *meta.ResourceConnectionSpec, out *ResourceConnectionSpec, s conversion.Scope) error {
	return autoConvert_meta_ResourceConnectionSpec_To_v1alpha1_ResourceConnectionSpec(in, out, s)
}

func autoConvert_v1alpha1_ResourceDescriptor_To_meta_ResourceDescriptor(in *ResourceDescriptor, out *meta.ResourceDescriptor, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_ResourceDescriptorSpec_To_meta_ResourceDescriptorSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_ResourceDescriptor_To_meta_ResourceDescriptor is an autogenerated conversion function.
func Convert_v1alpha1_ResourceDescriptor_To_meta_ResourceDescriptor(in *ResourceDescriptor, out *meta.ResourceDescriptor, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResourceDescriptor_To_meta_ResourceDescriptor(in, out, s)
}

func autoConvert_meta_ResourceDescriptor_To_v1alpha1_ResourceDescriptor(in *meta.ResourceDescriptor, out *ResourceDescriptor, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_meta_ResourceDescriptorSpec_To_v1alpha1_ResourceDescriptorSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_meta_ResourceDescriptor_To_v1alpha1_ResourceDescriptor is an autogenerated conversion function.
func Convert_meta_ResourceDescriptor_To_v1alpha1_ResourceDescriptor(in *meta.ResourceDescriptor, out *ResourceDescriptor, s conversion.Scope) error {
	return autoConvert_meta_ResourceDescriptor_To_v1alpha1_ResourceDescriptor(in, out, s)
}

func autoConvert_v1alpha1_ResourceDescriptorList_To_meta_ResourceDescriptorList(in *ResourceDescriptorList, out *meta.ResourceDescriptorList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]meta.ResourceDescriptor)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_ResourceDescriptorList_To_meta_ResourceDescriptorList is an autogenerated conversion function.
func Convert_v1alpha1_ResourceDescriptorList_To_meta_ResourceDescriptorList(in *ResourceDescriptorList, out *meta.ResourceDescriptorList, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResourceDescriptorList_To_meta_ResourceDescriptorList(in, out, s)
}

func autoConvert_meta_ResourceDescriptorList_To_v1alpha1_ResourceDescriptorList(in *meta.ResourceDescriptorList, out *ResourceDescriptorList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]ResourceDescriptor)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_meta_ResourceDescriptorList_To_v1alpha1_ResourceDescriptorList is an autogenerated conversion function.
func Convert_meta_ResourceDescriptorList_To_v1alpha1_ResourceDescriptorList(in *meta.ResourceDescriptorList, out *ResourceDescriptorList, s conversion.Scope) error {
	return autoConvert_meta_ResourceDescriptorList_To_v1alpha1_ResourceDescriptorList(in, out, s)
}

func autoConvert_v1alpha1_ResourceDescriptorSpec_To_meta_ResourceDescriptorSpec(in *ResourceDescriptorSpec, out *meta.ResourceDescriptorSpec, s conversion.Scope) error {
	if err := Convert_v1alpha1_ResourceID_To_meta_ResourceID(&in.Resource, &out.Resource, s); err != nil {
		return err
	}
	out.DisplayColumns = *(*[]meta.ResourceColumnDefinition)(unsafe.Pointer(&in.DisplayColumns))
	out.Connections = *(*[]meta.ResourceConnection)(unsafe.Pointer(&in.Connections))
	out.KeyTargets = *(*[]v1.TypeMeta)(unsafe.Pointer(&in.KeyTargets))
	return nil
}

// Convert_v1alpha1_ResourceDescriptorSpec_To_meta_ResourceDescriptorSpec is an autogenerated conversion function.
func Convert_v1alpha1_ResourceDescriptorSpec_To_meta_ResourceDescriptorSpec(in *ResourceDescriptorSpec, out *meta.ResourceDescriptorSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResourceDescriptorSpec_To_meta_ResourceDescriptorSpec(in, out, s)
}

func autoConvert_meta_ResourceDescriptorSpec_To_v1alpha1_ResourceDescriptorSpec(in *meta.ResourceDescriptorSpec, out *ResourceDescriptorSpec, s conversion.Scope) error {
	if err := Convert_meta_ResourceID_To_v1alpha1_ResourceID(&in.Resource, &out.Resource, s); err != nil {
		return err
	}
	out.DisplayColumns = *(*[]ResourceColumnDefinition)(unsafe.Pointer(&in.DisplayColumns))
	out.Connections = *(*[]ResourceConnection)(unsafe.Pointer(&in.Connections))
	out.KeyTargets = *(*[]v1.TypeMeta)(unsafe.Pointer(&in.KeyTargets))
	return nil
}

// Convert_meta_ResourceDescriptorSpec_To_v1alpha1_ResourceDescriptorSpec is an autogenerated conversion function.
func Convert_meta_ResourceDescriptorSpec_To_v1alpha1_ResourceDescriptorSpec(in *meta.ResourceDescriptorSpec, out *ResourceDescriptorSpec, s conversion.Scope) error {
	return autoConvert_meta_ResourceDescriptorSpec_To_v1alpha1_ResourceDescriptorSpec(in, out, s)
}

func autoConvert_v1alpha1_ResourceID_To_meta_ResourceID(in *ResourceID, out *meta.ResourceID, s conversion.Scope) error {
	out.Group = in.Group
	out.Version = in.Version
	out.Name = in.Name
	out.Kind = in.Kind
	out.Scope = meta.ResourceScope(in.Scope)
	return nil
}

// Convert_v1alpha1_ResourceID_To_meta_ResourceID is an autogenerated conversion function.
func Convert_v1alpha1_ResourceID_To_meta_ResourceID(in *ResourceID, out *meta.ResourceID, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResourceID_To_meta_ResourceID(in, out, s)
}

func autoConvert_meta_ResourceID_To_v1alpha1_ResourceID(in *meta.ResourceID, out *ResourceID, s conversion.Scope) error {
	out.Group = in.Group
	out.Version = in.Version
	out.Name = in.Name
	out.Kind = in.Kind
	out.Scope = ResourceScope(in.Scope)
	return nil
}

// Convert_meta_ResourceID_To_v1alpha1_ResourceID is an autogenerated conversion function.
func Convert_meta_ResourceID_To_v1alpha1_ResourceID(in *meta.ResourceID, out *ResourceID, s conversion.Scope) error {
	return autoConvert_meta_ResourceID_To_v1alpha1_ResourceID(in, out, s)
}
