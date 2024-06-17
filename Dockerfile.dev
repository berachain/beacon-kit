# syntax=docker/dockerfile:1
#
# Copyright (C) 2022, Berachain Foundation. All rights reserved.
# See the file LICENSE for licensing terms.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
# AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
# IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
# DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
# FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
# DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
# SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
# CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
# OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
# OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

#######################################################
###           Stage 0 - Build Arguments             ###
#######################################################

ARG GO_VERSION=1.22.4
ARG RUNNER_IMAGE=alpine:latest
ARG BUILD_TAGS="netgo,muslc,blst,bls12381,pebbledb"
ARG NAME=beacond
ARG APP_NAME=beacond
ARG DB_BACKEND=pebbledb
ARG CMD_PATH=./beacond/cmd

#######################################################
###         Stage 1 - Cache Go Modules              ###
#######################################################

FROM golang:${GO_VERSION}-alpine3.19 as mod-cache

WORKDIR /workdir

RUN apk add git

COPY ./beacond/go.mod ./beacond/go.sum ./beacond/
COPY ./mod/async/go.mod ./mod/async/go.sum ./mod/async/
COPY ./mod/beacon/go.mod ./mod/beacon/go.sum ./mod/beacon/
COPY ./mod/cli/go.mod ./mod/cli/go.sum ./mod/cli/
COPY ./mod/consensus-types/go.mod ./mod/consensus-types/go.sum ./mod/consensus-types/
COPY ./mod/config/go.mod ./mod/config/go.sum ./mod/config/
COPY ./mod/da/go.mod ./mod/da/go.sum ./mod/da/
COPY ./mod/engine-primitives/go.mod ./mod/engine-primitives/go.sum ./mod/engine-primitives/
COPY ./mod/execution/go.mod ./mod/execution/go.sum ./mod/execution/
COPY ./mod/interfaces/go.mod ./mod/interfaces/
COPY ./mod/log/go.mod ./mod/log/
COPY ./mod/node-api/go.mod ./mod/node-api/go.sum ./mod/node-api/
COPY ./mod/node-core/go.mod ./mod/node-core/go.sum ./mod/node-core/
COPY ./mod/p2p/go.mod ./mod/p2p/
COPY ./mod/payload/go.mod ./mod/payload/go.sum ./mod/payload/
COPY ./mod/primitives/go.mod ./mod/primitives/go.sum ./mod/primitives/
COPY ./mod/runtime/go.mod ./mod/runtime/go.sum ./mod/runtime/
COPY ./mod/state-transition/go.mod ./mod/state-transition/go.sum ./mod/state-transition/
COPY ./mod/storage/go.mod ./mod/storage/go.sum ./mod/storage/
COPY ./mod/errors/go.mod ./mod/errors/go.sum ./mod/errors/

RUN go work init
RUN go work use ./beacond
RUN go work use ./mod/async
RUN go work use ./mod/beacon
RUN go work use ./mod/cli
RUN go work use ./mod/config
RUN go work use ./mod/consensus-types
RUN go work use ./mod/da
RUN go work use ./mod/execution
RUN go work use ./mod/interfaces
RUN go work use ./mod/log
RUN go work use ./mod/node-api
RUN go work use ./mod/node-core
RUN go work use ./mod/p2p
RUN go work use ./mod/payload
RUN go work use ./mod/primitives
RUN go work use ./mod/engine-primitives
RUN go work use ./mod/runtime
RUN go work use ./mod/state-transition
RUN go work use ./mod/storage
RUN go work use ./mod/errors

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

#######################################################
###         Stage 2 - Build the Application         ###
#######################################################

FROM golang:${GO_VERSION}-alpine3.19 as builder

ARG GIT_VERSION
ARG GIT_COMMIT
ARG BUILD_TAGS

# Set the working directory
WORKDIR /workdir

# Consolidate RUN commands to reduce layers
RUN apk add ca-certificates build-base linux-headers git && \
    set -eux

# Copy the dependencies from the cache stage as well as the
# go.work file to the working directory
COPY --from=mod-cache /go/pkg /go/pkg
COPY --from=mod-cache ./workdir/go.work ./go.work

# Copy the rest of the source code
COPY ./mod ./mod
COPY ./beacond ./beacond

# Build args
ARG NAME
ARG APP_NAME
ARG DB_BACKEND
ARG CMD_PATH

# Build beacond
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    env NAME=${NAME} DB_BACKEND=${DB_BACKEND} APP_NAME=${APP_NAME} CGO_ENABLED=1 && \
    go build \
    -mod=readonly \
    -tags ${BUILD_TAGS} \
    -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=${NAME} \
    -X github.com/cosmos/cosmos-sdk/version.AppName=${APP_NAME} \
    -X github.com/cosmos/cosmos-sdk/version.Version=${GIT_VERSION} \
    -X github.com/cosmos/cosmos-sdk/version.Commit=${GIT_COMMIT} \
    -X github.com/cosmos/cosmos-sdk/version.BuildTags=${BUILD_TAGS} \
    -X github.com/cosmos/cosmos-sdk/types.DBBackend=$DB_BACKEND \
    -w -s -linkmode=external -extldflags '-Wl,-z,muldefs -static'" \
    -trimpath \
    -o /workdir/build/bin/beacond \
    ${CMD_PATH}

#######################################################
###        Stage 3 - Prepare the Final Image        ###
#######################################################

FROM ${RUNNER_IMAGE}

# Copy over built executable into a fresh container.
COPY --from=builder /workdir/build/bin/beacond /usr/bin
RUN mkdir -p /root/jwt /root/kzg && \
    apk add --no-cache bash sed curl

ENTRYPOINT [ "/usr/bin/beacond" ]
