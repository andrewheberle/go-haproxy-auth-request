# Authelia

## Overview

This is an example deployment that can be run using Docker Compose to integrate
with [Authelia](https://www.authelia.com/).

The only pre-requisites for this is to set hostnames and any required
hosts/DNS entries as required and provide a SSL certicate as "app.pem" as all
traffic must be secured from the browser to HAProxy for authentication to work
correctly (a self signed certificate is sufficient for testing).

Once this is in place, the containers may be started using Docker Compose:

```sh
docker compose up
```
