package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	addonv1alpha1 "sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/apis/v1alpha1"
)

// ClusterManagedApplicationSpec defines the desired state of ClusterManagedApplication
type ClusterManagedApplicationSpec struct {
	addonv1alpha1.CommonSpec `json:",inline"`
	addonv1alpha1.PatchSpec  `json:",inline"`

	Package string `json:"package,omitempty"`
}

// ClusterManagedApplicationStatus defines the observed state of ClusterManagedApplication
type ClusterManagedApplicationStatus struct {
	addonv1alpha1.CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ClusterManagedApplication is the Schema for the clustermanagedapplications API
type ClusterManagedApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterManagedApplicationSpec   `json:"spec,omitempty"`
	Status ClusterManagedApplicationStatus `json:"status,omitempty"`
}

var _ addonv1alpha1.CommonObject = &ClusterManagedApplication{}

func (o *ClusterManagedApplication) ComponentName() string {
	return o.Spec.Package
}

func (o *ClusterManagedApplication) CommonSpec() addonv1alpha1.CommonSpec {
	return o.Spec.CommonSpec
}

func (o *ClusterManagedApplication) PatchSpec() addonv1alpha1.PatchSpec {
	return o.Spec.PatchSpec
}

func (o *ClusterManagedApplication) GetCommonStatus() addonv1alpha1.CommonStatus {
	return o.Status.CommonStatus
}

func (o *ClusterManagedApplication) SetCommonStatus(s addonv1alpha1.CommonStatus) {
	o.Status.CommonStatus = s
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// ClusterManagedApplicationList contains a list of ClusterManagedApplication
type ClusterManagedApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterManagedApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterManagedApplication{}, &ClusterManagedApplicationList{})
}
