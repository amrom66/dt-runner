package api

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Model struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ModelSpec `json:"spec"`
}

type ModelSpec struct {
	Jobs    []Job   `json:"jobs"`
	Archive Archive `json:"archive"`
}

type ModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Model `json:"items"`
}

type Job struct {
	Name    string `json:"name,omitempty"`
	Prepare Step   `json:"prepare,omitempty"`
	Check   Step   `json:"check,omitempty"`
	Build   Step   `json:"build,omitempty"`
}

type Step struct {
	Image  string `json:"image,omitempty"`
	Script string `json:"script,omitempty"`
}

type Archive struct {
	Uri    string `json:"uri,omitempty"`
	Key    string `json:"key,omitempty"`
	Secret string `json:"secret,omitempty"`
}
