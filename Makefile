# Define color codes
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
RESET := \033[0m

# Define directories
DIRS := testToolset

# Define emoji
CHECK_MARK := ✅
CROSS_MARK := ❌
WASTEBASKET := 🗑️
ALCHEMY := ⚗️
LAB := 🧪

# go clean -cache removes the Go build cache and deletes $HOME/.cache/go-build on Linux.
.PHONY: clean
clean:
	@go clean --cache

# Test target to run tests in all specified directories.
.PHONY: test-all
test-all:
	@for dir in $(DIRS); do \
		echo "$(YELLOW)$(ALCHEMY)  Testing directory: $$dir$(RESET)"; \
		cd $$dir && go test -v ./...; \
		if [ $$? -eq 0 ]; then \
			echo "$(GREEN)$(CHECK_MARK) Tests passed in $$dir$(RESET)"; \
		else \
			echo "$(RED)$(CROSS_MARK) Tests failed in $$dir$(RESET)"; \
		fi; \
		cd ..; \
	done

# Default target: clean Go build cache and run tests in all directories.
.PHONY: all
all: clean test-all

# Run tests in the ./utilhub directory with verbose output and no timeout.
.PHONY: test-utilhub
test-utilhub:
	@go test -v ./utilhub -timeout=0 -run Test_AnsiColorOutput
