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

package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"

	"kmodules.xyz/resource-metadata/apis/meta/v1alpha1"
	"kmodules.xyz/resource-metadata/pkg/apiserver"

	"github.com/go-openapi/spec"
	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/authentication/request/anonymous"
	"k8s.io/apiserver/pkg/authorization/authorizerfactory"
	genericapifilters "k8s.io/apiserver/pkg/endpoints/filters"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/server"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericfilters "k8s.io/apiserver/pkg/server/filters"
	genericoptions "k8s.io/apiserver/pkg/server/options"
)

// ResourceMetadataServerOptions contains state for master/api server
type ResourceMetadataServerOptions struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions
	InsecureServing         *genericoptions.DeprecatedInsecureServingOptionsWithLoopback
	Audit                   *genericoptions.AuditOptions
	Features                *genericoptions.FeatureOptions
	ProcessInfo             *genericoptions.ProcessInfo
	Webhook                 *genericoptions.WebhookOptions

	StdOut io.Writer
	StdErr io.Writer
}

// NewResourceMetadataServerOptions returns a new ResourceMetadataServerOptions
func NewResourceMetadataServerOptions(out, errOut io.Writer) *ResourceMetadataServerOptions {
	sso := genericoptions.NewSecureServingOptions()

	// We are composing recommended options for an aggregated api-server,
	// whose client is typically a proxy multiplexing many operations ---
	// notably including long-running ones --- into one HTTP/2 connection
	// into this server.  So allow many concurrent operations.
	sso.HTTP2MaxStreamsPerConnection = 1000

	o := &ResourceMetadataServerOptions{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		InsecureServing:         NewInsecureServingOptions(),
		Audit:                   genericoptions.NewAuditOptions(),
		Features:                genericoptions.NewFeatureOptions(),
		ProcessInfo:             genericoptions.NewProcessInfo("resource-metadata-server", "resource-metadata"),
		Webhook:                 genericoptions.NewWebhookOptions(),

		StdOut: out,
		StdErr: errOut,
	}
	o.GenericServerRunOptions.CorsAllowedOriginList = []string{".*"}
	return o
}

// NewInsecureServingOptions gives default values for the kube-apiserver.
// TODO: switch insecure serving off by default
func NewInsecureServingOptions() *genericoptions.DeprecatedInsecureServingOptionsWithLoopback {
	o := genericoptions.DeprecatedInsecureServingOptions{
		BindAddress: net.ParseIP("127.0.0.1"),
		BindPort:    8080,
	}
	if port, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		o.BindPort = port
	}
	return o.WithLoopback()
}

// NewCommandStartResourceMetadataServer provides a CLI handler for 'start master' command
// with a default ResourceMetadataServerOptions.
func NewCommandStartResourceMetadataServer(defaults *ResourceMetadataServerOptions, stopCh <-chan struct{}) *cobra.Command {
	o := *defaults
	cmd := &cobra.Command{
		Short: "Launch a resource metadata API server",
		Long:  "Launch a resource metadata API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunResourceMetadataServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.GenericServerRunOptions.AddUniversalFlags(flags)
	o.InsecureServing.AddFlags(flags)
	o.Audit.AddFlags(flags)
	o.Features.AddFlags(flags)

	return cmd
}

// Validate validates ResourceMetadataServerOptions
func (o ResourceMetadataServerOptions) Validate(args []string) error {
	var errors []error
	errors = append(errors, o.GenericServerRunOptions.Validate()...)
	errors = append(errors, o.InsecureServing.Validate()...)
	errors = append(errors, o.Audit.Validate()...)
	errors = append(errors, o.Features.Validate()...)
	return utilerrors.NewAggregate(errors)
}

// Complete fills in fields required to have valid data
func (o *ResourceMetadataServerOptions) Complete() error {
	return nil
}

// Config returns config for the api server given ResourceMetadataServerOptions
func (o *ResourceMetadataServerOptions) Config() (*apiserver.Config, error) {
	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)

	if err := o.GenericServerRunOptions.ApplyTo(&serverConfig.Config); err != nil {
		return nil, err
	}
	serverConfig.ExternalAddress = fmt.Sprintf("0.0.0.0:%d", o.InsecureServing.BindPort)

	var insecureServingInfo *genericapiserver.DeprecatedInsecureServingInfo
	if err := o.InsecureServing.ApplyTo(&insecureServingInfo, &serverConfig.Config.LoopbackClientConfig); err != nil {
		return nil, err
	}
	serverConfig.Config.Authentication = genericapiserver.AuthenticationInfo{
		APIAudiences:      nil,
		Authenticator:     anonymous.NewAuthenticator(),
		SupportsBasicAuth: false,
	}
	serverConfig.Config.Authorization = genericapiserver.AuthorizationInfo{
		Authorizer: authorizerfactory.NewAlwaysAllowAuthorizer(),
	}
	if err := o.Audit.ApplyTo(&serverConfig.Config, serverConfig.ClientConfig, serverConfig.SharedInformerFactory, o.ProcessInfo, o.Webhook); err != nil {
		return nil, err
	}
	if err := o.Features.ApplyTo(&serverConfig.Config); err != nil {
		return nil, err
	}

	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(v1alpha1.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(apiserver.Scheme))
	serverConfig.OpenAPIConfig.Info.InfoProps = spec.InfoProps{
		Title:   "resource-metadata",
		Version: "v0.1.0",
		Contact: &spec.ContactInfo{
			Name:  "AppsCode Inc.",
			URL:   "https://appscode.com",
			Email: "kmodules@appscode.com",
		},
		License: &spec.License{
			Name: "Apache 2.0",
			URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
		},
	}

	config := &apiserver.Config{
		GenericConfig:       serverConfig,
		InsecureServingInfo: insecureServingInfo,
	}
	return config, nil
}

// RunResourceMetadataServer starts a new ResourceMetadataServer given ResourceMetadataServerOptions
func (o ResourceMetadataServerOptions) RunResourceMetadataServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	s2 := server.GenericAPIServer.PrepareRun()

	insecureHandlerChain := BuildInsecureHandlerChain(server.GenericAPIServer.UnprotectedHandler(), &config.GenericConfig.Config)
	if err := config.InsecureServingInfo.Serve(insecureHandlerChain, config.GenericConfig.RequestTimeout, stopCh); err != nil {
		return err
	}

	return s2.Run(stopCh)
}

// BuildInsecureHandlerChain sets up the server to listen to http. Should be removed.
func BuildInsecureHandlerChain(apiHandler http.Handler, c *server.Config) http.Handler {
	handler := apiHandler
	handler = genericapifilters.WithAudit(handler, c.AuditBackend, c.AuditPolicyChecker, c.LongRunningFunc)
	handler = genericapifilters.WithAuthentication(handler, server.InsecureSuperuser{}, nil, nil)
	handler = genericfilters.WithCORS(handler, c.CorsAllowedOriginList, nil, nil, nil, "true")
	handler = genericfilters.WithTimeoutForNonLongRunningRequests(handler, c.LongRunningFunc, c.RequestTimeout)
	handler = genericfilters.WithMaxInFlightLimit(handler, c.MaxRequestsInFlight, c.MaxMutatingRequestsInFlight, c.LongRunningFunc)
	handler = genericfilters.WithWaitGroup(handler, c.LongRunningFunc, c.HandlerChainWaitGroup)
	handler = genericapifilters.WithRequestInfo(handler, server.NewRequestInfoResolver(c))
	handler = genericfilters.WithPanicRecovery(handler)

	return handler
}
