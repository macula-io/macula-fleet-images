// caddy with the Linode DNS provider, for the stations' ACME DNS-01 challenge.
// The module set is fixed by go.mod and go.sum beside this file: no xcaddy, which
// resolves module versions at build time and records nothing.
package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	_ "github.com/caddy-dns/linode"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
)

func main() {
	caddycmd.Main()
}
