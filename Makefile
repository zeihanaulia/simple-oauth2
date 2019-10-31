DOCKER_NAMESPACE?=zeihanaulia
DOCKER_TAG?=latest

all: ## run all service
	go run main.go all

lint: ## Lint the files
	@golint -set_exit_status ./...

dep: ## Get the dependencies
	@go mod tidy

.PHONY: build
build: ## Build go binary
	CGO_ENABLED=0 installsuffix=cgo go build -ldflags="-s -w" -o bin/simple-oauth2-$(GOOS) main.go

.PHONE: docker
docker: ## Build docker
	GOOS=linux $(MAKE) build
	docker build -t $(DOCKER_NAMESPACE)/simple-oauth2:${DOCKER_TAG} .

.PHONY: run
run: ## run docker
	docker run --rm --name simple-oauth2 -p8080-8083:8080-8083 zeihanaulia/simple-oauth2

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'