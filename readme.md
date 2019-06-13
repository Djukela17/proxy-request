# Proxy Request

Simple request proxy

## Requirements
- `go` obviously (tested on version 1.12, but should work older versions)


## Setup

### Linux:

Tested on:
- ubuntu 18.04 bionic
- debian 9 stretch

Probably works on many others too

#### Preparing:

- Create your config file containing whitelisted IP addresses 
(default one is `default.config`), based on the `example.config` file

**Example config file:**
```
# Allowed IPs
125.234.432.43
```

- Build the binary using `go build` => will produce `proxy-request` binary

#### Starting the server:

```sh
./proxy-request -host :8080 -config my-file.config
```
**Flags:**
- `-host` - the address and port the server will listen on
- `-config` - the filepath to the config file containing whitelisted IPs


## Usage

For example, let's say server is running on remote with domain name: `my-domain-name.com`
Server is running on default port: `80`.

#### Request format
```
[server-address]/[protocol]:[url-with-params-to-access]
```

**Example:**
`my-domain.com/https:api.another-site.com?foo=bar`

### Accessing github API
Accessing `https://api.github.com/users` through `proxy-request` that is running
on remote server with domain name `my-domain-name.com`:

* Create a `GET` request in the following format:
```sh
curl my-domain-name.com/https:api.github.com/users/
```
* You will receive a regular github API response with a list of users

The same goes for other request types (`POST`, `PUT`, `DELETE`)


## Limitations

### Server

* The server cannot be started in TLS mode

### Requests

* Headers from the original request will not pe passed to the repeated request

### Config file

* Does not support IP blacklisting, only whitelisting
* The line must contain only IP or comment, no mixing