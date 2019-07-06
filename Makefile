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
BUILD_COMPOSE_FILE=build-compose.yml

# Build
compile:
	@$(ENV_STATIC_BUILD) $(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY)

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
