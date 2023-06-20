# Copyright 2023 KylinSoft  Co., Ltd.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

IMG_OPERATOR ?= housekeeper-operator:latest
IMG_CONTROLLER ?= housekeeper-controller:latest

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
CONTROLLER_TOOLS_VERSION ?= v0.9.2

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec


##@ Build
.PHONY: all
all: housekeeper-operator-manager housekeeper-controller-manager
# Build binary
housekeeper-operator-manager:
	go build -o bin/housekeeper-operator-manager operator/housekeeper-operator/main.go
housekeeper-controller-manager:
	go build -o bin/housekeeper-controller-manager operator/housekeeper-controller/main.go

# Build the docker image
.PHONY: docker-build
docker-build:  ## Build docker image with the housekeeper-operator-manager.
	docker build -t ${IMG_OPERATOR} .
	docker build -t ${IMG_CONTROLLER} .

.PHONY: docker-push
docker-push: ## Push docker image with the housekeeper-operator-manager.
	docker push ${IMG_OPERATOR}
	docker push ${IMG_CONTROLLER}

# ##@ Development
.PHONY: manifests
manifests: controller-gen ##Generate manifests e.g. CRD/RBAC
	$(CONTROLLER_GEN) rbac:roleName=update-manager-role crd paths="./..." output:crd:artifacts:config=config/crd

.PHONY: install
install: 	## Install CRD in a cluster
	kubectl apply -f config/crd

.PHONY: uninstall
uninstall:	## Uninstall CRD from a cluster
	kubectl delete -f config/crd

.PHONY: deploy
deploy: ## Deploy controller
	kubectl apply -f config/rbac
	kubectl apply -f config/manager

.PHONY: undeploy
undeploy: ## Undeploy controller
	kubectl delete -f config/rbac
	kubectl delete -f config/manager

## Location to install dependencies to
.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)
