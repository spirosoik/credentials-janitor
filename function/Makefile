.PHONY: build lint setup

OUTPUT:=bin
EXECUTABLE ?= janitor

all: build package 

build: ## Build for development only
	@echo "Building binary..."
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./$(OUTPUT}/$(EXECUTABLE) ./

clean: ## Clean
	@echo "Cleaning..."
	@go clean && rm -rf bin/ && rm -f $(EXECUTABLE).zip

lint: setup ## Linting the codebase
	@echo "Linting..."
	golint -set_exit_status ./...

package:
	@echo "Packing binary..."
	@zip $(EXECUTABLE).zip ./$(OUTPUT)/$(EXECUTABLE)

race: ## Run data race detector
	@echo "Race detecting..."
	@go test -race -short ./...

setup: ## Setup modules
	@echo "Lint installing..."
	@go get -u golang.org/x/lint/golint