/*
Copyright 2022.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MyCRDSpec defines the desired state of MyCRD
// TODO(user):用户自定义逻辑
type MyCRDSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DeploymentName string `json:"deploymentName"`
	Replicas       *int32 `json:"replicas"`
}

// MyCRDStatus defines the observed state of MyCRD
// TODO(user):用户自定义逻辑
type MyCRDStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	AvailableReplicas int32 `json:"availableReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MyCRD is the Schema for the mycrds API
// TODO(user):用户自定义逻辑
type MyCRD struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyCRDSpec   `json:"spec,omitempty"`
	Status MyCRDStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyCRDList contains a list of MyCRD
type MyCRDList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyCRD `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyCRD{}, &MyCRDList{})
}
