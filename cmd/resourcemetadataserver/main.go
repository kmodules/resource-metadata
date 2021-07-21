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

//go:generate make gen-bindata
package main

import (
	"os"

	"kmodules.xyz/resource-metadata/pkg/cmd/server"

	"gomodules.xyz/logs"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
)

func main() {
	stopCh := genericapiserver.SetupSignalHandler()
	options := server.NewResourceMetadataServerOptions(os.Stdout, os.Stderr)
	cmd := server.NewCommandStartResourceMetadataServer(options, stopCh)
	logs.Init(cmd, true)
	defer logs.FlushLogs()
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
