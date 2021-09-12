module sigs.k8s.io/cluster-addons/coredns

go 1.16

require (
	github.com/coredns/corefile-migration v1.0.14
	github.com/go-logr/logr v0.4.0
	github.com/pkg/errors v0.9.1
	golang.org/x/sys v0.0.0-20210906170528-6f6e22806c34 // indirect
	golang.org/x/tools v0.1.5 // indirect
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20210907124116-1249b4b381fc
	sigs.k8s.io/kustomize/kyaml v0.10.17
)
