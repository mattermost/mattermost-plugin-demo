# ====================================================================================
# Semantic Versioning
# ====================================================================================

# Used for semver bumping
PROTECTED_BRANCH := master
APP_NAME    := $(shell pluginctl manifest get '{{.Id}}')
CURRENT_VERSION := $(shell git describe --abbrev=0 --tags)
VERSION_PARTS := $(subst ., ,$(subst v,,$(subst -rc, ,$(CURRENT_VERSION))))
MAJOR := $(word 1,$(VERSION_PARTS))
MINOR := $(word 2,$(VERSION_PARTS))
PATCH := $(word 3,$(VERSION_PARTS))
RC := $(shell echo $(CURRENT_VERSION) | grep -oE 'rc[0-9]+' | sed 's/rc//')

# Check if current branch is protected
define check_protected_branch
	@current_branch=$$(git rev-parse --abbrev-ref HEAD); \
	if ! echo "$(PROTECTED_BRANCH)" | grep -wq "$$current_branch" && ! echo "$$current_branch" | grep -q "^release"; then \
		echo "Error: Tagging is only allowed from $(PROTECTED_BRANCH) or release branches. You are on $$current_branch branch."; \
		exit 1; \
	fi
endef

# Check if there are pending pulls
define check_pending_pulls
	@git fetch; \
	current_branch=$$(git rev-parse --abbrev-ref HEAD); \
	if [ "$$(git rev-parse HEAD)" != "$$(git rev-parse origin/$$current_branch)" ]; then \
		echo "Error: Your branch is not up to date with upstream. Please pull the latest changes before performing a release"; \
		exit 1; \
	fi
endef

# Prompt for approval
define prompt_approval
	@read -p "About to bump $(APP_NAME) to version $(1), approve? (y/n) " userinput; \
	if [ "$$userinput" != "y" ]; then \
		echo "Bump aborted."; \
		exit 1; \
	fi
endef

.PHONY: patch minor major patch-rc minor-rc major-rc

patch: ## to bump patch version (semver)
	$(call check_protected_branch)
	$(call check_pending_pulls)
	@$(eval PATCH := $(shell echo $$(($(PATCH)+1))))
	$(call prompt_approval,$(MAJOR).$(MINOR).$(PATCH))
	@echo Bumping $(APP_NAME) to Patch version $(MAJOR).$(MINOR).$(PATCH)
	git tag -s -a v$(MAJOR).$(MINOR).$(PATCH) -m "Bumping $(APP_NAME) to Patch version $(MAJOR).$(MINOR).$(PATCH)"
	git push origin v$(MAJOR).$(MINOR).$(PATCH)
	@echo Bumped $(APP_NAME) to Patch version $(MAJOR).$(MINOR).$(PATCH)

minor: ## to bump minor version (semver)
	$(call check_protected_branch)
	$(call check_pending_pulls)
	@$(eval MINOR := $(shell echo $$(($(MINOR)+1))))
	@$(eval PATCH := 0)
	$(call prompt_approval,$(MAJOR).$(MINOR).$(PATCH))
	@echo Bumping $(APP_NAME) to Minor version $(MAJOR).$(MINOR).$(PATCH)
	git tag -s -a v$(MAJOR).$(MINOR).$(PATCH) -m "Bumping $(APP_NAME) to Minor version $(MAJOR).$(MINOR).$(PATCH)"
	git push origin v$(MAJOR).$(MINOR).$(PATCH)
	@echo Bumped $(APP_NAME) to Minor version $(MAJOR).$(MINOR).$(PATCH)

major: ## to bump major version (semver)
	$(call check_protected_branch)
	$(call check_pending_pulls)
	$(eval MAJOR := $(shell echo $$(($(MAJOR)+1))))
	$(eval MINOR := 0)
	$(eval PATCH := 0)
	$(call prompt_approval,$(MAJOR).$(MINOR).$(PATCH))
	@echo Bumping $(APP_NAME) to Major version $(MAJOR).$(MINOR).$(PATCH)
	git tag -s -a v$(MAJOR).$(MINOR).$(PATCH) -m "Bumping $(APP_NAME) to Major version $(MAJOR).$(MINOR).$(PATCH)"
	git push origin v$(MAJOR).$(MINOR).$(PATCH)
	@echo Bumped $(APP_NAME) to Major version $(MAJOR).$(MINOR).$(PATCH)

patch-rc: ## to bump patch release candidate version (semver)
	$(call check_protected_branch)
	$(call check_pending_pulls)
	@$(eval RC := $(shell echo $$(($(RC)+1))))
	$(call prompt_approval,$(MAJOR).$(MINOR).$(PATCH)-rc$(RC))
	@echo Bumping $(APP_NAME) to Patch RC version $(MAJOR).$(MINOR).$(PATCH)-rc$(RC)
	git tag -s -a v$(MAJOR).$(MINOR).$(PATCH)-rc$(RC) -m "Bumping $(APP_NAME) to Patch RC version $(MAJOR).$(MINOR).$(PATCH)-rc$(RC)"
	git push origin v$(MAJOR).$(MINOR).$(PATCH)-rc$(RC)
	@echo Bumped $(APP_NAME) to Patch RC version $(MAJOR).$(MINOR).$(PATCH)-rc$(RC)

minor-rc: ## to bump minor release candidate version (semver)
	$(call check_protected_branch)
	$(call check_pending_pulls)
	@$(eval MINOR := $(shell echo $$(($(MINOR)+1))))
	@$(eval PATCH := 0)
	@$(eval RC := 1)
	$(call prompt_approval,$(MAJOR).$(MINOR).$(PATCH)-rc$(RC))
	@echo Bumping $(APP_NAME) to Minor RC version $(MAJOR).$(MINOR).$(PATCH)-rc$(RC)
	git tag -s -a v$(MAJOR).$(MINOR).$(PATCH)-rc$(RC) -m "Bumping $(APP_NAME) to Minor RC version $(MAJOR).$(MINOR).$(PATCH)-rc$(RC)"
	git push origin v$(MAJOR).$(MINOR).$(PATCH)-rc$(RC)
	@echo Bumped $(APP_NAME) to Minor RC version $(MAJOR).$(MINOR).$(PATCH)-rc$(RC)

major-rc: ## to bump major release candidate version (semver)
	$(call check_protected_branch)
	$(call check_pending_pulls)
	@$(eval MAJOR := $(shell echo $$(($(MAJOR)+1))))
	@$(eval MINOR := 0)
	@$(eval PATCH := 0)
	@$(eval RC := 1)
	$(call prompt_approval,$(MAJOR).$(MINOR).$(PATCH)-rc$(RC))
	@echo Bumping $(APP_NAME) to Major RC version $(MAJOR).$(MINOR).$(PATCH)-rc$(RC)
	git tag -s -a v$(MAJOR).$(MINOR).$(PATCH)-rc$(RC) -m "Bumping $(APP_NAME) to Major RC version $(MAJOR).$(MINOR).$(PATCH)-rc$(RC)"
	git push origin v$(MAJOR).$(MINOR).$(PATCH)-rc$(RC)
	@echo Bumped $(APP_NAME) to Major RC version $(MAJOR).$(MINOR).$(PATCH)-rc$(RC)
