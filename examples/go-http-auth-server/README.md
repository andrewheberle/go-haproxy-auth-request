# Authelia

## Overview

This is an example deployment that can be run using Docker Compose to integrate
with [go-http-auth-server](https://github.com/andrewheberle/go-http-auth-server).

The pre-requisites for this are:

* set hostnames and any required hosts/DNS entries as required
* provide a SSL certicate as "app.pem" as all traffic must be secured from the
  browser to HAProxy for authentication to work correctly (a self signed
  certificate is sufficient for testing).
* generate a self-signed (RSA only) certificate for the SAML Service Provider
  component as `samlsp.crt` and `samlsp.key` for certificate and key
  respectively.

Once this is in place, the containers may be started using Docker Compose:

```sh
docker compose up
```
