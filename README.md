# bcrypt_password_hasher

This is a small tool to create base64-encoded bcrypt password hashes.
It is intended to be used for [PollGo!](https://github.com/Top-Ranger/pollgo) and [ResponseGo!](https://github.com/Top-Ranger/responsego), but can be used for other purposes (e.g. for [Caddy](https://caddyserver.com/docs/caddyfile/directives/basicauth)).

## Building

```
go build
```

## Running

```
./bcrypt_password_hash
```

Use ```-help``` to get an overview over all options.

## Licence

Apache-2.0