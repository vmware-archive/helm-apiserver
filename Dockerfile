FROM golang:1.9-alpine3.7 as gobuild
WORKDIR /go/src/github.com/bitnami-labs/helm-apiserver
COPY . .
RUN go install ./cmd/...

FROM alpine:3.7
MAINTAINER Angus Lees <gus@bitnami.com>
RUN apk --no-cache add ca-certificates
COPY --from=gobuild /go/bin/apiserver /go/bin/controller-manager /usr/bin/
EXPOSE 443
CMD ["apiserver"]
