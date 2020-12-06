#!/bin/bash

set -e
#set -x

kustomize build ../crd

echo ""
echo "---"
echo ""

# LoadRestrictionsNone needed to load RBAC/CRD from parent dir
kustomize build .


#cat ../rbac/role.yaml