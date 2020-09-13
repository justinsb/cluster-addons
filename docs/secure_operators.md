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

## Don't use the SubjectAccessReview on metrics

Exposing metrics is indeed suboptimal from a security point of view.  However, many clusters and addons can likely tolerate the risk.  Using the kubernetes ServiceAccountToken has its own risks, because that token is a bearer token.

As such, NetworkPolicy is a nice alternative for security sensitive installations.

## Consider using a single dedicated namespace

If installing a cluster-scoped operator, consider using a dedicated namespace.  We can then also use a Role instead of a ClusterRole to further reduce permissions.

TODO: Should we name them a particular way?  Should we label them a particular way?

## Write to status, not spec

Most operators don't / shouldn't need to write back to the spec, instead they should be writing only to status.

Try to remove the create/update/patch/delete permission on the object itself.
