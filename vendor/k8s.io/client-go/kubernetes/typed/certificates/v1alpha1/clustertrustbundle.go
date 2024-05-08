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

package v1alpha1

import (
	"context"
	json "encoding/json"
	"time"

	"fmt"

	v1alpha1 "k8s.io/api/certificates/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	certificatesv1alpha1 "k8s.io/client-go/applyconfigurations/certificates/v1alpha1"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

// ClusterTrustBundlesGetter has a method to return a ClusterTrustBundleInterface.
// A group's client should implement this interface.
type ClusterTrustBundlesGetter interface {
	ClusterTrustBundles() ClusterTrustBundleInterface
}

// ClusterTrustBundleInterface has methods to work with ClusterTrustBundle resources.
type ClusterTrustBundleInterface interface {
	Create(ctx context.Context, clusterTrustBundle *v1alpha1.ClusterTrustBundle, opts v1.CreateOptions) (*v1alpha1.ClusterTrustBundle, error)
	Update(ctx context.Context, clusterTrustBundle *v1alpha1.ClusterTrustBundle, opts v1.UpdateOptions) (*v1alpha1.ClusterTrustBundle, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ClusterTrustBundle, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ClusterTrustBundleList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ClusterTrustBundle, err error)
	Apply(ctx context.Context, clusterTrustBundle *certificatesv1alpha1.ClusterTrustBundleApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.ClusterTrustBundle, err error)
	ClusterTrustBundleExpansion
}

// clusterTrustBundles implements ClusterTrustBundleInterface
type clusterTrustBundles struct {
	client rest.Interface
}

// newClusterTrustBundles returns a ClusterTrustBundles
func newClusterTrustBundles(c *CertificatesV1alpha1Client) *clusterTrustBundles {
	return &clusterTrustBundles{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterTrustBundle, and returns the corresponding clusterTrustBundle object, and an error if there is any.
func (c *clusterTrustBundles) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ClusterTrustBundle, err error) {
	result = &v1alpha1.ClusterTrustBundle{}
	err = c.client.Get().
		Resource("clustertrustbundles").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterTrustBundles that match those selectors.
func (c *clusterTrustBundles) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ClusterTrustBundleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ClusterTrustBundleList{}
	err = c.client.Get().
		Resource("clustertrustbundles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterTrustBundles.
func (c *clusterTrustBundles) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("clustertrustbundles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a clusterTrustBundle and creates it.  Returns the server's representation of the clusterTrustBundle, and an error, if there is any.
func (c *clusterTrustBundles) Create(ctx context.Context, clusterTrustBundle *v1alpha1.ClusterTrustBundle, opts v1.CreateOptions) (result *v1alpha1.ClusterTrustBundle, err error) {
	result = &v1alpha1.ClusterTrustBundle{}
	err = c.client.Post().
		Resource("clustertrustbundles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterTrustBundle).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a clusterTrustBundle and updates it. Returns the server's representation of the clusterTrustBundle, and an error, if there is any.
func (c *clusterTrustBundles) Update(ctx context.Context, clusterTrustBundle *v1alpha1.ClusterTrustBundle, opts v1.UpdateOptions) (result *v1alpha1.ClusterTrustBundle, err error) {
	result = &v1alpha1.ClusterTrustBundle{}
	err = c.client.Put().
		Resource("clustertrustbundles").
		Name(clusterTrustBundle.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterTrustBundle).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the clusterTrustBundle and deletes it. Returns an error if one occurs.
func (c *clusterTrustBundles) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("clustertrustbundles").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterTrustBundles) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("clustertrustbundles").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched clusterTrustBundle.
func (c *clusterTrustBundles) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ClusterTrustBundle, err error) {
	result = &v1alpha1.ClusterTrustBundle{}
	err = c.client.Patch(pt).
		Resource("clustertrustbundles").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied clusterTrustBundle.
func (c *clusterTrustBundles) Apply(ctx context.Context, clusterTrustBundle *certificatesv1alpha1.ClusterTrustBundleApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.ClusterTrustBundle, err error) {
	if clusterTrustBundle == nil {
		return nil, fmt.Errorf("clusterTrustBundle provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(clusterTrustBundle)
	if err != nil {
		return nil, err
	}
	name := clusterTrustBundle.Name
	if name == nil {
		return nil, fmt.Errorf("clusterTrustBundle.Name must be provided to Apply")
	}
	result = &v1alpha1.ClusterTrustBundle{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("clustertrustbundles").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
