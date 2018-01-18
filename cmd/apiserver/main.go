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

package main

import (
	"github.com/bitnami-labs/helm-apiserver/pkg/apis"
	"github.com/bitnami-labs/helm-apiserver/pkg/cmd/server"
	"github.com/bitnami-labs/helm-apiserver/pkg/openapi"

	// Make sure glide gets these dependencies
	_ "github.com/go-openapi/loads"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"

	// Enable cloud provider auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	version := "v0"
	server.StartApiServer("/registry/bitnami.com", apis.GetAllApiBuilders(), openapi.GetOpenAPIDefinitions, "Helm apiserver", version)
}
