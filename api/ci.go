package api

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Ci struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CiSpec `json:"spec"`
}

type CiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ci `json:"items"`
}

type CiSpec struct {
	Model  string `json:"model"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	On     On     `json:"on"`
}

type On struct {
	Schedule string  `json:"schedule"`
	Events   []Event `json:"events"`
}
type Event struct {
	Push   string `json:"push"`
	Commit string `json:"commit"`
}
