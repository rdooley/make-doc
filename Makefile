SHELL=/bin/bash

BUILD_DIR=build
PROJECT_NAME=make-doc ## @build Project name
DEFAULT_TARGET=$(BUILD_DIR)/$(PROJECT_NAME)

.PHONY: clean build install

clean: ## @build Clean stuff
		rm -rf $(BUILD_DIR)

build: ## @build Build the actual thing
		mkdir -p $(BUILD_DIR)
		# default build
		go build -o $(DEFAULT_TARGET)
		# build for specific operating systems
		BUILD_DIR=$(BUILD_DIR) PROJECT_NAME=$(PROJECT_NAME) ./build.sh "linux/amd64" "darwin/amd64"

install: build
	cp $(DEFAULT_TARGET) /usr/local/bin
