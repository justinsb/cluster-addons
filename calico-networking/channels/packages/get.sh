#!/bin/bash

set -e

VERSION=3.16.1

MANIFEST=`curl https://docs.projectcalico.org/manifests/calico.yaml`


mkdir -p caliconetworking/${VERSION}/

echo "$MANIFEST" \
  | kaml filter -v --kind=Namespace --kind=ServiceAccount --kind=ClusterRole --kind=ClusterRoleBinding --kind=Role --kind=RoleBinding  \
  | kaml filter -v --kind=CustomResourceDefinition  \
  | cat > caliconetworking/${VERSION}/manifest.yaml 

mkdir -p caliconetworking-rbac/${VERSION}/
echo "$MANIFEST" \
  | kaml filter --kind=Namespace --kind=ServiceAccount --kind=ClusterRole --kind=ClusterRoleBinding --kind=Role --kind=RoleBinding  \
  | kaml filter -v --kind=CustomResourceDefinition  \
  | cat > caliconetworking-rbac/${VERSION}/manifest.yaml 

mkdir -p caliconetworking-crds/${VERSION}/
echo "$MANIFEST" \
  | kaml filter -v --kind=Namespace --kind=ServiceAccount --kind=ClusterRole --kind=ClusterRoleBinding --kind=Role --kind=RoleBinding  \
  | kaml filter --kind=CustomResourceDefinition  \
  | cat > caliconetworking-crds/${VERSION}/manifest.yaml 
