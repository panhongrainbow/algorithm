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

# Default target to run tests in all directories.
.PHONY: all
all: test

# Test target to run tests in all specified directories.
.PHONY: test
test:
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

# Clean up any generated files or artifacts.
.PHONY: clean
clean:
	@for dir in $(DIRS); do \
		echo "$(YELLOW)$(WASTEBASKET) Cleaning directory: $$dir$(RESET)"; \
		cd $$dir && go clean; \
		cd ..; \
	done
