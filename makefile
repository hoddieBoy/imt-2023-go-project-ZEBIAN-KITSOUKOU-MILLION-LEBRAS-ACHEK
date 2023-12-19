.PHONY: help clean build install

##@ Help
help: ## Display this current help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ General
clean: ## clean all go cache of all module of the project
	@go clean -cache -modcache -i -r

build: ## build all go module of the project
	@go build -v -o . ./...

install: ## install all go module of the project
	@go install -v ./...
