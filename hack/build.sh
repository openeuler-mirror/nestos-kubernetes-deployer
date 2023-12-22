#!/bin/sh
# Copyright 2023 KylinSoft  Co., Ltd.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# 	http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

version() { IFS="."; printf "%03d%03d%03d" $@; unset IFS;}

minimum_go_version=1.17
current_go_version=$(go version | cut -d " " -f 3)

if [ "$(version "${current_go_version#go}")" -lt "$(version "$minimum_go_version")" ]; then
     echo "Go version should be greater or equal to the $minimum_go_version"
     exit 1
fi

echo "building..."
sudo go build -mod=vendor -tags release --ldflags="-w -s" -o nkd nkd.go