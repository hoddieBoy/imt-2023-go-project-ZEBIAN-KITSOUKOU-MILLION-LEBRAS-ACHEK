##@ General
clean: ## clean all go cache of all module of the project
	@go clean -cache -modcache -i -r

build: ## build all go module of the project
	@go build -v -o . ./...

install: ## install all go module of the project
	@go install -v ./...