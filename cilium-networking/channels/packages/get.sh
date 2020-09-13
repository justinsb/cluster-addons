#!/bin/bash

set -e

VERSION=1.8.3

MANIFEST=`curl kubectl create -f https://raw.githubusercontent.com/cilium/cilium/v1.8.3/install/kubernetes/quick-install.yaml`

mkdir -p ciliumnetworking/${VERSION}/

echo "$MANIFEST" \
  | kaml filter -v --kind=Namespace --kind=ServiceAccount --kind=ClusterRole --kind=ClusterRoleBinding --kind=Role --kind=RoleBinding  \
  | cat > ciliumnetworking/${VERSION}/manifest.yaml 

mkdir -p ciliumnetworking-rbac/${VERSION}/

echo "$MANIFEST" \
  | kaml filter --kind=Namespace --kind=ServiceAccount --kind=ClusterRole --kind=ClusterRoleBinding --kind=Role --kind=RoleBinding  \
  | cat > ciliumnetworking-rbac/${VERSION}/manifest.yaml 
