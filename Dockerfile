#
# Copyright (c) 2021 Jiangxing Intelligence
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

ARG BASE=golang:1.18-alpine3.16
FROM ${BASE} AS builder

ARG ADD_BUILD_TAGS=""
ARG MAKE="make -e ADD_BUILD_TAGS=$ADD_BUILD_TAGS build"

ARG ALPINE_PKG_BASE="make git openssh-client gcc libc-dev zeromq-dev libsodium-dev"
ARG ALPINE_PKG_EXTRA=""

# set the working directory
WORKDIR /device-gpio

# Replicate the APK repository override.
# If it is no longer necessary to avoid the CDN mirros we should consider dropping this as it is brittle.
RUN sed -e 's/dl-cdn[.]alpinelinux.org/dl-4.alpinelinux.org/g' -i~ /etc/apk/repositories
# Install our build time packages.
RUN apk add --update --no-cache ${ALPINE_PKG_BASE} ${ALPINE_PKG_EXTRA}

COPY go.mod vendor* ./
RUN [ ! -d "vendor" ] && go mod download all || echo "skipping..."

COPY . .
RUN ${MAKE}

# Next image - Copy built Go binary into new workspace
FROM alpine:3.16

LABEL license='SPDX-License-Identifier: Apache-2.0' \
  copyright='Copyright (c) 2021: Jiangxing Intelligence'

RUN sed -e 's/dl-cdn[.]alpinelinux.org/dl-4.alpinelinux.org/g' -i~ /etc/apk/repositories
RUN apk add --update --no-cache zeromq

WORKDIR /
COPY --from=builder /device-gpio/Attribution.txt /
COPY --from=builder /device-gpio/LICENSE /
COPY --from=builder /device-gpio/cmd/ /

EXPOSE 59910

ENTRYPOINT ["/device-gpio"]
CMD ["-cp=consul.http://edgex-core-consul:8500", "--registry", "--confdir=/res"]
