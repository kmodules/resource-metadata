# Copyright AppsCode Inc. and Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL=/bin/bash -o pipefail

GO_PKG   := kmodules.xyz
REPO     := $(notdir $(shell pwd))
BIN      := resourcemetadataserver
COMPRESS ?= no

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
# CRD_OPTIONS          ?= "crd:trivialVersions=true"
CRD_OPTIONS          ?= "crd:allowDangerousTypes=true,crdVersions={v1}"
CODE_GENERATOR_IMAGE ?= ghcr.io/appscode/gengo:release-1.32
API_GROUPS           ?= core:v1alpha1 identity:v1alpha1 management:v1alpha1 meta:v1alpha1 node:v1alpha1 ui:v1alpha1

# Where to push the docker image.
REGISTRY ?= appscode

# This version-strategy uses git tags to set the version string
git_branch       := $(shell git rev-parse --abbrev-ref HEAD)
git_tag          := $(shell git describe --exact-match --abbrev=0 2>/dev/null || echo "")
commit_hash      := $(shell git rev-parse --verify HEAD)
commit_timestamp := $(shell date --date="@$$(git show -s --format=%ct)" --utc +%FT%T)

VERSION          := $(shell git describe --tags --always --dirty)
version_strategy := commit_hash
ifdef git_tag
	VERSION := $(git_tag)
	version_strategy := tag
else
	ifeq (,$(findstring $(git_branch),master HEAD))
		ifneq (,$(patsubst release-%,,$(git_branch)))
			VERSION := $(git_branch)
			version_strategy := branch
		endif
	endif
endif

###
### These variables should not need tweaking.
###

SRC_PKGS := apis cmd crds hub pkg
SRC_DIRS := $(SRC_PKGS) *.go # directories which hold app source (not vendored)

DOCKER_PLATFORMS := linux/amd64 linux/arm linux/arm64
BIN_PLATFORMS    := $(DOCKER_PLATFORMS)

# Used internally.  Users should pass GOOS and/or GOARCH.
OS   := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

BASEIMAGE_PROD   ?= gcr.io/distroless/static-debian12
BASEIMAGE_DBG    ?= debian:12

IMAGE            := $(REGISTRY)/$(BIN)
VERSION_PROD     := $(VERSION)
VERSION_DBG      := $(VERSION)-dbg
TAG              := $(VERSION)_$(OS)_$(ARCH)
TAG_PROD         := $(TAG)
TAG_DBG          := $(VERSION)-dbg_$(OS)_$(ARCH)

GO_VERSION       ?= 1.24
BUILD_IMAGE      ?= ghcr.io/appscode/golang-dev:$(GO_VERSION)

OUTBIN = bin/$(OS)_$(ARCH)/$(BIN)
ifeq ($(OS),windows)
  OUTBIN = bin/$(OS)_$(ARCH)/$(BIN).exe
endif

# Directories that we need created to build/test.
BUILD_DIRS  := bin/$(OS)_$(ARCH)     \
               .go/bin/$(OS)_$(ARCH) \
               .go/cache

DOCKERFILE_PROD  = Dockerfile.in
DOCKERFILE_DBG   = Dockerfile.dbg

DOCKER_REPO_ROOT := /go/src/$(GO_PKG)/$(REPO)

# If you want to build all binaries, see the 'all-build' rule.
# If you want to build all containers, see the 'all-container' rule.
# If you want to build AND push all containers, see the 'all-push' rule.
all: fmt build

# For the following OS/ARCH expansions, we transform OS/ARCH into OS_ARCH
# because make pattern rules don't match with embedded '/' characters.

build-%:
	@$(MAKE) build                        \
	    --no-print-directory              \
	    GOOS=$(firstword $(subst _, ,$*)) \
	    GOARCH=$(lastword $(subst _, ,$*))

container-%:
	@$(MAKE) container                    \
	    --no-print-directory              \
	    GOOS=$(firstword $(subst _, ,$*)) \
	    GOARCH=$(lastword $(subst _, ,$*))

push-%:
	@$(MAKE) push                         \
	    --no-print-directory              \
	    GOOS=$(firstword $(subst _, ,$*)) \
	    GOARCH=$(lastword $(subst _, ,$*))

all-build: $(addprefix build-, $(subst /,_, $(BIN_PLATFORMS)))

all-container: $(addprefix container-, $(subst /,_, $(DOCKER_PLATFORMS)))

all-push: $(addprefix push-, $(subst /,_, $(DOCKER_PLATFORMS)))

version:
	@echo version=$(VERSION)
	@echo version_strategy=$(version_strategy)
	@echo git_tag=$(git_tag)
	@echo git_branch=$(git_branch)
	@echo commit_hash=$(commit_hash)
	@echo commit_timestamp=$(commit_timestamp)

