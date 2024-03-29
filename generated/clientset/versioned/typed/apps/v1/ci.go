/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	v1 "dt-runner/api/apps/v1"
	scheme "dt-runner/generated/clientset/versioned/scheme"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CisGetter has a method to return a CiInterface.
// A group's client should implement this interface.
type CisGetter interface {
	Cis(namespace string) CiInterface
}

// CiInterface has methods to work with Ci resources.
type CiInterface interface {
	Create(ctx context.Context, ci *v1.Ci, opts metav1.CreateOptions) (*v1.Ci, error)
	Update(ctx context.Context, ci *v1.Ci, opts metav1.UpdateOptions) (*v1.Ci, error)
	UpdateStatus(ctx context.Context, ci *v1.Ci, opts metav1.UpdateOptions) (*v1.Ci, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Ci, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.CiList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Ci, err error)
	CiExpansion
}

// cis implements CiInterface
type cis struct {
	client rest.Interface
	ns     string
}

// newCis returns a Cis
func newCis(c *AppsV1Client, namespace string) *cis {
	return &cis{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the ci, and returns the corresponding ci object, and an error if there is any.
func (c *cis) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Ci, err error) {
	result = &v1.Ci{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cis").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Cis that match those selectors.
func (c *cis) List(ctx context.Context, opts metav1.ListOptions) (result *v1.CiList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.CiList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested cis.
func (c *cis) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("cis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a ci and creates it.  Returns the server's representation of the ci, and an error, if there is any.
func (c *cis) Create(ctx context.Context, ci *v1.Ci, opts metav1.CreateOptions) (result *v1.Ci, err error) {
	result = &v1.Ci{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("cis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ci).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a ci and updates it. Returns the server's representation of the ci, and an error, if there is any.
func (c *cis) Update(ctx context.Context, ci *v1.Ci, opts metav1.UpdateOptions) (result *v1.Ci, err error) {
	result = &v1.Ci{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("cis").
		Name(ci.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ci).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *cis) UpdateStatus(ctx context.Context, ci *v1.Ci, opts metav1.UpdateOptions) (result *v1.Ci, err error) {
	result = &v1.Ci{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("cis").
		Name(ci.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ci).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the ci and deletes it. Returns an error if one occurs.
func (c *cis) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cis").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *cis) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cis").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched ci.
func (c *cis) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Ci, err error) {
	result = &v1.Ci{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("cis").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
