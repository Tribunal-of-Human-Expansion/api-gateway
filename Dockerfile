# --- build ---
FROM golang:1.21-alpine AS build
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/gateway .

# --- runtime ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget \
    && adduser -D -u 65532 gateway
COPY --from=build /out/gateway /usr/local/bin/gateway
USER 65532:65532
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/health || exit 1
ENTRYPOINT ["/usr/local/bin/gateway"]
