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

// CiliumNetworkingSpec defines the desired state of CiliumNetworking
type CiliumNetworkingSpec struct {
	addonv1alpha1.CommonSpec `json:",inline"`
	addonv1alpha1.PatchSpec  `json:",inline"`

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// CiliumNetworkingStatus defines the observed state of CiliumNetworking
type CiliumNetworkingStatus struct {
	addonv1alpha1.CommonStatus `json:",inline"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
 // +kubebuilder:resource:scope=Cluster 

// CiliumNetworking is the Schema for the ciliumnetworkings API
type CiliumNetworking struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CiliumNetworkingSpec   `json:"spec,omitempty"`
	Status CiliumNetworkingStatus `json:"status,omitempty"`
}

var _ addonv1alpha1.CommonObject = &CiliumNetworking{}

func (o *CiliumNetworking) ComponentName() string {
	return "ciliumnetworking"
}

func (o *CiliumNetworking) CommonSpec() addonv1alpha1.CommonSpec {
	return o.Spec.CommonSpec
}

func (o *CiliumNetworking) PatchSpec() addonv1alpha1.PatchSpec {
	return o.Spec.PatchSpec
}

func (o *CiliumNetworking) GetCommonStatus() addonv1alpha1.CommonStatus {
	return o.Status.CommonStatus
}

func (o *CiliumNetworking) SetCommonStatus(s addonv1alpha1.CommonStatus) {
	o.Status.CommonStatus = s
}

// +kubebuilder:object:root=true
 // +kubebuilder:resource:scope=Cluster 

// CiliumNetworkingList contains a list of CiliumNetworking
type CiliumNetworkingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CiliumNetworking `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CiliumNetworking{}, &CiliumNetworkingList{})
}
