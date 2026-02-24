FROM golang:1.23.6-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod ./
COPY *.go ./
COPY templates/ templates/
COPY static/ static/
RUN FINAL_HASH=$(git clone --depth=1 https://github.com/uSwapExchange/uswap-zero.git /tmp/repo 2>/dev/null && cd /tmp/repo && git rev-parse HEAD || echo "unknown") && \
    rm -rf /tmp/repo && \
    FINAL_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") && \
    CGO_ENABLED=0 GOOS=linux go build \
      -ldflags "-s -w -X main.commitHash=${FINAL_HASH} -X main.buildTime=${FINAL_TIME}" \
      -o /uswap-zero

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /uswap-zero /uswap-zero
EXPOSE 3000
ENTRYPOINT ["/uswap-zero"]
