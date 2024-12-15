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

ARG GO_VERSION=1.23.0
ARG RUNNER_IMAGE=alpine:3.20
ARG BUILD_TAGS="netgo,muslc,blst,bls12381,pebbledb"
ARG NAME=beacond
ARG APP_NAME=beacond
ARG DB_BACKEND=pebbledb
ARG CMD_PATH=./cmd/beacond

#######################################################
###         Stage 1 - Cache Go Modules              ###
#######################################################

FROM golang:${GO_VERSION}-alpine3.20 AS mod-cache

WORKDIR /workdir

RUN apk add --no-cache git

# Download Go modules
COPY ./go.mod ./go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

#######################################################
###         Stage 2 - Build the Application         ###
#######################################################

FROM golang:${GO_VERSION}-alpine3.20 AS builder

ARG GIT_VERSION
ARG GIT_COMMIT
ARG BUILD_TAGS

# Set the working directory
WORKDIR /workdir

# Consolidate RUN commands to reduce layers
RUN apk add --no-cache --update \
    ca-certificates \
    build-base

# Copy the dependencies from the cache stage
COPY --from=mod-cache /go/pkg /go/pkg

# Copy all the source code (this will ignore files/dirs in .dockerignore)
COPY ./ ./

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

# Build args
ARG APP_NAME

# Copy over built executable into a fresh container
COPY --from=builder /workdir/build/bin/${APP_NAME} /usr/bin/${APP_NAME}

# TODO: We should un hood this part, its very specific
# to our kurtosis setup.
RUN mkdir -p /root/jwt /root/kzg && \
    apk add --no-cache bash sed curl

EXPOSE 26656
EXPOSE 26657

ENTRYPOINT [ "beacond" ]