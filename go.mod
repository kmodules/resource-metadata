module kmodules.xyz/resource-metadata

go 1.12

require (
	github.com/appscode/go v0.0.0-20200323182826-54e98e09185a
	github.com/emicklei/dot v0.11.0
	github.com/go-openapi/spec v0.19.3
	github.com/gobuffalo/flect v0.2.1
	github.com/google/gofuzz v1.1.0
	github.com/hashicorp/golang-lru v0.5.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
	gomodules.xyz/version v0.1.0
	k8s.io/api v0.18.3
	k8s.io/apiextensions-apiserver v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/apiserver v0.18.3
	k8s.io/client-go v0.18.3
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6
	kmodules.xyz/apiversion v0.2.0
	kmodules.xyz/client-go v0.0.0-20200521065424-173e32c78a20
	kmodules.xyz/crd-schema-fuzz v0.0.0-20200521005638-2433a187de95
	sigs.k8s.io/yaml v1.2.0
)

replace (
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.19.0-alpha.0.0.20200520235721-10b58e57a423
	k8s.io/apiserver => github.com/kmodules/apiserver v0.18.4-0.20200521000930-14c5f6df9625
)
