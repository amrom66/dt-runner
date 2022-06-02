package api

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Model struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ModelSpec `json:"spec"`
}

type ModelSpec struct {
	Replicas int `json:"replicas"`
}

type ModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Model `json:"items"`
}
