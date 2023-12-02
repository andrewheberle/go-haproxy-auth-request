# go-haproxy-auth-request

[![Go Report Card](https://goreportcard.com/badge/github.com/andrewheberle/go-haproxy-auth-request?logo=go&style=flat-square)](https://goreportcard.com/report/github.com/andrewheberle/go-haproxy-auth-request)

This is service written in Go to integrate [Authelia](https://www.authelia.com/) with HAProxy using [SPOE](https://www.haproxy.org/download/2.8/doc/SPOE.txt) as an alternative to [haproxy-auth-request](https://github.com/TimWolla/haproxy-auth-request/) via Lua.

## Install

```sh
go install github.com/andrewheberle/go-haproxy-auth-request/cmd/haproxy-auth-request
```

### Docker

```sh
docker run -p 3000:3000 ghcr.io/andrewheberle/go-haproxy-auth-request
```

**Note:** The container is built with the environment variable `AUTH_LISTEN=:3000` set. 

## Usage

The service can be managed via the following command line flags:

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

In addition you may set these values via environment variables with the `AUTH_` prefix, such as:

* `AUTH_URL=http://authelia:9091/api/verify`
* `AUTH_METHOD=GET`

## Returned Variables

A number of variables are returned to HAProxy that can be used to determine if
authentication was successful and the identity of the user in question.

These variables are available for the whole HTTP transaction (request and
response):

* `txn.<PREFIX>.response_successful` : Set to `true` if the authentication
  request was successful or `false` otherwise.
* `txn.<PREFIX>.response_code` : Set to the HTTP response code of the
  authentication request.
* `txn.<PREFIX>.response_redirect` : Set to `true` if the response was a
  redirect response code (ie 30X) or if a redirect for authentication is
  required (ie a response code of 401).
* `txn.<PREFIX>.response_location` : Set to value of the `Location` header in
  the response to the authentication request. This is used to redirect the
  user to Authelia for authentication.

The following variables are made available in the request scope only to
minimise memory usage:

* `req.<PREFIX>.response_header.*` : Set to value the relevant authentication
  request header, such as `Remote-User` or `Remote-Email` in order to identify
  the authenticated user.

  The name of the variable will be all in lower case with any dash (`-`)
  replaced with an underscore (`_`).

  So the `Remote-User` reponse header is passed back to HAProxy as
  `req.<PREFIX>.response_header.remote_user`.

In the above `<PREFIX>` is set to the value of the spoe agent name in `spoe.cfg`.

## Handling Authentication

The general concepts to implement in your HAProxy configuration are:

1. Add the headers Authelia expects:
  * X-Forward-For
  * X-Forwarded-Proto
  * X-Forwarded-Host
  * X-Forwarded-Uri
  * X-Forwarded-Method
2. Send the SPOE request
3. Handle the response and redirect for authentication, allow or deny access

The following HAProxy configuration snippet shows this process:

```text
# a protected backend
backend be_protected
        # set required headers
        http-request set-header X-Forward-For %[src]
	http-request set-header X-Forwarded-Proto %[ssl_fc,iif(https,http)]
	http-request set-header X-Forwarded-Host %[req.hdr(Host)]
	http-request set-header X-Forwarded-Uri %[capture.req.uri]
	http-request set-header X-Forwarded-Method %[capture.req.method]

	# set up spoe filter
	filter spoe engine auth-request config /usr/local/etc/haproxy/spoe.cfg

	# send to spoe and act on response
	http-request send-spoe-group auth-request auth-request-group
	http-request redirect location %[var(txn.auth_request.response_location)] if { var(txn.auth_request.response_redirect) -m bool } !{ var(txn.auth_request.response_successful) -m bool }
	http-request deny if !{ var(txn.auth_request.response_successful) -m bool }

        # have your server(s) here
```

## Examples

See the [Examples](examples/README.md) directory for more information.
