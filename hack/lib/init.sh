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

set -o errexit
set -o nounset
set -o pipefail

KOLE_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
KOLE_PACKAGE_NAME=${KOLE_ROOT##*src/}
KOLE_MOD="$(head -1 $KOLE_ROOT/go.mod | awk '{print $2}')"
KOLE_OUTPUT_DIR=${KOLE_ROOT}/_output
KOLE_BIN_DIR=${KOLE_OUTPUT_DIR}/bin

GIT_COMMIT=$(git rev-parse --short HEAD)
GIT_COMMIT_SHORT=$GIT_COMMIT
GIT_VERSION=${GIT_VERSION:-v0.2.0}
BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
REPO=${REPO:-openyurt}
TAG=${TAG:-${GIT_COMMIT_SHORT}}

source "${KOLE_ROOT}/hack/lib/binary.sh"
source "${KOLE_ROOT}/hack/lib/release-images.sh"
source "${KOLE_ROOT}/hack/lib/generate.sh"
