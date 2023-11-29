# go-haproxy-auth-request

[![Go Report Card](https://goreportcard.com/badge/github.com/andrewheberle/go-http-auth-request?logo=go&style=flat-square)](https://goreportcard.com/report/github.com/andrewheberle/go-http-auth-request)

This is service written in Go to integrate [Authelia](https://www.authelia.com/) with HAProxy using [SPOE](https://www.haproxy.org/download/2.8/doc/SPOE.txt) as an alternative to [haproxy-auth-request](https://github.com/TimWolla/haproxy-auth-request/) via Lua.

## Install

```sh
go install github.com/andrewheberle/go-haproxy-auth-request/cmd/haproxy-auth-request
```

## Usage

```sh
./haproxy-auth-request --help
Usage of haproxy-auth-request.exe:
      --cookie string      Session cookie name (default "authelia_session")
      --debug              Enable debug logging
      --listen string      Listen address (default "127.0.0.1:3000")
      --timeout duration   Timeout for verification (default 5s)
      --url string         URL to perform verification against (default "http://127.0.0.1:9091/api/verify")
```

## Example

See the [Examples](examples/README.md) directory for more information.
