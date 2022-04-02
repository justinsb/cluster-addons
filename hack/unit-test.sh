#!/bin/bash

# Copyright 2020 The Kubernetes Authors.
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

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE}")/..
cd "${REPO_ROOT}"

export GO111MODULE=on

echo "****** Testing CoreDNS Operator ******"
make test -C coredns

echo "****** Testing Dashboard Operator ******"
make test -C dashboard

echo "****** Testing Flannel Operator ******"
make test -C flannel

echo "****** Testing Kube-Proxy Operator ******"
KUBERNETES_SERVICE_HOST= make test -C kubeproxy

echo "****** Testing Metrics-server Operator ******"
make test -C metrics-server

echo "****** Testing Bootstrap Utility ******"
make test -C bootstrap

echo "****** Testing Node-Local-DNS Operator ******"
make test -C nodelocaldns

echo "****** Testing tools/rbac-gen ******"
pushd tools/rbac-gen
go test ./...
popd
