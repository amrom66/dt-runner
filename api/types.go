package api

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Ci struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CiSpec   `json:"spec"`
	Status CiStatus `json:"status,omitempty"`
}

type CiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ci `json:"items"`
}

// Model should be replaced with ModelTemplate
type CiSpec struct {
	Model     string            `json:"model,omitempty"`
	Repo      string            `json:"repo,omitempty"`
	Branch    string            `json:"branch,omitempty"`
	On        On                `json:"on,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

type CiStatus struct {
	Histroy []string `json:"history,omitempty"`
}

type On struct {
	Schedule string   `json:"schedule"`
	Events   []string `json:"events"`
}

type Model struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelSpec   `json:"spec"`
	Status ModelStatus `json:"status,omitempty"`
}

type ModelSpec struct {
	Tasks     []Task            `json:"tasks,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

type ModelStatus struct {
	StartTime  string
	Completime string
	Succeeded  bool
}

type ModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Model `json:"items"`
}

type Task struct {
	Name    string   `json:"name,omitempty"`
	Image   string   `json:"image,omitempty"`
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}
