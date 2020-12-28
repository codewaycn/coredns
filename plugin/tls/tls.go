package tls

import (
	"coredns/core/dnsserver"
	"coredns/plugin"
	"coredns/plugin/pkg/tls"

	"github.com/caddyserver/caddy"
)

func init() {
	caddy.RegisterPlugin("tls", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	config := dnsserver.GetConfig(c)

	if config.TLSConfig != nil {
		return plugin.Error("tls", c.Errf("TLS already configured for this server instance"))
	}

	for c.Next() {
		args := c.RemainingArgs()
		if len(args) < 2 || len(args) > 3 {
			return plugin.Error("tls", c.ArgErr())
		}
		tls, err := tls.NewTLSConfigFromArgs(args...)
		if err != nil {
			return plugin.Error("tls", err)
		}
		config.TLSConfig = tls
	}
	return nil
}