# Generate code for Kubernetes types
.PHONY: clientset
clientset:
	@docker run --rm	                                 \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
		--env HTTP_PROXY=$(HTTP_PROXY)                   \
		--env HTTPS_PROXY=$(HTTPS_PROXY)                 \
		$(CODE_GENERATOR_IMAGE)                          \
		deepcopy-gen                                     \
			--go-header-file "./hack/license/go.txt"       \
			--input-dirs "$(GO_PKG)/$(REPO)/apis/shared"   \
			--output-file-base zz_generated.deepcopy
	@docker run --rm                                   \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
		--env HTTP_PROXY=$(HTTP_PROXY)                   \
		--env HTTPS_PROXY=$(HTTPS_PROXY)                 \
		$(CODE_GENERATOR_IMAGE)                          \
		/go/src/k8s.io/code-generator/generate-groups.sh \
			deepcopy,client                                \
			$(GO_PKG)/$(REPO)/client                       \
			$(GO_PKG)/$(REPO)/apis                         \
			"$(API_GROUPS)"                                \
			--go-header-file "./hack/license/go.txt"
# 	@docker run --rm                                            \
# 		-u $$(id -u):$$(id -g)                                    \
# 		-v /tmp:/.cache                                           \
# 		-v $$(pwd):$(DOCKER_REPO_ROOT)                            \
# 		-w $(DOCKER_REPO_ROOT)                                    \
# 		--env HTTP_PROXY=$(HTTP_PROXY)                            \
# 		--env HTTPS_PROXY=$(HTTPS_PROXY)                          \
# 		$(CODE_GENERATOR_IMAGE)                                   \
# 		/go/src/k8s.io/code-generator/generate-internal-groups.sh \
# 			"conversion"                                            \
# 			$(GO_PKG)/$(REPO)/client                                \
# 			$(GO_PKG)/$(REPO)/apis                                  \
# 			$(GO_PKG)/$(REPO)/apis                                  \
# 			"$(API_GROUPS)"                                         \
# 			--go-header-file "./hack/license/go.txt" --extra-dirs k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1
# 	# for both CRD and EAS types
# 	@docker run --rm                                          \
# 		-u $$(id -u):$$(id -g)                                  \
# 		-v /tmp:/.cache                                         \
# 		-v $$(pwd):$(DOCKER_REPO_ROOT)                          \
# 		-w $(DOCKER_REPO_ROOT)                                  \
# 		--env HTTP_PROXY=$(HTTP_PROXY)                          \
# 		--env HTTPS_PROXY=$(HTTPS_PROXY)                        \
# 		$(CODE_GENERATOR_IMAGE)                                 \
# 		/go/src/k8s.io/code-generator/generate-groups.sh        \
# 			all                                                   \
# 			$(GO_PKG)/$(REPO)/client                              \
# 			$(GO_PKG)/$(REPO)/apis                                \
# 			"$(API_GROUPS)"                                       \
# 			--go-header-file "./hack/license/go.txt"

# Generate openapi schema
.PHONY: openapi
openapi: openapi-shared $(addprefix openapi-, $(subst :,_, $(API_GROUPS)))
openapi-shared:
	@echo "Generating openapi schema"
	@docker run --rm	                                    \
		-u $$(id -u):$$(id -g)                              \
		-v /tmp:/.cache                                     \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                      \
		-w $(DOCKER_REPO_ROOT)                              \
		--env HTTP_PROXY=$(HTTP_PROXY)                      \
		--env HTTPS_PROXY=$(HTTPS_PROXY)                    \
		$(CODE_GENERATOR_IMAGE)                             \
		openapi-gen                                         \
			--v 1 --logtostderr                               \
			--go-header-file "./hack/license/go.txt"          \
			--input-dirs "$(GO_PKG)/$(REPO)/apis/shared"      \
			--output-package "$(GO_PKG)/$(REPO)/apis/shared"  \
			--report-filename /tmp/violation_exceptions.list
openapi-%:
	@echo "Generating openapi schema for $(subst _,/,$*)"
	@mkdir -p .config/api-rules
	@docker run --rm                                   \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
		--env HTTP_PROXY=$(HTTP_PROXY)                   \
		--env HTTPS_PROXY=$(HTTPS_PROXY)                 \
		$(CODE_GENERATOR_IMAGE)                          \
		openapi-gen                                      \
			--v 1 --logtostderr                            \
			--go-header-file "./hack/license/go.txt"       \
			--input-dirs "$(GO_PKG)/$(REPO)/apis/$(subst _,/,$*),$(GO_PKG)/$(REPO)/apis/shared,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/api/resource,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/version,k8s.io/api/core/v1,k8s.io/apimachinery/pkg/util/intstr,kmodules.xyz/client-go/api/v1,go.bytebuilders.dev/catalog/api/v1alpha1,kmodules.xyz/offshoot-api/api/v1" \
			--output-package "$(GO_PKG)/$(REPO)/apis/$(subst _,/,$*)" \
			--report-filename .config/api-rules/violation_exceptions.list

