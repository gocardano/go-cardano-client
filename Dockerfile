FROM golang:1.15.2 as builder
ARG VERSION
ENV VERSION=${VERSION:-snapshot}
WORKDIR /go/src/github.com/masterjk/go-cardano-client
COPY . /go/src/github.com/masterjk/go-cardano-client

RUN cd /go/src/github.com/masterjk/go-cardano-client \
  && GO111MODULE=off CGO_ENABLED=0 go install -ldflags "-X main.version=${VERSION}" github.com/masterjk/go-cardano-client/cmd/cli

FROM alpine
COPY --from=builder /go/bin/cli /go-cardano-client

ENTRYPOINT ["/go-cardano-client"]
