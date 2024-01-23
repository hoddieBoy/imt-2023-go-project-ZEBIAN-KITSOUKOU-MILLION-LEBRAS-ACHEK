.PHONY: help clean build install lint lint-fix
BASE_LINT_CMD=docker run -t --rm -v $$(pwd):/app -w /app golangci/golangci-lint golangci-lint run -v

##@ Help
help: ## Display this current help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ General
init: ## init the project
	@go mod tidy
	@docker compose up -d

clean: ## clean all go cache of all module of the project
	@go clean -cache -modcache -i -r

build: ## build all go module of the project
	@go build -v -o . ./...

install: ## install all go module of the project
	@go install -v ./...

lint: ## run golangci-lint on the project
	@$(BASE_LINT_CMD)

lint-fix: ## run golangci-lint on the project and fix issues automatically
	@$(BASE_LINT_CMD) --fix

run: ## run all main application of the project
	@./temperature config/temperature.yaml
	@./humidity config/humidity.yaml
	@./pressure config/pressure.yaml
	@./wind config/wind.yaml
	@./storage config/storage.yaml
	@./alert-manager config/alert-manager.yaml
	@./service

test: ## run all test of the project
	@go test -v ./test/...