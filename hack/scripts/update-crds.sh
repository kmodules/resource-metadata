#!/bin/bash

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

update_crds() {
    if test -d "$1"; then
        echo "refreshing $1"
        pushd "$1"
        git fetch origin --prune --tags -f
        git reset --hard HEAD
        git checkout master
        git pull origin master
        popd
        go run cmd/import-crds/main.go --input="$1"
        echo "-------------------------------------------------"
    fi
}

update_crds $HOME/go/src/kubeops.dev/supervisor/crds

update_crds $HOME/go/src/kubeops.dev/supervisor/crds

update_crds $HOME/go/src/kubeops.dev/ui-server/crds
update_crds $HOME/go/src/stash.appscode.dev/ui-server/crds
update_crds $HOME/go/src/go.openviz.dev/grafana-tools/crds

update_crds $HOME/go/src/k8s.io/api/crds
update_crds $HOME/go/src/k8s.io/kube-aggregator/crds

update_crds $HOME/go/src/github.com/prometheus-operator/prometheus-operator/example/prometheus-operator-crd
update_crds $HOME/go/src/github.com/cert-manager/cert-manager/deploy/crds
update_crds $HOME/go/src/voyagermesh.dev/apimachinery/crds

update_crds $HOME/go/src/stash.appscode.dev/apimachinery/crds
update_crds $HOME/go/src/kmodules.xyz/custom-resources/crds
update_crds $HOME/go/src/kubedb.dev/apimachinery/crds

update_crds $HOME/go/src/kubevault.dev/apimachinery/crds
update_crds $HOME/go/src/sigs.k8s.io/secrets-store-csi-driver/charts/secrets-store-csi-driver/crds

update_crds $HOME/go/src/github.com/kubernetes-csi/external-snapshotter/client/config/crd

update_crds $HOME/go/src/k8s.io/autoscaler/vertical-pod-autoscaler/deploy/vpa-v1-crd-gen.yaml
update_crds $HOME/go/src/k8s.io/autoscaler/vertical-pod-autoscaler/deploy/vpa-v1-crd.yaml

# self-update
go run cmd/import-crds/main.go --input=$HOME/go/src/kmodules.xyz/resource-metadata/crds
