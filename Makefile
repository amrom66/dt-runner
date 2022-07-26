SHELL := /bin/bash

define GEN_HELP_INFO
# 生成
#
# Example:
# make gen
endef

.PHONY: gen
ifeq ($(PRINT_HELP),y)
gen:
	@echo "$$GEN_HELP_INFO"
else
gen:
	@echo "will generate files"
endif

define BUILD_HELP_INFO
# 编译
#
# Example:
# make build
endef

.PHONY: build
ifeq ($(PRINT_HELP),y)
build:
	@echo "$$BUILD_HELP_INFO"
else
build:
	@echo "will build files"
	go build -o dt-runner
endif

define BUILDAMD_HELP_INFO
# 编译amd64
#
# Example:
# make buildamd
endef

.PHONY: buildamd
ifeq ($(PRINT_HELP),y)
buildamd:
	@echo "$$BUILDAMD_HELP_INFO"
else
buildamd:
	@echo "will build files"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dt-runner-amd64
endif