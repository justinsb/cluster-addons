module sigs.k8s.io/cluster-addons/coredns

go 1.13

require (
	github.com/coredns/corefile-migration v1.0.10
	github.com/go-logr/logr v0.2.1
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.19.2
	sigs.k8s.io/controller-runtime v0.7.0-alpha.5
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20200512162422-ce639cbf6d4c
	sigs.k8s.io/kustomize/api v0.3.2
)

replace sigs.k8s.io/kubebuilder-declarative-pattern => github.com/justinsb/kubebuilder-declarative-pattern v0.0.0-20201207135810-eefed5cfa7d8

//replace sigs.k8s.io/kubebuilder-declarative-pattern => ../../kubebuilder-declarative-pattern
