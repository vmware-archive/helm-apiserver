package impersonate_test

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/admission"
	clienttesting "k8s.io/client-go/testing"

	helmv1alpha1 "github.com/bitnami-labs/helm-apiserver/pkg/apis/helm/v1alpha1"
)

func TestVerifyImpersonatorPlugin(t *testing.T) {
	var scenarios = []struct {
		admissionInput         helmv1alpha1.Release
		admissionInputKind     schema.GroupVersionKind
		admissionInputResource schema.GroupVersionResource
		admissionMustFail      bool
	}{
		// Empty Impersonate
		{
			admissionInput:         helmv1alpha1.Release{},
			admissionInputKind:     helmv1alpha1.SchemeGroupVersion.WithKind("Release"),
			admissionInputResource: helmv1alpha1.Resource("releases"),
			admissionMustFail:      true,
		},
		// Impersonate as self
		{
			admissionInput:         helmv1alpha1.Release{},
			admissionInputKind:     helmv1alpha1.SchemeGroupVersion.WithKind("Release"),
			admissionInputResource: helmv1alpha1.Resource("releases"),
			admissionMustFail:      false,
		},
		// Impersonate as allowed impersonatee
		{
			admissionInput:         helmv1alpha1.Release{},
			admissionInputKind:     helmv1alpha1.SchemeGroupVersion.WithKind("Release"),
			admissionInputResource: helmv1alpha1.Resource("releases"),
			admissionMustFail:      false,
		},
		// Impersonate as disallowed impersonatee
		{
			admissionInput:         helmv1alpha1.Release{},
			admissionInputKind:     helmv1alpha1.SchemeGroupVersion.WithKind("Release"),
			admissionInputResource: helmv1alpha1.Resource("releases"),
			admissionMustFail:      true,
		},
	}

	for idx, scenario := range scenarios {
		// prepare
		client := &fake.Clientset{}
		client.AddReactor("get", "subjectaccessreview", func(action clienttesting.Action) (bool, runtime.Object, error) {
			return true, &scenarios.sarOutput, nil
		})

		target, err := impersonate.New()
		if err != nil {
			t.Fatalf("scenario %d: failed to create impersonate admission plugin due to %v", idx, err)
		}

		init, err := helminitializer.New()
		if err != nil {
			t.Fatalf("scenario %d: failed to create helm plugin initializer due to %v", idx, err)
		}
		init.Initialize(target)

		if err := admission.Validate(target); err != nil {
			t.Fatalf("scenario %d: failed to initialize helm admission plugin due to %v", idx, err)
		}

		// act
		err = target.Admit(admission.NewAttributesRecord(
			&scenario.admissionInput,
			nil,
			scenario.admissionInputKind,
			scenario.admissionInput.ObjectMeta.Namespace,
			"",
			scenario.admissionInputResource,
			"",
			admission.Create,
			nil),
		)

		// validate
		if scenario.admissionMustFail && err == nil {
			t.Errorf("scenario %d: expected an error but got nothing", idx)
		}

		if !scenario.admissionMustFail && err != nil {
			t.Errorf("scenario %d: admission plugin returned unexpected error %v", idx, err)
		}
	}
}
