//go:generate make gen-bindata
package main

import (
	"flag"
	"os"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog"
	"kmodules.xyz/client-go/logs"
	"kmodules.xyz/resource-metadata/pkg/cmd/server"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	stopCh := genericapiserver.SetupSignalHandler()
	options := server.NewResourceMetadataServerOptions(os.Stdout, os.Stderr)
	cmd := server.NewCommandStartResourceMetadataServer(options, stopCh)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	logs.ParseFlags()
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
