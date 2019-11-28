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

package apiserver

import (
	"strconv"

	"kmodules.xyz/resource-metadata/apis/meta"
	"kmodules.xyz/resource-metadata/apis/meta/install"
	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/pkg/registry/meta/graphfinder"
	"kmodules.xyz/resource-metadata/pkg/registry/meta/pathfinder"
	"kmodules.xyz/resource-metadata/pkg/registry/meta/resourcedescriptor"

	v "github.com/appscode/go/version"
	semver "gomodules.xyz/version"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	install.Install(Scheme)

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	// Place you custom config here.
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig       *genericapiserver.RecommendedConfig
	InsecureServingInfo *genericapiserver.DeprecatedInsecureServingInfo
}

// ResourceMetadataServer contains state for a Kubernetes cluster master/api server.
type ResourceMetadataServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	GenericConfig       genericapiserver.CompletedConfig
	insecureServingInfo *genericapiserver.DeprecatedInsecureServingInfo
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		cfg.InsecureServingInfo,
	}

	c.GenericConfig.Version = &version.Info{
		GitVersion: v.Version.Version,
		GitCommit:  v.Version.CommitHash,
		BuildDate:  v.Version.BuildTimestamp,
		GoVersion:  v.Version.GoVersion,
		Compiler:   v.Version.Compiler,
		Platform:   v.Version.Platform,
	}
	if ver, err := semver.NewVersion(v.Version.Version); err == nil {
		c.GenericConfig.Version.Major = strconv.FormatInt(ver.Major(), 10)
		c.GenericConfig.Version.Minor = strconv.FormatInt(ver.Minor(), 10)
	}

	return CompletedConfig{&c}
}

// New returns a new instance of ResourceMetadataServer from the given config.
func (c completedConfig) New() (*ResourceMetadataServer, error) {
	genericServer, err := c.GenericConfig.New("resource-metadata-server", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	s := &ResourceMetadataServer{
		GenericAPIServer: genericServer,
	}

	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(meta.GroupName, Scheme, metav1.ParameterCodec, Codecs)

	v1alpha1storage := map[string]rest.Storage{}
	v1alpha1storage[v1alpha1.ResourceResourceDescriptors] = resourcedescriptor.NewStorage()
	v1alpha1storage[v1alpha1.ResourcePathFinders] = pathfinder.NewStorage()
	v1alpha1storage[v1alpha1.ResourceGraphFinders] = graphfinder.NewStorage()
	apiGroupInfo.VersionedResourcesStorageMap["v1alpha1"] = v1alpha1storage

	if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
		return nil, err
	}

	return s, nil
}