# Generate CRD manifests
.PHONY: gen-crds
gen-crds:
	@echo "Generating CRD manifests"
	@docker run --rm 	                    \
		-u $$(id -u):$$(id -g)              \
		-v /tmp:/.cache                     \
		-v $$(pwd):$(DOCKER_REPO_ROOT)      \
		-w $(DOCKER_REPO_ROOT)              \
	    --env HTTP_PROXY=$(HTTP_PROXY)    \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)  \
		$(CODE_GENERATOR_IMAGE)             \
		controller-gen                      \
			$(CRD_OPTIONS)                    \
			paths="./apis/..."                \
			output:crd:artifacts:config=crds

.PHONY: manifests
manifests: gen-crds patch-schema

.PHONY: patch-schema
patch-schema: $(BUILD_DIRS)
	@echo "patching crds/core.k8s.appscode.com_podviews.yaml"
	@kubectl patch -f crds/core.k8s.appscode.com_podviews.yaml -p "$$(cat hack/podview-schema-patch.json)" --type=json --local=true -o yaml > bin/core.k8s.appscode.com_podviews.yaml
	@mv bin/core.k8s.appscode.com_podviews.yaml crds/core.k8s.appscode.com_podviews.yaml

.PHONY: gen
gen: clientset manifests openapi

fmt: $(BUILD_DIRS)
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    $(BUILD_IMAGE)                                          \
	    /bin/bash -c "                                          \
	        REPO_PKG=$(GO_PKG)                                  \
	        ./hack/fmt.sh $(SRC_DIRS)                           \
	    "
	go run ./cmd/check-edge-label/main.go
	go run ./cmd/resource-fmt/main.go

build: $(OUTBIN)

# The following structure defeats Go's (intentional) behavior to always touch
# result files, even if they have not changed.  This will still run `go` but
# will not trigger further work if nothing has actually changed.

$(OUTBIN): .go/$(OUTBIN).stamp
	@true

# This will build the binary under ./.go and update the real binary iff needed.
.PHONY: .go/$(OUTBIN).stamp
.go/$(OUTBIN).stamp: $(BUILD_DIRS)
	@echo "making $(OUTBIN)"
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    $(BUILD_IMAGE)                                          \
	    /bin/bash -c "                                          \
	        ARCH=$(ARCH)                                        \
	        OS=$(OS)                                            \
	        VERSION=$(VERSION)                                  \
	        version_strategy=$(version_strategy)                \
	        git_branch=$(git_branch)                            \
	        git_tag=$(git_tag)                                  \
	        commit_hash=$(commit_hash)                          \
	        commit_timestamp=$(commit_timestamp)                \
	        ./hack/build.sh                                     \
	    "
	@if [ $(COMPRESS) = yes ] && [ $(OS) != darwin ]; then        \
		echo "compressing $(OUTBIN)";                               \
		docker run                                                  \
		    -i                                                      \
		    --rm                                                    \
		    -u $$(id -u):$$(id -g)                                  \
		    -v $$(pwd):/src                                         \
		    -w /src                                                 \
		    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
		    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
		    -v $$(pwd)/.go/cache:/.cache                            \
		    --env HTTP_PROXY=$(HTTP_PROXY)                          \
		    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
		    $(BUILD_IMAGE)                                          \
		    upx --brute /go/$(OUTBIN);                              \
	fi
	@if ! cmp -s .go/$(OUTBIN) $(OUTBIN); then \
	    mv .go/$(OUTBIN) $(OUTBIN);            \
	    date >$@;                              \
	fi
	@echo

# Used to track state in hidden files.
DOTFILE_IMAGE    = $(subst /,_,$(IMAGE))-$(TAG)

container: bin/.container-$(DOTFILE_IMAGE)-PROD bin/.container-$(DOTFILE_IMAGE)-DBG
bin/.container-$(DOTFILE_IMAGE)-%: bin/$(OS)_$(ARCH)/$(BIN) $(DOCKERFILE_%)
	@echo "container: $(IMAGE):$(TAG_$*)"
	@sed                                  \
		-e 's|{ARG_BIN}|$(BIN)|g'           \
		-e 's|{ARG_ARCH}|$(ARCH)|g'         \
		-e 's|{ARG_OS}|$(OS)|g'             \
		-e 's|{ARG_FROM}|$(BASEIMAGE_$*)|g' \
		$(DOCKERFILE_$*) > bin/.dockerfile-$*-$(OS)_$(ARCH)
	@DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --platform $(OS)/$(ARCH) --load --pull -t $(IMAGE):$(TAG_$*) -f bin/.dockerfile-$*-$(OS)_$(ARCH) .
	@docker images -q $(IMAGE):$(TAG_$*) > $@
	@echo

