// Package core registers the server and all plugins we support.
package core

import (
	// plug in the server
	_ "coredns/core/dnsserver"
	_ "coredns/plugin/nacos"
)
