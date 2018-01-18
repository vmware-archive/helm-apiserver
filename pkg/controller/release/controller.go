
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


package release

import (
	"log"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"

	"github.com/bitnami-labs/helm-apiserver/pkg/apis/helm/v1alpha1"
	"github.com/bitnami-labs/helm-apiserver/pkg/controller/sharedinformers"
	listers "github.com/bitnami-labs/helm-apiserver/pkg/client/listers_generated/helm/v1alpha1"
)

// +controller:group=helm,version=v1alpha1,kind=Release,resource=releases
type ReleaseControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Release
	lister listers.ReleaseLister
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *ReleaseControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing releases labels
	c.lister = arguments.GetSharedInformers().Factory.Helm().V1alpha1().Releases().Lister()
}

// Reconcile handles enqueued messages
func (c *ReleaseControllerImpl) Reconcile(u *v1alpha1.Release) error {
	// Implement controller logic here
	log.Printf("Running reconcile Release for %s\n", u.Name)
	return nil
}

func (c *ReleaseControllerImpl) Get(namespace, name string) (*v1alpha1.Release, error) {
	return c.lister.Releases(namespace).Get(name)
}
