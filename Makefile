PROJECT_NAME=go-cardano-client
PROJECT_SRCDIR=github.com/masterjk/${PROJECT_NAME}

VERSION="0.0.1-$(shell git rev-parse --short=8 HEAD)"
DOCKER_ARGS="--rm -u $(shell id -u) -e GOCACHE=/tmp/"

GOLANG_IMAGE="golang:1.15.2"

.PHONY: default fmt vet test coverage build

default: container

fmt:
	@echo ➭ Running go fmt
	@docker run "${DOCKER_ARGS}" -v `pwd`:/go/src/${PROJECT_SRCDIR} \
		-w /go/src/${PROJECT_SRCDIR} ${GOLANG_IMAGE} \
		go fmt ./... | read 1>&2 && exit 1 || true

vet:
	@echo ➭ Running go vet
	@docker run "${DOCKER_ARGS}" -v `pwd`:/go/src/${PROJECT_SRCDIR} \
		-w /go/src/${PROJECT_SRCDIR} ${GOLANG_IMAGE} go vet ./...

test:
	@echo ➭ Running go test
	@docker run "${DOCKER_ARGS}" \
		-v `pwd`:/go/src/${PROJECT_SRCDIR} \
		-w /go/src/${PROJECT_SRCDIR} ${GOLANG_IMAGE} go test ./...

coverage:
	@echo ➭ Running go test coverage
	@docker run "${DOCKER_ARGS}" \
		-v `pwd`:/go/src/${PROJECT_SRCDIR} \
		-w /go/src/${PROJECT_SRCDIR} ${GOLANG_IMAGE} go test -coverprofile=.cov ./...;  go tool cover -func .cov

build:
	@echo ➭ Building ${PROJECT_NAME}
	@docker run "${DOCKER_ARGS}" \
		-e GOOS=darwin \
		-e GOARCH=amd64 \
		-e GO111MODULE=off \
		-e CGO_ENABLED=0 \
		-v `pwd`:/go/src/${PROJECT_SRCDIR} \
		-w /go/src/${PROJECT_SRCDIR} ${GOLANG_IMAGE} go build \
			-o ${PROJECT_NAME} \
			-ldflags "-X main.version=${VERSION}" \
			github.com/masterjk/go-cardano-client/cmd/cli
