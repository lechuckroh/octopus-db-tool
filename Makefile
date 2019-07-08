GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOCOVER=$(GOCMD) tool cover
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

ENV_GOMOD_ON=GO111MODULE=on
ENV_STATIC_BUILD=CGO_ENABLED=0
GO_VENDOR_OPT=-mod=vendor
GOBUILD_OPT=$(GO_VENDOR_OPT) -v
GOTEST_OPT=$(GO_VENDOR_OPT) -v

TEST_DIR=./...
BINARY=oct
BINARY_WINDOWS=oct.exe
BINARY_LINUX=oct-linux
BINARY_MACOS=oct-macos
BUILD_COMPOSE_FILE=build-compose.yml

# Build
compile:
	@$(ENV_STATIC_BUILD) $(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY)

compile-all: compile-windows compile-linux compile-macos

compile-windows:
	@$(ENV_STATIC_BUILD) GOOS=windows GOARCH=amd64 $(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY_WINDOWS)
compile-linux:
	@$(ENV_STATIC_BUILD) GOOS=linux GOARCH=amd64 $(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY_LINUX)
compile-macos:
	@$(ENV_STATIC_BUILD) GOOS=darwin GOARCH=amd64 $(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY_MACOS)

compile-docker:
	@USER_NAME=`id -un` GROUP_NAME=`id -gn` docker-compose -f $(BUILD_COMPOSE_FILE) run --rm compile

compile-rmi:
	@USER_NAME=`id -un` GROUP_NAME=`id -gn` docker-compose -f $(BUILD_COMPOSE_FILE) down --rmi local || true

# Test
test:
	@$(ENV_GOMOD_ON) $(GOTEST) $(GOTEST_OPT) -count=1 $(TEST_DIR)

# Clean
clean:
	@$(GOCLEAN)
	@rm -f $(BINARY)

# Run
run: compile
	@./$(BINARY)

# Install dependencies to vendor/
vendor:
	@$(GOMOD) vendor
vendor-update:
	@$(GOGET) -u
