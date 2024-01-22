FROM golang:1.21@sha256:5f5d61dcb58900bc57b230431b6367c900f9982b583adcabf9fa93fd0aa5544a AS builder

COPY . /build

RUN cd /build && \
    go build ./cmd/haproxy-auth-request

FROM gcr.io/distroless/base-debian12:nonroot@sha256:684dee415923cb150793530f7997c96b3cef006c868738a2728597773cf27359

COPY --from=builder /build/haproxy-auth-request /app/haproxy-auth-request

ENV AUTH_LISTEN=":3000"

ENTRYPOINT [ "/app/haproxy-auth-request" ]
