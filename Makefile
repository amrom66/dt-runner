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