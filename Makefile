SHELL := /bin/bash

SCRIPT_PATH := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

TARGETS_DIR = targets
TARGETS = $(SCRIPT_PATH)$(TARGETS_DIR)

socksproxy: $(TARGETS)/socksproxy
$(TARGETS)/socksproxy: main.go
	@echo "Build socksproxy..."
	@go build -o $(TARGETS)/socksproxy main.go

run-socksproxy: socksproxy
	$(TARGETS)/socksproxy > $(TARGETS)/socksproxy.log

run-socksproxy-quiet: socksproxy
	$(TARGETS)/socksproxy &> $(TARGETS)/socksproxy.log &

clean:
	rm -rf $(TARGETS)