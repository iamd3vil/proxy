# Proxy

This is a fast TCP/UDP proxy written in Golang with support for TLS in case of TCP.

Docker image can be found at [https://hub.docker.com/r/iamd3vil/proxy](https://hub.docker.com/r/iamd3vil/proxy).

## Usage

```bash
$ ./proxy -config /path/to/config
```

## Configuration

There is a `config.sample.toml` provided in the repo for the sample configuration.

```toml
[proxy]
# Can be either "tcp" or "tls" or "udp"
type = "tcp" 
source = "127.0.0.1:4500"
destination = "127.0.0.1:5500"

# Only required in case of tls
cert = "/path/to/cert"
key = "/path/to/certkey"
```