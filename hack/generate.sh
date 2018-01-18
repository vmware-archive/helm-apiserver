#!/bin/sh

set -e

pkg=github.com/bitnami-labs/helm-apiserver

join() {
    sep="$1"; shift
    IFS="$sep" ret="$*"; unset IFS
    echo "$ret"
}

set -x

apiregister-gen \
    --input-dirs $pkg/pkg/apis/... \
    --input-dirs $pkg/pkg/controller/...

conversion-gen \
    --input-dirs $pkg/pkg/apis/helm/v1alpha1 \
    --input-dirs $pkg/pkg/apis/helm \
    -o $GOPATH/src \
    --go-header-file hack/boilerplate.go.txt \
    -O zz_generated.conversion \
    --extra-peer-dirs k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/conversion,k8s.io/apimachinery/pkg/runtime

deepcopy-gen \
    --input-dirs $pkg/pkg/apis/helm/v1alpha1 \
    --input-dirs $pkg/pkg/apis/helm \
    -o $GOPATH/src \
    --go-header-file hack/boilerplate.go.txt \
    -O zz_generated.deepcopy

apipkgs=$(
    go list -f '{{.ImportComment}}' \
       ./vendor/k8s.io/api/... \
       ./vendor/k8s.io/apimachinery/pkg/apis/meta/v1 \
       ./vendor/k8s.io/apimachinery/pkg/api/resource \
       ./vendor/k8s.io/apimachinery/pkg/version \
       ./vendor/k8s.io/apimachinery/pkg/runtime \
       ./vendor/k8s.io/apimachinery/pkg/util/intstr \
)

openapi-gen \
    --input-dirs $pkg/pkg/apis/helm/v1alpha1 \
    --input-dirs $pkg/pkg/apis/helm \
    -o $GOPATH/src \
    --go-header-file hack/boilerplate.go.txt \
    -i $(join , $apipkgs) \
    --output-package $pkg/pkg/openapi

defaulter-gen \
    --input-dirs $pkg/pkg/apis/helm/v1alpha1 \
    --input-dirs $pkg/pkg/apis/helm \
    -o $GOPATH/src \
    --go-header-file hack/boilerplate.go.txt \
    -O zz_generated.defaults \
    --extra-peer-dirs=k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/conversion,k8s.io/apimachinery/pkg/runtime

client-gen \
    -o /home/gus/go/src \
    --go-header-file hack/boilerplate.go.txt \
    --input-base $pkg/pkg/apis \
    --input helm/v1alpha1 \
    --clientset-path $pkg/pkg/client/clientset_generated \
    --clientset-name clientset

client-gen \
    -o $GOPATH/src \
    --go-header-file hack/boilerplate.go.txt \
    --input-base $pkg/pkg/apis \
    --input helm \
    --clientset-path $pkg/pkg/client/clientset_generated \
    --clientset-name internalclientset

lister-gen \
    --input-dirs $pkg/pkg/apis/helm/v1alpha1 \
    --input-dirs $pkg/pkg/apis/helm \
    -o $GOPATH/src \
    --go-header-file hack/boilerplate.go.txt \
    --output-package $pkg/pkg/client/listers_generated

informer-gen \
    --input-dirs $pkg/pkg/apis/helm/v1alpha1 \
    --input-dirs $pkg/pkg/apis/helm \
    -o $GOPATH/src \
    --go-header-file hack/boilerplate.go.txt \
    --output-package $pkg/pkg/client/informers_generated \
    --listers-package $pkg/pkg/client/listers_generated \
    --versioned-clientset-package $pkg/pkg/client/clientset_generated/clientset \
    --internal-clientset-package $pkg/pkg/client/clientset_generated/internalclientset
