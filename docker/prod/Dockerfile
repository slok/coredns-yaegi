# Set also `ARCH` ARG here so we can use it on all the `FROM`s. 
ARG ARCH

FROM golang:1.22-alpine as build-stage

ARG COREDNS_VERSION="1.11.1"
ARG VERSION

RUN apk --no-cache add \
    g++ \
    git \
    make \
    curl \
    bash

# Compile, we are settings the plugin just before forward plugin.
WORKDIR /tmp
RUN git clone --depth 1 --branch v${COREDNS_VERSION} https://github.com/coredns/coredns && \
    cd coredns && \
    go get github.com/slok/coredns-yaegi/@${VERSION} && \
    sed -i '/loop:loop/a yaegi:github.com/slok/coredns-yaegi' ./plugin.cfg && \
    make


# Although we are on an specific architecture (normally linux/amd64) our go binary has been built for
# ${ARCH} specific architecture.
# To make portable our building process we base our final image on that same architecture as the binary 
# to obtain a resulting ${ARCH} image independently where we are building this image.
FROM gcr.io/distroless/static:nonroot-${ARCH}

COPY --from=build-stage /tmp/coredns/coredns /usr/local/bin/coredns

ENTRYPOINT ["/usr/local/bin/coredns"]