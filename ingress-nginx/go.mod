module github.com/kubernetes-sigs/cluster-addons

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	golang.org/x/tools v0.0.0-20200910143807-b484961fa2c7 // indirect
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.2
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20200816135617-dbfe418e405f
)
