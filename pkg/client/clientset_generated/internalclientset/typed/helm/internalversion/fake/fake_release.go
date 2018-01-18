/*
Copyright 2017 The helm-apiserver Authors.

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
package fake

import (
	helm "github.com/bitnami-labs/helm-apiserver/pkg/apis/helm"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeReleases implements ReleaseInterface
type FakeReleases struct {
	Fake *FakeHelm
	ns   string
}

var releasesResource = schema.GroupVersionResource{Group: "helm.bitnami.com", Version: "", Resource: "releases"}

var releasesKind = schema.GroupVersionKind{Group: "helm.bitnami.com", Version: "", Kind: "Release"}

// Get takes name of the release, and returns the corresponding release object, and an error if there is any.
func (c *FakeReleases) Get(name string, options v1.GetOptions) (result *helm.Release, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(releasesResource, c.ns, name), &helm.Release{})

	if obj == nil {
		return nil, err
	}
	return obj.(*helm.Release), err
}

// List takes label and field selectors, and returns the list of Releases that match those selectors.
func (c *FakeReleases) List(opts v1.ListOptions) (result *helm.ReleaseList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(releasesResource, releasesKind, c.ns, opts), &helm.ReleaseList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &helm.ReleaseList{}
	for _, item := range obj.(*helm.ReleaseList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested releases.
func (c *FakeReleases) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(releasesResource, c.ns, opts))

}

// Create takes the representation of a release and creates it.  Returns the server's representation of the release, and an error, if there is any.
func (c *FakeReleases) Create(release *helm.Release) (result *helm.Release, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(releasesResource, c.ns, release), &helm.Release{})

	if obj == nil {
		return nil, err
	}
	return obj.(*helm.Release), err
}

// Update takes the representation of a release and updates it. Returns the server's representation of the release, and an error, if there is any.
func (c *FakeReleases) Update(release *helm.Release) (result *helm.Release, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(releasesResource, c.ns, release), &helm.Release{})

	if obj == nil {
		return nil, err
	}
	return obj.(*helm.Release), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeReleases) UpdateStatus(release *helm.Release) (*helm.Release, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(releasesResource, "status", c.ns, release), &helm.Release{})

	if obj == nil {
		return nil, err
	}
	return obj.(*helm.Release), err
}

// Delete takes name of the release and deletes it. Returns an error if one occurs.
func (c *FakeReleases) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(releasesResource, c.ns, name), &helm.Release{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeReleases) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(releasesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &helm.ReleaseList{})
	return err
}

// Patch applies the patch and returns the patched release.
func (c *FakeReleases) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *helm.Release, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(releasesResource, c.ns, name, data, subresources...), &helm.Release{})

	if obj == nil {
		return nil, err
	}
	return obj.(*helm.Release), err
}
