kubebuilder init  --plugins="go/v3,declarative,klogr" --project-version=3 \
  --domain x-k8s.io --repo sigs.k8s.io/cluster-addons/coredns \
  --owner "The Kubernetes Authors" --license apache2

kubebuilder create api --group addons --kind CoreDNS --version v1alpha1
