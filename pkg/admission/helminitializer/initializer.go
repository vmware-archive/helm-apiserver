package helminitializer

import (
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/admission/initializer"
	"k8s.io/client-go/kubernetes"
)

type pluginInitializer struct {
	client kubernetes.Interface
}

var _ admission.PluginInitializer = pluginInitializer{}

func New(client kubernetes.Interface) (pluginInitializer, error) {
	return pluginInitializer{
		client: client,
	}, nil
}

func (i pluginInitializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(initializer.WantsExternalKubeClientSet); ok {
		wants.SetExternalKubeClientSet(i.client)
	}
}
