# Builds the zenad image used by `make localnet-start` (see contrib/images/Makefile).
ARG IMG_TAG=latest
# Keep in sync with the `go` directive in go.mod.
ARG GO_VERSION=1.23.8

# Compile the zenad binary
FROM golang:${GO_VERSION}-alpine AS zenad-builder
WORKDIR /work
ENV PACKAGES="curl build-base git bash file linux-headers eudev-dev"
RUN apk add --no-cache $PACKAGES

COPY go.mod go.sum* ./
COPY zenad/go.mod zenad/go.sum* ./zenad/
RUN go mod download && (cd zenad && go mod download)

COPY . .
RUN LEDGER_ENABLED=false COSMOS_BUILD_OPTIONS="staticlink" BUILD_TAGS=muslc make build
RUN echo "Checking binary linkage..." \
    && file /work/build/zenad \
    && (file /work/build/zenad | grep -q "statically linked" || echo "Warning: Binary may not be statically linked")

FROM alpine:$IMG_TAG AS run
RUN apk add --no-cache build-base jq bash curl
RUN addgroup -g 1025 nonroot
RUN adduser -D nonroot -u 1025 -G nonroot

# Set up the runtime environment
EXPOSE 26656 26657 1317 9090 26660 8545 8100
STOPSIGNAL SIGTERM
VOLUME /zenad
WORKDIR /zenad

# Copy the wrapper script and binary to expected locations
COPY contrib/images/zenad-env/wrapper.sh /usr/bin/wrapper.sh
COPY --from=zenad-builder /work/build/zenad /zenad/
COPY --from=zenad-builder /work/build/zenad /usr/local/bin/

# Set proper ownership and permissions before switching to nonroot user
RUN chown nonroot:nonroot /usr/bin/wrapper.sh && chmod +x /usr/bin/wrapper.sh
RUN chown -R nonroot:nonroot /zenad

USER nonroot

ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start", "--log_format", "plain", "--minimum-gas-prices", "0.0001aznnt", "--json-rpc.enable", "true", "--json-rpc.api", "eth,txpool,personal,net,debug,web3", "--chain-id", "local-4221"]
