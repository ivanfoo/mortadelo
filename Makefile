# Default shell
SHELL := /bin/bash

# Project
PROJECT := mortadelo
COMMANDS := mortadelo

# Go commands
GOCMD = go
GOGET = $(GOCMD) get -v -t
GOTEST = $(GOCMD) test -v

# Coverage
COVERAGE_REPORT = coverage.txt
COVERAGE_PROFILE = profile.out
COVERAGE_MODE = atomic

# Env
BUILD_PATH := $(PWD)

# Artifacts
ARTIFACTS_PATH := $(PWD)/artifacts
BUILD ?= $(shell date +"%m-%d-%Y_%H_%M_%S")
COMMIT ?= $(shell git log --format='%H' -n 1 | cut -c1-10)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
PKG_OS = darwin linux
PKG_ARCH = amd64

# Travis-CI
ifneq ($(origin CI), undefined)
	COMMIT := $(shell echo $(TRAVIS_COMMIT) | cut -c1-10)
	BRANCH := $(TRAVIS_BRANCH)
	BUILD_PATH := $(TRAVIS_BUILD_DIR)
	ARTIFACTS_PATH := $(TRAVIS_BUILD_DIR)/artifacst
endif

dependencies:
	$(GOGET) ./...

test-coverage:
	cd $(BUILD_PATH); \
	echo "" > $(COVERAGE_REPORT); \
	for dir in `find . -name "*.go" | grep -o '.*/' | sort -u`; do \
		$(GOTEST) $$dir -coverprofile=$(COVERAGE_PROFILE) -covermode=$(COVERAGE_MODE); \
		if [ $$? != 0 ]; then \
			exit 2; \
		fi; \
		if [ -f $(COVERAGE_PROFILE) ]; then \
			cat $(COVERAGE_PROFILE) >> $(COVERAGE_REPORT); \
			rm $(COVERAGE_PROFILE); \
		fi; \
	done; \

packages:
	for os in $(PKG_OS); do \
		for arch in $(PKG_ARCH); do \
			cd $(BUILD_PATH); \
			mkdir -p $(ARTIFACTS_PATH)/$(PROJECT)_$${os}_$${arch}; \
			for cmd in $(COMMANDS); do \
				if [ -d "$${cmd}" ]; then \
					cd $${cmd}; \
				fi; \
				GOOS=$${os} GOARCH=$${arch} $(GOCMD) build -ldflags \
				"-X main.version=$(BRANCH) -X main.build=$(BUILD) -X main.commit=$(COMMIT)" \
				-o "$(ARTIFACTS_PATH)/$(PROJECT)_$${os}_$${arch}/`basename $${PWD}`" .; \
				cd $(BUILD_DIR); \
			done; \
			cd $(ARTIFACTS_PATH); \
			tar -cvzf $(PROJECT)_$(BRANCH)_$${os}_$${arch}.tar.gz $(PROJECT)_$${os}_$${arch}/; \
		done; \
	done; \

