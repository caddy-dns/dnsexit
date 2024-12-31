
---

DNSExit module for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with DNSExit.

## Caddy module name

```
dns.providers.dnsexit
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "dnsexit",
				"api_token": "YOUR_DNSEXIT_API_TOKEN"
			}
		}
	}
}
```

or with the Caddyfile:

```
# globally
{
	acme_dns dnsexit ...
}
```

```
# one site
tls {
	dns dnsexit ...
}
```
