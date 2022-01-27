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

# project_info generates the project information and the corresponding valuse 
# for 'ldflags -X' option
project_info() {
    PROJECT_INFO_PKG=${KOLE_MOD}/pkg/projectinfo
    echo "-X ${PROJECT_INFO_PKG}.gitVersion=${GIT_VERSION}"
    echo "-X ${PROJECT_INFO_PKG}.gitCommit=${GIT_COMMIT}"
    echo "-X ${PROJECT_INFO_PKG}.buildDate=${BUILD_DATE}"
}


build_binary() {
    local goflags goldflags gcflags
    goldflags="${GOLDFLAGS:--s -w $(project_info)}"
    gcflags="${GOGCFLAGS:-}"
    goflags=${GOFLAGS:-}

    local bin_dir=${KOLE_BIN_DIR}
    rm -rf ${bin_dir}/*
    mkdir -p ${bin_dir}
    allCommand=$(ls $KOLE_ROOT/cmd/)
    for subCommand in $allCommand; do
        echo $subCommand
        echo ""
        (
            local bin_name=$subCommand
    
            echo "Building ${bin_name}"
            cd  $KOLE_ROOT/cmd/$subCommand
            go build -o ${bin_dir}/${bin_name} \
            -ldflags "${goldflags:-}" \
            -gcflags "${gcflags:-}" ${goflags} 
        )
    done
}


