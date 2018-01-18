package impersonate

import (
	"errors"
	"fmt"
	"io"

	authzv1 "k8s.io/api/authorization/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/initializer"
	"k8s.io/client-go/kubernetes"

	helmv1alpha1 "github.com/bitnami-labs/helm-apiserver/pkg/apis/helm/v1alpha1"
)

const (
	// PluginName is the name of the plugin
	PluginName = "VerifyImpersonator"
)

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register(PluginName, func(config io.Reader) (admission.Interface, error) {
		return New()
	})
}

// Impersonator is an implementation of admission.Interface.
// It verifies the user is allowed to set Release.Spec.Impersonate
type impersonator struct {
	*admission.Handler
	client kubernetes.Interface
}

func (i *impersonator) Admit(a admission.Attributes) error {
	if a.GetSubresource() != "" {
		return nil
	}

	switch o := a.GetObject().(type) {
	case *helmv1alpha1.Release:
		imp := o.Spec.Impersonate

		// Must specify either user xor serviceAccount
		if (imp.User != "" && imp.ServiceAccountName != "") ||
			(imp.User == "" && imp.ServiceAccountName == "") {
			return admission.NewForbidden(a, fmt.Errorf("exactly one of impersonate.user or impersonate.serviceAccountName is required"))
		}

		// Is the API user allowed to impersonate the enclosed user?
		rsrcs := []authzv1.ResourceAttributes{}
		if imp.ServiceAccountName != "" {
			rsrcs = append(rsrcs, authzv1.ResourceAttributes{
				Group:     "",
				Verb:      "impersonate",
				Resource:  "serviceaccounts",
				Namespace: a.GetNamespace(),
				Name:      imp.ServiceAccountName,
				//FIXME?: Name: serviceaccount.MakeUsername(a.GetNamespace(), imp.ServiceAccountName)
			})
		}
		if imp.User != "" {
			rsrcs = append(rsrcs, authzv1.ResourceAttributes{
				Group:    "",
				Verb:     "impersonate",
				Resource: "users",
				Name:     imp.User,
			})
		}
		for _, g := range imp.Groups {
			rsrcs = append(rsrcs, authzv1.ResourceAttributes{
				Group:    "",
				Verb:     "impersonate",
				Resource: "groups",
				Name:     g,
			})
		}
		for k, values := range imp.Extra {
			for _, v := range values {
				rsrcs = append(rsrcs, authzv1.ResourceAttributes{
					Group:       authzv1.SchemeGroupVersion.Group,
					Version:     authzv1.SchemeGroupVersion.Version,
					Verb:        "impersonate",
					Resource:    "userextras",
					Subresource: k,
					Name:        v,
				})
			}
		}

		user := a.GetUserInfo()

		// Dear golang, you know this is a no-op .. right?
		extra := map[string]authzv1.ExtraValue{}
		for k, v := range user.GetExtra() {
			extra[k] = authzv1.ExtraValue(v)
		}

		sar := authzv1.SubjectAccessReview{
			Spec: authzv1.SubjectAccessReviewSpec{
				User:               user.GetName(),
				Groups:             user.GetGroups(),
				UID:                user.GetUID(),
				Extra:              extra,
				ResourceAttributes: nil, // set below
			},
		}

		for _, rsrc := range rsrcs {
			sar.Spec.ResourceAttributes = &rsrc
			res, err := i.client.AuthorizationV1().
				SubjectAccessReviews().
				Create(&sar)

			if err != nil {
				return apierrors.NewInternalError(err)
			}
			if res.Status.EvaluationError != "" {
				return apierrors.NewInternalError(errors.New(res.Status.EvaluationError))
			}
			if !res.Status.Allowed {
				return admission.NewForbidden(a, errors.New(res.Status.Reason))
			}
		}
	default:
		return nil
	}

	return nil
}

var _ = initializer.WantsExternalKubeClientSet(&impersonator{})

func (i *impersonator) SetExternalKubeClientSet(client kubernetes.Interface) {
	i.client = client
}

func (i *impersonator) Validate() error {
	if i.client == nil {
		return fmt.Errorf("missing kubernetes client")
	}
	return nil
}

// New creates a new instance of the admission plugin
func New() (admission.Interface, error) {
	return &impersonator{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}, nil
}
