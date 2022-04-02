# Cluster Addons

Cluster Addons is a sub-project of [SIG-Cluster-Lifecycle](https://github.com/kubernetes/community/tree/master/sig-cluster-lifecycle). Addon management has been a problem of cluster tooling for a long time.

This sub-project wants to figure out the best way to install, manage and deliver cluster addons.

In this repository we explore ideas for all of the above. [Cluster addon operators](https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/addons/2492-Addons-via-Operators) in particular.

With cluster addon operators, we are exploring a kubernetes-native way of managing addons using CRDs(Custom Resource Definitions) and controllers where the controllers encode how best to manage the addon. Installing and managing an addon could be as simple as creating a custom resource.

We are also exploring other concepts such as managing different simple addons with a single controller to reduce consumption of resources and make it dead simple to manage addons. We welcome you to create your own addon operators and try using the ones created by this community.

## Frequently asked questions

> What is this?

Born out of the discussion in the [original KEP PR](https://github.com/kubernetes/enhancements/pull/746), we set up the sub-project with the goal to explore addon operators, since then we took on a number of other challenges.

> What is this not?

This sub-project is not interested in maintaining all cluster addons. Here we want to create some design patterns, some libraries, some supporting tooling, so everybody can easily create their own operators.

Not everything will need a cluster addon. Not everyone will want to use an operator.

> What is a cluster addon?

The lifecycle of a cluster addon is managed alongside the lifecycle of the cluster. Typically it has to be upgraded/downgraded when you move to a newer Kubernetes version. We want to use operators for this: a CRD describes the addon, and then the code which installs whatever the addon does, controlled by the CRD.

> How do I build my own cluster addon operator?

We have created a tutorial on how to create your own addon operator [here](https://github.com/kubernetes-sigs/cluster-addons/tree/master/walkthrough.md)

> What's your current agenda and timeline?

We:

- created an [addon operator we deemed as straight-forward](https://github.com/kubernetes-sigs/cluster-addons/tree/master/coredns), so we have actual code to look at.
- wrote an [installer library](https://github.com/kubernetes-sigs/cluster-addons/tree/master/installer) to install addons into the cluster
- added support for [addon operators to kubebuilder](https://github.com/kubernetes-sigs/kubebuilder/tree/master/pkg/plugins/golang/declarative/v1)
- started work on integrating the addon installer into kubeadm and kOps
- are building consensus around a pattern of a "base" manifest with an overlay of parameter-driven changes (usually through code)

> Who does this?

Cluster addons is a community project. If you're interested in building this, please get in touch. We're all ears!

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

Check out up to date information about where discussions and meetings happen on
the [community page of SIG Cluster Lifecycle](https://github.com/kubernetes/community/tree/master/sig-cluster-lifecycle).

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).
