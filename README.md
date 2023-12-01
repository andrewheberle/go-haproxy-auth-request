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

## Returned Variables

A number of variables are returned to HAProxy that can be used to determine if
authentication was successful and the identity of the user in question.

These variables are:

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
* `req.<PREFIX>.response_header.*` : Set to value the relevant authentication
  request header, such as `Remote-User` or `Remote-Email` in order to identify
  the authenticated user.

  The name of the variable will be all in lower case with any dash (`-`)
  replaced with an underscore (`_`).

  So the `Remote-User` reponse header is passed back to HAProxy as
  `req.<PREFIX>.response_header.remote_user`.

In the above `<PREFIX>` is set to the value of the spoe agent name in `spoe.cfg`.

## Example

See the [Examples](examples/README.md) directory for more information.
