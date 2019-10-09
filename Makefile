.PHONY: default clean checks test build

export GO111MODULE=on

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

BUILD_DATE := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

default: clean checks test build

test: clean
	go test -v -cover ./...

clean:
	rm -rf dist/ cover.out

build: clean
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -ldflags '-X "github.com/onaci/prometheus-export-traefik-endpoints/cmd.version=${VERSION}" -X "github.com/onaci/prometheus-export-traefik-endpoints/cmd.commit=${SHA}" -X "github.com/onaci/prometheus-export-traefik-endpoints/cmd.date=${BUILD_DATE}"' -o prometheus-export-traefik-endpoints

checks:
	golangci-lint run

doc:
	go run . doc

publish-images:
	#seihon publish -v "$(TAG_NAME)" -v "latest" --image-name onaci/prometheus-export-traefik-endpoints --dry-run=false
	docker build -t onaci/prometheus-export-traefik-endpoints .
	docker push onaci/prometheus-export-traefik-endpoints