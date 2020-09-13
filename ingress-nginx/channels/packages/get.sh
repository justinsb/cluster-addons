#!/bin/bash

set -e

VERSION=2.15.0

MANIFEST=`curl https://raw.githubusercontent.com/kubernetes/ingress-nginx/ingress-nginx-${VERSION}/deploy/static/provider/cloud/deploy.yaml`

echo "$MANIFEST" \
  | kaml remove-labels helm.sh/chart app.kubernetes.io/managed-by \
  | kaml add-labels app.kubernetes.io/version=${VERSION} \
  | kaml filter -v --kind=Namespace --kind=ServiceAccount --kind=ClusterRole --kind=ClusterRoleBinding --kind=Role --kind=RoleBinding  \
  | kaml filter -v --kind=ValidatingWebhookConfiguration --kind=Job  \
  | kaml filter -v --kind=Service --name=ingress-nginx-controller-admission \
  | kaml remove-volume webhook-cert \
  | sed -e 's@- --validating-webhook@#- --validating-webhook@g' \
  | cat > ingressnginx/${VERSION}/manifest.yaml 

echo "$MANIFEST" \
  | kaml remove-labels helm.sh/chart app.kubernetes.io/managed-by \
  | kaml add-labels app.kubernetes.io/version=${VERSION} \
  | kaml filter --kind=Namespace --kind=ServiceAccount --kind=ClusterRole --kind=ClusterRoleBinding --kind=Role --kind=RoleBinding  \
  | kaml filter -v --kind=Role --kind=RoleBinding --kind=ClusterRole --kind=ClusterRoleBinding --kind=ServiceAccount --name=ingress-nginx-admission  \
  | cat > ingressnginx-rbac/${VERSION}/split_rbac.yaml 
