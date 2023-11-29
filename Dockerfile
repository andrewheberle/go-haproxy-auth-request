FROM golang:1.21@sha256:57bf74a970b68b10fe005f17f550554406d9b696d10b29f1a4bdc8cae37fd063 AS builder

COPY . /build

RUN cd /build && \
    go build ./cmd/haproxy-auth-request

FROM gcr.io/distroless/base-debian12:nonroot@sha256:5a779e9c2635dbea68ae7988f398f95686ccde186cd2abf51207e41ed2ec51f4

COPY --from=builder /build/haproxy-auth-request /app/haproxy-auth-request

ENV AUTH_LISTEN=":3000"

ENTRYPOINT [ "/app/haproxy-auth-request" ]
