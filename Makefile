GO = go
GOPKGS = ./cmd/... ./pkg/...
GOFLAGS =
GOTESTFLAGS = $(GOFLAGS) -race
GOBUILDFLAGS = $(GOFLAGS)

all: bin/apiserver bin/controller-manager

bin/%:
	$(GO) build $(GOBUILDFLAGS) -o $@ ./cmd/$*

bin/apiserver: $(shell hack/godeps.sh ./cmd/apiserver)
bin/controller-manager: $(shell hack/godeps.sh ./cmd/controller-manager)

test:
	$(GO) test $(GOTESTFLAGS) $(GOPKGS)

fmt:
	$(GOFMT) -s -w $(shell $(GO) list -f '{{$d := .Dir}}{{range .GoFiles}}{{$d}}/{{.}} {{end}}' $(GOPKGS))

vet:
	$(GO) vet $(GOPKGS)

generated-files:
	hack/generate.sh

deploy-config:
	# FIXME: Need better installation procedure (cert process)
	apiserver-boot build config --name helm-apiserver --namespace kube-system --image bitnami/helm-apiserver:latest
