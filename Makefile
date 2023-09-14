# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOIMPORTS ?= go run -modfile hack/go.mod golang.org/x/tools/cmd/goimports
DOCKER_ORG ?= registry.harbor.learn.tapsme.org/convention-service
LATEST_TAG := $(shell git describe --tags --abbrev=0)
DEV_IMAGE_LOCATION ?= harbor-repo.vmware.com/tanzu_practice/conventions/multi-purpose-convention-server-bundle-repo
PROMOTION_IMAGE_LOCATION ?= projects.registry.vmware.com/tanzu_practice/conventions/multi-purpose-convention-server-bundle-repo

.PHONY: all
all: test

.PHONY: build
build: test ## Build the project
	go build ./...

.PHONY: test
test: fmt vet ## Run tests
	go test ./... -coverprofile cover.out

.PHONY: fmt
fmt: ## Run go fmt against code
	$(GOIMPORTS) --local github.com/x95castle1/multi-purpose-convention-server -w .

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: image
image:
	pack build --publish $(DOCKER_ORG)/multi-purpose-convention:latest

.PHONY: install
install: test ## Install conventions server
	kubectl apply -f install-server/server-it.yaml

.PHONY: uninstall
uninstall: ## Uninstall conventions server
	kubectl delete -f install-server/server-it.yaml

.PHONY: restart
restart: ## Kill the convention pods and allow them to be restarted
	kubectl get pods -n multi-purpose-convention | grep webhook | awk '{print $$1}' | xargs kubectl delete pods -n multi-purpose-convention

.PHONY: applyw
applyw:
	kubectl apply -f ./examples/workload/.

.PHONY: unapplyw
unapplyw:
	kubectl delete -f ./examples/workload/.

.PHONY: applyp
applyp:
	tanzu package repository add multi-purpose-conventions-repository \
	--url projects.registry.vmware.com/tanzu_practice/conventions/multi-purpose-convention-server-bundle-repo:$(LATEST_TAG) \
	--namespace tap-install \
	--yes
	tanzu package install multi-purpose-convention-server  \
  	--package multi-purpose-convention-server.conventions.tanzu.vmware.com \
  	--values-file ./examples/package/values.yaml \
  	--version $(LATEST_TAG) \
  	--namespace tap-install \
		--yes

.PHONY: unapplyp
unapplyp:
	tanzu package installed delete multi-purpose-convention-server -n tap-install --yes
	tanzu package repository delete multi-purpose-conventions-repository -n tap-install --yes

.PHONY: package
package:
	kctrl package release --chdir ./carvel -v $(LATEST_TAG) --tag $(LATEST_TAG) --repo-output ./packagerepository -y
	kctrl package repo release --chdir carvel/packagerepository -v $(LATEST_TAG) -y

.PHONY: promote
promote:
	imgpkg --tty copy -b $(DEV_IMAGE_LOCATION):$(LATEST_TAG) --to-repo $(PROMOTION_IMAGE_LOCATION) --registry-response-header-timeout 1m --registry-retry-count 2

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

