module sigs.k8s.io/addon-operators/nodelocaldns

go 1.13

replace sigs.k8s.io/kubebuilder-declarative-pattern => ../../../sigs.k8s.io/kubebuilder-declarative-pattern

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v0.0.0-20190918200256-06eb1244587a
	k8s.io/klog v0.3.3
	sigs.k8s.io/controller-runtime v0.3.0
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0
)
