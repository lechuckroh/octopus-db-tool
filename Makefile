.PHONY: vendor

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
GOBUILD_BUILD_DATE_OPT=-ldflags "-X main.buildDateVersion=`date -u +.b%y%m%d-%H%M%S`"
GOBUILD_OPT=$(GO_VENDOR_OPT) -v $(GOBUILD_BUILD_DATE_OPT)
GOTEST_OPT=$(GO_VENDOR_OPT) -v

TEST_DIR=./...
BINARY=oct
BINARY_WINDOWS=oct.exe
BINARY_LINUX=oct
BINARY_MACOS=oct
BUILD_COMPOSE_FILE=build-compose.yml

# Build
compile:
	@$(ENV_STATIC_BUILD) $(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY)

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

package: package-macos package-linux package-windows

package-windows: compile-windows
#	upx -9 $(BINARY_WINDOWS)
#	zip -m oct-win64.zip $(BINARY_WINDOWS)
	gzip $(BINARY_WINDOWS) && mv $(BINARY_WINDOWS).gz oct-win64.gz
package-linux: compile-linux
#	upx -9 $(BINARY_LINUX)
#	zip -m oct-linux64.zip $(BINARY_LINUX)
	gzip $(BINARY_LINUX) && mv $(BINARY_LINUX).gz oct-linux64.gz
package-macos: compile-macos
#	zip -m oct-darwin64.zip $(BINARY_MACOS)
	gzip $(BINARY_MACOS) && mv $(BINARY_MACOS).gz oct-darwin64.gz

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

# Bump version
# see https://github.com/Shyp/bump_version
bump-patch:
	bump_version patch main.go

# local build with drone
drone-exec:
	drone exec --secret-file secrets.txt
drone-exec-tag:
	DRONE_BUILD_EVENT=tag \
	DRONE_REPO_OWNER=lechuckroh \
	DRONE_REPO_NAME=octopus-db-tool \
	DRONE_COMMIT_REF=1-test \
	drone exec --secret-file secrets.txt --event tag
