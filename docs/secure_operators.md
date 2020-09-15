# Limiting RBAC scope

## Pre-create the permissions needed by the workload

Previously, we suggested adding the required permissions to the operator role.  However, this gives the operator very broad permissions.

A better approach is to instead pre-create the RBAC permissions needed for installing the workload.  We don't lose flexibility by doing so, because to change these permissions in the old model we would still need to change the operator manifest.

Whatever method you are using to install the operator, that should also install the workload-RBAC manifest.

Note: this doesn't work for per-namespace addons, where we likely want to create a Role, RoleBinding & ServiceAccount per namespace.

## Use a StatefulSet instead of a Deployment

Operators have normally used Deployments with leader-election, which recovers very quickly from failure.

However, it requires additional RBAC permissions (on the lease object) and it consumes more resources, because there is a second container.

Instead, where recovery time is not critical (and where you are not in the critical path for launching a replacement Node / Pod), consider using a StatefulSet instead.

Try using   podManagementPolicy: "Parallel"


## Don't use the SubjectAccessReview on metrics

Exposing metrics is indeed suboptimal from a security point of view.  However, many clusters and addons can likely tolerate the risk.  Using the kubernetes ServiceAccountToken has its own risks, because that token is a bearer token.

As such, NetworkPolicy is a nice alternative for security sensitive installations.

## Consider using a single dedicated namespace

If installing a cluster-scoped operator, consider using a dedicated namespace.  We can then also use a Role instead of a ClusterRole to further reduce permissions.

TODO: Should we name them a particular way?  Should we label them a particular way?

## Write to status, not spec

Most operators don't / shouldn't need to write back to the spec, instead they should be writing only to status.

Try to remove the create/update/patch/delete permission on the object itself.

## Clean up manifests

```
rm -rf config/certmanager
rm -rf config/prometheus
rm -rf config/webhook
rm -rf config/crd/patches/
rm -rf config/crd/kustomizeconfig.yaml

cat <<EOF > config/crd/kustomization.yaml
# This kustomization.yaml is not intended to be run by itself,
# since it depends on service name and namespace that are out of this kustomize package.
# It should be run by config/default
resources:
- bases/addons.x-k8s.io_caliconetworkings.yaml
# +kubebuilder:scaffold:crdkustomizeresource
EOF

rm -rf config/manager/default.yaml

cat <<EOF > config/manager/kustomization.yaml
resources:
- statefulset.yaml
EOF


rm -rf config/rbac/auth_proxy*
rm -rf config/rbac/leader_election*

go run . --yaml ~/k8s/src/sigs.k8s.io/cluster-addons/calico-networking/channels/packages/caliconetworking/3.16.1/manifest.yaml --supervisory

cat <<EOF > config/rbac/kustomization.yaml
resources:
- role.yaml
- role_binding.yaml
- deploy.yaml
- deploy_binding.yaml
EOF


```

