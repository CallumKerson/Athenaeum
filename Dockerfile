FROM golang:1.20-alpine as base
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | /bin/ash -s v1.50.1

FROM base as dependencies
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

FROM dependencies as src
COPY cmd/ ./cmd
COPY internal/ ./internal
COPY pkg/ ./pkg

FROM src as test
RUN CGO_ENABLED=0 GOOS=linux go test -v ./...

FROM test as builder
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /bin/ ./...

FROM gcr.io/distroless/static-debian11:nonroot AS production

WORKDIR /app
COPY --from=builder /bin/athenaeum /usr/bin/athenaeum
COPY --from=busybox:1.35.0-uclibc /bin/wget /usr/bin/wget

EXPOSE 8080

HEALTHCHECK --interval=3s \
    --timeout=2s \
    --start-period=5s \
    CMD ["/usr/bin/wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]

ENTRYPOINT ["/usr/bin/athenaeum"]
