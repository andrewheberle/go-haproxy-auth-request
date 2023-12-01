# go-haproxy-auth-request

[![Go Report Card](https://goreportcard.com/badge/github.com/andrewheberle/go-haproxy-auth-request?logo=go&style=flat-square)](https://goreportcard.com/report/github.com/andrewheberle/go-haproxy-auth-request)

This is service written in Go to integrate [Authelia](https://www.authelia.com/) with HAProxy using [SPOE](https://www.haproxy.org/download/2.8/doc/SPOE.txt) as an alternative to [haproxy-auth-request](https://github.com/TimWolla/haproxy-auth-request/) via Lua.

## Install

```sh
go install github.com/andrewheberle/go-haproxy-auth-request/cmd/haproxy-auth-request
```

## Usage

```sh
./haproxy-auth-request --help
Usage of haproxy-auth-request.exe:
      --debug              Enable debug logging
      --headers strings    HTTP Headers to return on success (default [authorization,proxy-authorization,remote-user,remote-groups,remote-name,remote-email])
      --listen string      Listen address (default "127.0.0.1:3000")
      --method string      HTTP Method for authentication request (default "HEAD")
      --timeout duration   Timeout for verification (default 5s)
      --url string         URL to perform verification against (default "http://127.0.0.1:9091/api/authz/forward-auth")
```

## Example

See the [Examples](examples/README.md) directory for more information.
