FROM golang:1.9 as gobuild
WORKDIR /go/src/github.com/bitnami-labs/helm-apiserver
COPY . .
RUN go install ./cmd/...

FROM bitnami/minideb:stretch
MAINTAINER Angus Lees <gus@bitnami.com>
RUN install_packages ca-certificates
COPY --from=gobuild /go/bin/apiserver /go/bin/controller-manager /usr/bin/
EXPOSE 443
CMD ["apiserver"]
