#!/bin/bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
YAML_FILE="$CURRENT_DIR/manifest.yaml"

CONFIG_BASE_DIR="$CURRENT_DIR/../base"
CONFIG_CRD_DIR="$CURRENT_DIR/../crd"
CONFIG_RBAC_DIR="$CURRENT_DIR/../rbac"

cat << EOF >$YAML_FILE
apiVersion: v1
kind: Namespace
metadata:
  name: kole
---
EOF

function create_yaml_file() {
    dir=$1
    for file in $dir/*
    do
        file_name=$(basename $file)
        if ! test -f $file ; then
            #echo "$file_name is not a file"
            continue
        fi
        extension="${file##*.}"
        if [ "$extension" != "yaml" ]; then
            #echo "$file_name is not a yaml file"
            continue
        fi
        cat "$file" >> $YAML_FILE
        echo "---" >> $YAML_FILE
    done
}

create_yaml_file $CONFIG_BASE_DIR
create_yaml_file $CONFIG_CRD_DIR
create_yaml_file $CONFIG_RBAC_DIR
