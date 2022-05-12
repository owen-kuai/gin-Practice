#GOBUILD = CGO_ENABLED=0 GO111MODULE=on go build
SED = sed -i
GIT_SHA=$(shell git rev-parse HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_TIME=$(shell TZ=UTC-8 date -Isecond)

ifeq ($(shell uname), Darwin)
	SED += ""
endif

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

swag: fmt vet
	swag init --generalInfo ./cmd/app/main.go --output ./docs/app_swagger

# Run tests
test: fmt vet
	go test -timeout 120s ./... -coverprofile cover.out

build:swag
	go build -a -ldflags '-extldflags "-static"' -o bin/app ./cmd/app
