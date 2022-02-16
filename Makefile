# Copyright 2022 The OpenYurt Authors.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: binary fmt vet release controller-gen

all: binary 

# Build binaries in the host environment
#
# ARGS:
#   GOARCH: go GOARCH.
#   GOOS: go GOOS 
#
# Examples:
#   make binary
#   or
#   GOARCH=amd64 GOOS=linux make binary

binary: 
	bash hack/make-rules/binary.sh


# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Build binaries and push docker images.  
# NOTE: this rule can take time, as we build binaries inside containers
#
# ARGS:
#   IMAGES: It is used to define your private image, default is  openyurt/kole:v1 
#   MQTT5_SERVER: It is used to define the address of mqtt5 server, default is mqtt://8.142.157.229:1883. The default address may not be available, you must specify the address of your own MQTT Server.
#
# Examples:
#   make release
#   or 
#   IMAGES="openyurt/kole:v2" make release
#   or
#   MQTT5_SERVER="mqtt://8.142.157.111:1883" IMAGES="openyurt/kole:v2" make release
#    
release: fmt vet generate
	bash hack/make-rules/release-images.sh
	bash hack/make-rules/manifest.sh

clean: 
	-rm -Rf _output
	-rm -Rf $(GOPATH)/bin/controller-gen-kole

generate: controller-gen
	bash hack/make-rules/generate.sh

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen-kole))
	$(warning controller-gen-kole not found, need to get)	
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	echo "replace sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.7.0" >> go.mod ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.7.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	mv $(GOPATH)/bin/controller-gen $(GOPATH)/bin/controller-gen-kole ;\
	} \
	$(warning install controller-gen-kole successfully)	
else
	$(warning controller-gen-kole has installed)
endif
