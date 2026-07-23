FROM --platform=$BUILDPLATFORM golang:1.25.12-alpine3.23 AS builder

WORKDIR /build

COPY go.mod .
COPY main.go .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o /app .

FROM scratch

COPY --from=builder /app /app

ENTRYPOINT ["/app"]
