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

package v1alpha1

import (
	"log"

	authzv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/endpoints/request"

	"github.com/bitnami-labs/helm-apiserver/pkg/apis/helm"
)

const (
	// StableRepoUrl is the URL of the "stable" chart repo
	StableRepoUrl = "https://kubernetes-charts.storage.googleapis.com"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Release
// +k8s:openapi-gen=true
// +resource:path=releases,strategy=ReleaseStrategy
type Release struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ReleaseSpec `json:"spec,omitempty"`
	// +optional
	Status ReleaseStatus `json:"status,omitempty"`
}

// ReleaseSpec defines the desired state of Release
type ReleaseSpec struct {
	// RepoURL is the URL of the repository. Defaults to stable repo.
	RepoURL string `json:"repoUrl,omitempty"`
	// ChartName is the name of the chart within the repo
	ChartName string `json:"chartName,omitempty"`
	// Version is the chart version
	Version string `json:"version,omitempty"`
	// Values is a string containing (unparsed) YAML values
	Values string `json:"values,omitempty"`
	// Impersonate is the user to impersonate when manipulating resources
	Impersonate ImpersonateSpec `json:"impersonate,omitempty"`
}

type ImpersonateSpec struct {
	User               string                        `json:"user,omitempty"`
	Groups             []string                      `json:"groups,omitempty"`
	Extra              map[string]authzv1.ExtraValue `json:"extra,omitempty"`
	ServiceAccountName string                        `json:"serviceAccountName,omitempty"`
}

// ReleaseStatus defines the observed state of Release
type ReleaseStatus struct {
	Phase string `json:"phase,omitempty"`
	Notes string `json:"notes,omitempty"`
}

// Validate checks that an instance of Release is well formed
func (ReleaseStrategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*helm.Release)
	log.Printf("Validating fields for Release %s\n", o.Name)
	errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	return errors
}

// DefaultingFunction sets default Release field values
func (ReleaseSchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*Release)
	// set default field values here
	log.Printf("Defaulting fields for Release %s\n", obj.Name)
	if obj.Spec.RepoURL == "" {
		obj.Spec.RepoURL = StableRepoUrl
	}
}
