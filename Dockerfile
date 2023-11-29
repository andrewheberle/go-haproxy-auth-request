FROM golang:1.21@sha256:57bf74a970b68b10fe005f17f550554406d9b696d10b29f1a4bdc8cae37fd063 AS builder

COPY . /build

RUN cd /build && \
    go build ./cmd/haproxy-auth-request

FROM gcr.io/distroless/base-debian12@sha256:1dfdb5ed7d9a66dcfc90135b25a46c25a85cf719b619b40c249a2445b9d055f5

COPY --from=builder /build/haproxy-auth-request /app/haproxy-auth-request

ENV AUTH_LISTEN=":3000"

ENTRYPOINT [ "/app/haproxy-auth-request" ]
