# registering user-service using consul discovery service

* [consul](https://www.consul.io/)
* [user-service](https://github.com/mariiatuzovska/consul-discovery-service/blob/master/user-service)

## Registration

1. Building *user-service* `$make user` in a new terminal.
2. Starting *user-service* `$make user-start` (available at the following address *127.0.0.1:8181*).
3. Startning *consul agent* in development mode `$make consul-start` in a new terminal.
4. Connecting to *user-service* `$make consul-connect-user` in a new terminal.
5. Connecting to *web* `$make consul-connect-web` (now, *user-service* is available at the following address *127.0.0.1:9191*) in a new terminal.

## Access

1. Creating an intention to deny access from web to user `$make intention-create` in a new terminal. The connection to *127.0.0.1:9191* will fail.
2. Deleting the intention `$make intention-delete`.

## Health checks 

1. `$make user-health`
2. `$make web-health`