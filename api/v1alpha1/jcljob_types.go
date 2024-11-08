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

// JCLJobSpec defines the desired state of JCLJob
type JCLJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Path of DataSet on the Mainframe with JCL to be submitted (e.g. USERID.SOURCE(SOMEJCL))
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	DSPath string `json:"dsPath,omitempty"`

	// Path of DataSet on the Mainframe with JCL to be submitted (e.g. USERID.SOURCE(SOMEJCL))
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	USSPath string `json:"ussPath,omitempty"`

	// JCL script to be submitted
	//+kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable"
	JCL string `json:"jcl,omitempty"`
}

// JCLJobStatus defines the observed state of JCLJob
type JCLJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	JobID      string             `json:"jobid,omitempty"`
	JobName    string             `json:"jobname,omitempty"`
	Status     string             `json:"status,omitempty"`
	ReturnCode string             `json:"retcode,omitempty"`
	StartedAt  metav1.Time        `json:"startedAt,omitempty" protobuf:"bytes,8,opt,name=creationTimestamp"`
	FinishedAt metav1.Time        `json:"finishedAt,omitempty" protobuf:"bytes,8,opt,name=creationTimestamp"`
	SpoolFiles []JCLJobSpoolFiles `json:"spools,omitempty"`
}

type JCLJobSpoolFiles struct {
	SpoolID  string `json:"spoolid,omitempty"`
	StepName string `json:"stepname,omitempty"`
	DDName   string `json:"ddname,omitempty"`
	Data     string `json:"data,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// JCLJob is the Schema for the jcljobs API
type JCLJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	//+required
	Spec JCLJobSpec `json:"spec,omitempty"`

	Status JCLJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JCLJobList contains a list of JCLJob
type JCLJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JCLJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JCLJob{}, &JCLJobList{})
}
