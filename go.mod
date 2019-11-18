module kmodules.xyz/resource-metadata

go 1.12

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/appscode/go v0.0.0-20191025021232-311ac347b3ef
	github.com/emicklei/dot v0.10.1
	github.com/go-openapi/jsonpointer v0.19.0 // indirect
	github.com/go-openapi/jsonreference v0.19.0 // indirect
	github.com/go-openapi/spec v0.19.0
	github.com/go-openapi/swag v0.17.2 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.6 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/json-iterator/go v1.1.6
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/sirupsen/logrus v1.4.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20190424203555-c05e17bb3b2d // indirect
	golang.org/x/net v0.0.0-20190424112056-4829fb13d2c6 // indirect
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a // indirect
	golang.org/x/sys v0.0.0-20190425145619-16072639606e // indirect
	golang.org/x/text v0.3.1 // indirect
	gomodules.xyz/version v0.1.0
	k8s.io/api v0.0.0-20190503110853-61630f889b3c
	k8s.io/apiextensions-apiserver v0.0.0-20190516231611-bf6753f2aa24
	k8s.io/apimachinery v0.0.0-20190508063446-a3da69d3723c
	k8s.io/apiserver v0.0.0-20190516230822-f89599b3f645
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/component-base v0.0.0-20190424053038-9fe063da3132 // indirect
	k8s.io/klog v0.3.1
	k8s.io/kube-openapi v0.0.0-20190502190224-411b2483e503
	k8s.io/utils v0.0.0-20190506122338-8fab8cb257d5 // indirect
	kmodules.xyz/client-go v0.0.0-20191107203555-78563b72c377
	sigs.k8s.io/yaml v1.1.0
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.34.0
	git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.4.2+incompatible
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.0.0-20190508045248-a52a97a7a2bf
	k8s.io/apiserver => github.com/kmodules/apiserver v0.0.0-20190811223248-5a95b2df4348
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190314002645-c892ea32361a
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190314000054-4a91899592f4
	k8s.io/klog => k8s.io/klog v0.3.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/utils => k8s.io/utils v0.0.0-20190514214443-0a167cbac756
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v0.0.0-20190302045857-e85c7b244fd2
)
