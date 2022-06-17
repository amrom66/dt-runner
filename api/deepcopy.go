package api

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *Model) DeepCopyInto(out *Model) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = ModelSpec{
		Tasks:     in.Spec.Tasks,
		Variables: in.Spec.Variables,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Model) DeepCopyObject() runtime.Object {
	out := Model{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *ModelList) DeepCopyObject() runtime.Object {
	out := ModelList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Model, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *Ci) DeepCopyInto(out *Ci) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = CiSpec{
		Model:     in.Spec.Model,
		Repo:      in.Spec.Repo,
		Branch:    in.Spec.Branch,
		On:        in.Spec.On,
		Variables: in.Spec.Variables,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Ci) DeepCopyObject() runtime.Object {
	out := Ci{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *CiList) DeepCopyObject() runtime.Object {
	out := CiList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Ci, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
