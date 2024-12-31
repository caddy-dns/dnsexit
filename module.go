package dnsexit

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/dnsexit"
)

// Provider lets Caddy read and manipulate DNS records hosted by DNSExit.
type Provider struct{ *dnsexit.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.dnsexit",
		New: func() caddy.Module { return &Provider{new(dnsexit.Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.Provider.APIKey = caddy.NewReplacer().ReplaceAll(p.Provider.APIKey, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	providername [<api_token>] {
//	    api_token <api_token>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.APIKey = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if p.Provider.APIKey != "" {
					return d.Err("API token already set")
				}
				if d.NextArg() {
					p.Provider.APIKey = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.APIKey == "" {
		return d.Err("missing API token")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
