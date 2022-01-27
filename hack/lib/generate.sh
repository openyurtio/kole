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

#!/usr/bin/env bash

set -x


set -o errexit
set -o nounset
set -o pipefail

CUSTOM_RESOURCE_GROUP="lite"
CUSTOM_RESOURCE_VERSION="v1alpha1"


function generate() {
    (
        cd $KOLE_ROOT
        local crdPath="./config/crd"
        local rbacPath="./config/rbac"
        controller-gen crd:crdVersions=v1 paths="./pkg/apis/..." output:crd:artifacts:config=${crdPath}
        controller-gen rbac:roleName=kole paths="./pkg/..." output:rbac:artifacts:config=${rbacPath}

        go mod tidy

        go mod vendor

        chmod +x ./vendor/k8s.io/code-generator/generate-groups.sh

        ./vendor/k8s.io/code-generator/generate-groups.sh \
        "all"  \
        $KOLE_PACKAGE_NAME/pkg/client \
        $KOLE_PACKAGE_NAME/pkg/apis \
        $CUSTOM_RESOURCE_GROUP:$CUSTOM_RESOURCE_VERSION \
        --go-header-file $KOLE_ROOT/hack/tools/boilerplate.go.txt \
    )
}

