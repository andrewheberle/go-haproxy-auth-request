FROM golang:1.21@sha256:672a2286da3ee7a854c3e0a56e0838918d0dbb1c18652992930293312de898a6 AS builder

COPY . /build

RUN cd /build && \
    go build ./cmd/haproxy-auth-request

FROM gcr.io/distroless/base-debian12:nonroot@sha256:4f20cde3246b0192549d6547a0e4cb6dbb84df7e0fa1cfaabbe9be75f532d5c7

COPY --from=builder /build/haproxy-auth-request /app/haproxy-auth-request

ENV AUTH_LISTEN=":3000"

ENTRYPOINT [ "/app/haproxy-auth-request" ]
