FROM golang:1.22@sha256:ef61a20960397f4d44b0e729298bf02327ca94f1519239ddc6d91689615b1367 AS builder

COPY . /build

RUN cd /build && \
    go build ./cmd/haproxy-auth-request

FROM gcr.io/distroless/base-debian12:nonroot@sha256:3d46a42c741d469811b95d0b6ec9d10c15ebb1be7c5eaa989d429d91b066d78c

COPY --from=builder /build/haproxy-auth-request /app/haproxy-auth-request

ENV AUTH_LISTEN=":3000"

ENTRYPOINT [ "/app/haproxy-auth-request" ]
