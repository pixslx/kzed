/*
Copyright 2024.

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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SequentialDataSetSpec defines the desired state of SequentialDataSet
type SequentialDataSetSpec struct {

	// The block size for the data set (for example, 6160). Default value: 6160.
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	BlockSize int `json:"block-size,omitempty"`

	// The block size for the data set (for example, 6160)
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	DataClass string `json:"data-class,omitempty"`

	// The device type, also known as 'unit'
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	DeviceType string `json:"device-type,omitempty"`

	// The number of directory blocks (for example, 25). Default value: 5.
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	DirectoryBlocks int `json:"directory-blocks,omitempty"`

	// The SMS management class to use for the allocation.
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	ManagementClass string `json:"management-class,omitempty"`

	// The primary space allocation (for example, 5). Default value: 1.
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	PrimarySpace int `json:"primary-space,omitempty"`

	// The record format for the data set (for example, FB for "Fixed Block"). Default value: FB.
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	RecordFormat string `json:"record-format,omitempty"`

	// The logical record length. Analogous to the length of a line (for example, 80). Default value: 80.
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	RecordLength int `json:"record-length,omitempty"`

	// The secondary space allocation (for example, 1).
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	SecondarySpace int `json:"secondary-space,omitempty"`

	// The size of the data set (specified as nCYL or nTRK - where n is the number of cylinders or tracks). Sets the primary allocation (the secondary allocation becomes ~10% of the primary).
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	Size string `json:"size,omitempty"`

	// The SMS storage class to use for the allocation.
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	StorageClass string `json:"storage-class,omitempty"`

	// The volume serial (VOLSER) on which you want the data set to be placed. A VOLSER	is analogous to a drive name on a PC.
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	VolumeSerial string `json:"volume-serial,omitempty"`

	// Name of an existing data set to base your new data set's properties on
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	Like string `json:"like,omitempty"`
}

// SequentialDataSetStatus defines the observed state of SequentialDataSet
type SequentialDataSetStatus struct {
	Status     string      `json:"status,omitempty"`
	CreatedAt  metav1.Time `json:"createdAt,omitempty" protobuf:"bytes,8,opt,name=creationTimestamp"`
	LastSyncAt metav1.Time `json:"lastSyncAt,omitempty" protobuf:"bytes,8,opt,name=creationTimestamp"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SequentialDataSet is the Schema for the sequentialdatasets API
type SequentialDataSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SequentialDataSetSpec   `json:"spec,omitempty"`
	Data   string                  `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`
	Status SequentialDataSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SequentialDataSetList contains a list of SequentialDataSet
type SequentialDataSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SequentialDataSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SequentialDataSet{}, &SequentialDataSetList{})
}
