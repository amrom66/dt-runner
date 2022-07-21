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

