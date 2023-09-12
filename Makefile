# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOIMPORTS ?= go run -modfile hack/go.mod golang.org/x/tools/cmd/goimports
DOCKER_ORG ?==registry.harbor.learn.tapsme.org/convention-service

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
	$(GOIMPORTS) --local github.com/x95castle1/probes-convention-service -w .

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: image
image:
	pack build --publish $(DOCKER_ORG)/simple-conventions:latest

.PHONY: install
install: test ## Install conventions server
	kubectl apply -f install-server/server-it.yaml

.PHONY: uninstall
uninstall: ## Uninstall conventions server
	kubectl delete -f install-server/server-it.yaml

.PHONY: restart
restart: ## Kill the convention pods and allow them to be restarted
	kubectl get pods -n simple-convention | grep simple | awk '{print $$1}' | xargs kubectl delete pods -n simple-convention

.PHONY: apply
apply:
	kubectl delete workload -n dev app-golang-kpack --ignore-not-found
	tanzu apps workload create app-golang-kpack \
		--namespace dev \
  		--git-branch main \
  		--git-repo https://github.com/carto-run/app-golang-kpack \
		--param-yaml annotations='{"x95castle1.org/readinessProbe":"{\"httpGet\":{\"path\":\"/healthz\",\"port\":8080},\"initialDelaySeconds\":5,\"periodSeconds\":5}","x95castle1.org/livenessProbe":"{\"exec\":{\"command\":[\"cat\",\"/tmp/healthy\"]},\"initialDelaySeconds\":5,\"periodSeconds\":5}"}' \
  		--label apps.tanzu.vmware.com/has-tests=true \
  		--label app.kubernetes.io/part-of=app-golang-kpack \
  		--param-yaml testing_pipeline_matching_labels='{"apps.tanzu.vmware.com/pipeline":"golang-pipeline"}' \
  		--type web \
  		--yes

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

