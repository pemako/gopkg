.PHONY: check delete push replace tidy help
.DEFAULT_GOAL := help

specTag ?=
delTag ?=
# 获取当前最大的版本号
currentVersion := $(shell git tag --list 'v*' --sort HEAD | sort -r | head -n 1)
# 对当前最大的版本号进行 +1 仅操作第三位版本号
nextVersion := $(shell echo $(currentVersion) | awk -F '.' '{print $$1 "." $$2 "." $$3+1}')

# 使用条件判断来判断 specTag 是否为空，如果不为空则替换 nextVersion 的值为 specTag
$(if $(strip $(specTag)), \
	$(eval nextVersion := $(specTag)), \
)

check: ## Run govulncheck
	./scripts/build.sh --mode check

delete: ## Run delete tag
	./scripts/build.sh --mode delete --version $(delTag)

push: ## Run push tag
	@while [ -z "$$commitMsg" ]; do \
		read -p "Enter commit message > " commitMsg; \
		if [ -z "$$commitMsg" ]; then \
			echo "Commit message cannot by empty. Please try again."; \
		fi; \
	done; \
	echo "You entered msg is $$commitMsg"; \
	./scripts/build.sh --mode push --version $(nextVersion) --commitMsg "$$commitMsg"

replace: ## Run replace tag
	@echo "spec tag is" $(specTag)
	@echo "will next version is $(nextVersion) replace current version $(currentVersion)"

	./scripts/build.sh --mode replace --version $(currentVersion) --nextVersion $(nextVersion)

tidy: ## Run mod tidy
	./scripts/build.sh --mode tidy

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
