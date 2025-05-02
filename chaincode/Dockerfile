# SPDX-License-Identifier: Apache-2.0
ARG GO_VER=1.23.0
ARG ALPINE_VER=3.19

FROM golang:${GO_VER}-alpine${ALPINE_VER} AS build

RUN apk add --no-cache \
	bash \
	binutils-gold \
  dumb-init \
	gcc \
	git \
	make \
	musl-dev

ADD . $GOPATH/digital-identity/chaincode
WORKDIR $GOPATH/digital-identity/chaincode

RUN go install ./...

FROM golang:${GO_VER}-alpine${ALPINE_VER}

LABEL org.opencontainers.image.title="Go SmartContract"
LABEL org.opencontainers.image.description="Hyperledger Fabric Go Smart contract for Kubernetes chaincode builder"
LABEL org.opencontainers.image.source="https://github.com/hyperledger-labs/fabric-builder-k8s/tree/main/samples/go-contract"

COPY --from=build /usr/bin/dumb-init /usr/bin/dumb-init
COPY --from=build /go/bin/chaincode /usr/bin/chaincode

WORKDIR /home/chaincode
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "exec /usr/bin/chaincode -peer.address=$CORE_PEER_ADDRESS"]