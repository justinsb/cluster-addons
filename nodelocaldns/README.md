# NodeLocalDNS operator

## Created with:

```bash
export GO111MODULE=on
export KUBEBUILDER_ENABLE_PLUGINS=on

kubebuilder init --domain k8s.io --license apache2 --owner "The Kubernetes authors" --pattern addon

kubebuilder create api --group addons --version v1alpha1 --kind NodeLocalDNS --pattern addon --resource --controller --namespaced

```




We have to do two things:

```bash
  # Replace the sed configurations with variable values.
  sed -i -e "s/__PILLAR__DNS__DOMAIN__/${DNS_DOMAIN}/g" "${localdns_file}"
  sed -i -e "s/__PILLAR__DNS__SERVER__/${DNS_SERVER_IP}/g" "${localdns_file}"
  sed -i -e "s/__PILLAR__LOCAL__DNS__/${LOCAL_DNS_IP}/g" "${localdns_file}"
```

```bash
      if [[ "${ENABLE_NODELOCAL_DNS:-}" == "true" ]]; then
        dns_args=("--cluster-dns=${LOCAL_DNS_IP}" "--cluster-domain=${DNS_DOMAIN}")
      else
        dns_args=("--cluster-dns=${DNS_SERVER_IP}" "--cluster-domain=${DNS_DOMAIN}")
      fi
```

The sed substitutions are straightfoward, the cluster-dns change ... less so.

Instead we can use the approach described here:
https://github.com/kubernetes/dns/pull/280/files , which has since been
simplified by nodelocal-dns itself and is now the default!


We must replace the same 3 variables in the configmap and the daemonset
args / livenessProbe.  So we'll use string substitution technique.  Though this
isn't preferred in general, it is often the best choice for replacing values in
flags or configmaps.  (A better approach is to use CRDs for this configuration -
maybe even the CRD for the addon itself)

So we will need to perform the following substitutions:

`__PILLAR__DNS__DOMAIN__` we could expose this as a field, but most people use
cluster.local.  Also not clear if we should default to cluster.local, or default
to the value we see in `/etc/resolv.conf`.  For now, let's just default to cluster.local.

`__PILLAR__DNS__SERVER__` is the IP for DNS as configured in kubelet; we could expose
this as a field again, but we can also get it by querying for the
`kube-dns` service IP.

`__PILLAR__LOCAL__DNS__` is the IP for the local interception; we could expose
this as a field, but for now we can lock it to 169.254.20.10 (as used in the
kubernetes bash configuration)


Note that some variables are replaced by the nodelocal-dns agent itself:

https://github.com/kubernetes/dns/blob/960a9860b283aaadb9314ff4025cd677f0ec21a1/cmd/node-cache/app/configmap.go#L67-L69

__PILLAR__UPSTREAM__SERVERS__ is read from kube-dns config or from /etc/resolv.conf
__PILLAR__CLUSTER__DNS__ is defaulted from IP for kube-dns-upstream
__PILLAR__LOCAL__DNS__ is replaced from --localip arg


## TODOs

* Should we expose fields for more options?
* Should we have a toleration so this runs on master nodes?