push: bin/.push-$(DOTFILE_IMAGE)-PROD bin/.push-$(DOTFILE_IMAGE)-DBG
bin/.push-$(DOTFILE_IMAGE)-%: bin/.container-$(DOTFILE_IMAGE)-%
	@docker push $(IMAGE):$(TAG_$*)
	@echo "pushed: $(IMAGE):$(TAG_$*)"
	@echo

.PHONY: docker-manifest
docker-manifest: docker-manifest-PROD docker-manifest-DBG
docker-manifest-%:
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create -a $(IMAGE):$(VERSION_$*) $(foreach PLATFORM,$(DOCKER_PLATFORMS),$(IMAGE):$(VERSION_$*)_$(subst /,_,$(PLATFORM)))
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push $(IMAGE):$(VERSION_$*)

.PHONY: test
test: unit-tests

unit-tests: $(BUILD_DIRS)
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    --net=host                                              \
	    -v $(HOME)/.kube:/.kube                                 \
	    -v $(HOME)/.minikube:$(HOME)/.minikube                  \
	    -v $(HOME)/.credentials:$(HOME)/.credentials            \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    --env KUBECONFIG=$(KUBECONFIG)                          \
	    $(BUILD_IMAGE)                                          \
	    /bin/bash -c "                                          \
	        ARCH=$(ARCH)                                        \
	        OS=$(OS)                                            \
	        VERSION=$(VERSION)                                  \
	        KUBECONFIG=$${KUBECONFIG#$(HOME)}                   \
	        ./hack/test.sh $(SRC_PKGS)                          \
	    "

ADDTL_LINTERS   := gofmt,goimports,unparam

.PHONY: lint
lint: $(BUILD_DIRS)
	@echo "running linter"
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    --env GO111MODULE=on                                    \
	    --env GOFLAGS="-mod=vendor"                             \
	    $(BUILD_IMAGE)                                          \
	    golangci-lint run --enable $(ADDTL_LINTERS) --max-same-issues=100 --timeout=60m --exclude-files="generated.*\.go$\" --exclude-dirs-use-default --exclude-dirs=client,vendor

$(BUILD_DIRS):
	@mkdir -p $@

.PHONY: dev
dev: gen fmt push

.PHONY: verify
verify: verify-gen verify-modules

.PHONY: verify-modules
verify-modules:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor
	@if !(git diff --exit-code HEAD); then \
		echo "go module files are out of date"; exit 1; \
	fi

.PHONY: verify-gen
verify-gen: gen fmt
	@if !(git diff --exit-code HEAD); then \
		echo "files are out of date, run make gen fmt"; exit 1; \
	fi

.PHONY: add-license
add-license:
	@echo "Adding license header"
	@docker run --rm 	                                 \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
		--env HTTP_PROXY=$(HTTP_PROXY)                   \
		--env HTTPS_PROXY=$(HTTPS_PROXY)                 \
		$(BUILD_IMAGE)                                   \
		ltag -t "./hack/license" --excludes "vendor contrib" -v

.PHONY: check-license
check-license:
	@echo "Checking files for license header"
	@docker run --rm 	                                 \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
		--env HTTP_PROXY=$(HTTP_PROXY)                   \
		--env HTTPS_PROXY=$(HTTPS_PROXY)                 \
		$(BUILD_IMAGE)                                   \
		ltag -t "./hack/license" --excludes "vendor contrib" --check -v

.PHONY: ci
ci: verify check-license lint build unit-tests #cover

.PHONY: qa
qa:
	@if [ "$$APPSCODE_ENV" = "prod" ]; then                                              \
		echo "Nothing to do in prod env. Are you trying to 'release' binaries to prod?"; \
		exit 1;                                                                          \
	fi
	@if [ "$(version_strategy)" = "tag" ]; then               \
		echo "Are you trying to 'release' binaries to prod?"; \
		exit 1;                                               \
	fi
	@$(MAKE) clean all-push docker-manifest --no-print-directory

.PHONY: release
release:
	@if [ "$$APPSCODE_ENV" != "prod" ]; then      \
		echo "'release' only works in PROD env."; \
		exit 1;                                   \
	fi
	@if [ "$(version_strategy)" != "tag" ]; then                    \
		echo "apply tag to release binaries and/or docker images."; \
		exit 1;                                                     \
	fi
	@$(MAKE) clean all-push docker-manifest --no-print-directory

.PHONY: clean
clean:
	rm -rf .go bin

.PHONY: publish-icons
publish-icons:
	@echo "publishing icons"
	gsutil -m rsync -d -r $$(pwd)/icons gs://cdn.appscode.com/k8s/icons
	gsutil -m acl ch -u AllUsers:R -r gs://cdn.appscode.com/k8s/icons
