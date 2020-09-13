/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	addonv1alpha1 "sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/apis/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CalicoNetworkingSpec defines the desired state of CalicoNetworking
type CalicoNetworkingSpec struct {
	addonv1alpha1.CommonSpec `json:",inline"`
	addonv1alpha1.PatchSpec  `json:",inline"`

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// CalicoNetworkingStatus defines the observed state of CalicoNetworking
type CalicoNetworkingStatus struct {
	addonv1alpha1.CommonStatus `json:",inline"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
 // +kubebuilder:resource:scope=Cluster 

// CalicoNetworking is the Schema for the caliconetworkings API
type CalicoNetworking struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CalicoNetworkingSpec   `json:"spec,omitempty"`
	Status CalicoNetworkingStatus `json:"status,omitempty"`
}

var _ addonv1alpha1.CommonObject = &CalicoNetworking{}

func (o *CalicoNetworking) ComponentName() string {
	return "caliconetworking"
}

func (o *CalicoNetworking) CommonSpec() addonv1alpha1.CommonSpec {
	return o.Spec.CommonSpec
}

func (o *CalicoNetworking) PatchSpec() addonv1alpha1.PatchSpec {
	return o.Spec.PatchSpec
}

func (o *CalicoNetworking) GetCommonStatus() addonv1alpha1.CommonStatus {
	return o.Status.CommonStatus
}

func (o *CalicoNetworking) SetCommonStatus(s addonv1alpha1.CommonStatus) {
	o.Status.CommonStatus = s
}

// +kubebuilder:object:root=true
 // +kubebuilder:resource:scope=Cluster 

// CalicoNetworkingList contains a list of CalicoNetworking
type CalicoNetworkingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CalicoNetworking `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CalicoNetworking{}, &CalicoNetworkingList{})
}
