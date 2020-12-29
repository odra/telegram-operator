/*
Copyright 2020.

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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BotMessageSpec defines the desired state of BotMessage
type BotMessageSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Text   string             `json:"text"`
	Image  string             `json:"image"`
	Secret v1.SecretReference `json:"secret"`
}

type BotMessageConditionType string

var (
	BotMessageSending BotMessageConditionType = "Sending"
	BotMessageSent    BotMessageConditionType = "Sent"
	BotMessageError   BotMessageConditionType = "Error"
	BotMessageNew     BotMessageConditionType = ""
)

// BotMessageStatus defines the observed state of BotMessage
type BotMessageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Type    BotMessageConditionType `json:"type"`
	Status  metav1.ConditionStatus  `json:"status"`
	Message string                  `json:"message"`
	Reason  string                  `json:"reason"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// BotMessage is the Schema for the botmessages API
type BotMessage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BotMessageSpec   `json:"spec,omitempty"`
	Status BotMessageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BotMessageList contains a list of BotMessage
type BotMessageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BotMessage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BotMessage{}, &BotMessageList{})
}
