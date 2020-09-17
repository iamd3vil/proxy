# Proxy

This is a fast TCP/TLS proxy written in Golang with support for Automatic certificates for TLS using Letsencrypt. :boom: :boom:

`proxy` uses `dns-01` challenge for Letsencrypt. Currently `proxy` supports only Cloudflare but the plan is to add a lot of DNS providers.

Docker image can be found at [https://hub.docker.com/r/iamd3vil/proxy](https://hub.docker.com/r/iamd3vil/proxy).

## Usage

```bash
$ ./proxy -config /path/to/config
```

## Configuration

There is a `config.sample.toml` provided in the repo for the sample configuration.

```toml
[proxy]
type = "tls" # has to be either "tcp" or "tls"
source = "127.0.0.1:4500"
destination = "127.0.0.1:8888"

[tls]
disable_automatic = false # make this true to turn off automatic certs fetching and renewals from Letsencrypt
domain = "example.com" # domains for automatic https
certs_path = "/path/to/store/certs" # path where the automatic letsencrypt certs are stored
cloudflare_api_token = "auth-token" # API token for cloudflare
email = "hello@example.com" # Email for letsencrypt

# Only set these if you want to provide certs manually, i.e disable_automatic is true.
cert = "/path/to/cert"
key = "/path/to/key"
```
