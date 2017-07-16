# acadia
Simple home file server

## Setup
Create `~/.acadia.json`

``` json
{
	"serverRoot": "path-to-data-directory"
	"serverKey": "path-to-encryption-key"
	"serverCert": "path-to-encryption-certificate"
}
```

You'll need to a certificate and key. For example,
follow
[these instructions](https://jamielinux.com/docs/openssl-certificate-authority/sign-server-and-client-certificates.html).

## Running
``` bash
godep restore
go build
./acadia
```
