
DNSExit module for Caddy
========================

This module adds DNSExit support to Caddy for DNS record operations and ACME DNS-01 challenges.

## Resolver recommendation (important)

For ACME DNS-01, set recursive resolvers explicitly in your Caddy `tls` block:

```caddyfile
tls {
	resolvers 1.1.1.1 8.8.8.8
}
```

Use public recursive resolvers for challenge zone discovery and propagation checks.
Avoid system stub/split-DNS resolver paths (for example `127.0.0.53`) and avoid authoritative DNSExit nameservers in the ACME resolver list.

## Module name

```
dns.providers.dnsexit
```

## Supported configuration

The provider currently supports a single setting:

- `api_token`

Set `DNSEXIT_API_KEY` in your environment (or `.env` when using Docker Compose), and reference it from config. Do not paste secrets directly into the Caddyfile.

## Caddyfile examples

Global ACME DNS provider:

```caddyfile
{
	acme_dns dnsexit {
		api_token {$DNSEXIT_API_KEY}
	}
}

example.com {
	respond "hello"
}
```

Per-site DNS provider:

```caddyfile
example.com {
	tls {
		dns dnsexit {
			api_token {$DNSEXIT_API_KEY}
		}
	}
	respond "hello"
}
```

Multi-site pattern with a shared ACME block:

```caddyfile
(tls_dnsexit_common) {
	tls {
		dns dnsexit {
			api_token {$DNSEXIT_API_KEY}
		}
	}
}

example.com {
	import tls_dnsexit_common
	respond "site one"
}

www.example.com {
	import tls_dnsexit_common
	respond "site two"
}
```

JSON config example:

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "dnsexit",
				"api_token": "{env.DNSEXIT_API_KEY}"
			}
		}
	}
}
```

## Run in Docker (without publishing images)

This repository includes a Dockerfile so users can build their own image locally.

Quick copy-paste example files are available in `examples/docker`.

### Fastest path (example folder)

From this repository root:

```bash
cd examples/docker
cp .env.example .env
# edit .env and Caddyfile for your real domain/token
docker compose up -d --build
docker compose logs -f caddy
```

1. Build image from this repository root:

```bash
docker build -t caddy-dnsexit:local .
```

2. Run container with your Caddyfile and DNSExit token:

```bash
cp examples/docker/.env.example .env
# edit .env and set DNSEXIT_API_KEY
docker run -d --name caddy-dnsexit \
	-p 80:80 -p 443:443 \
	-v "$PWD/Caddyfile:/etc/caddy/Caddyfile:ro" \
	-v caddy_data:/data \
	-v caddy_config:/config \
	--env-file .env \
	caddy-dnsexit:local
```

3. Check logs:

```bash
docker logs -f caddy-dnsexit
```

## Minimal docker-compose example

```yaml
services:
  caddy:
    image: caddy-dnsexit:local
    build: .
    ports:
      - "80:80"
      - "443:443"
		env_file:
			- .env
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy_data:/data
      - caddy_config:/config

volumes:
  caddy_data:
  caddy_config:
```

## Troubleshooting

- `missing API token`: ensure `api_token` is set, or environment variable expansion resolves correctly.
- Secret handling: keep API keys in environment variables / `.env`; avoid committing `.env` files.
- DNS challenge zone detection errors (`REFUSED` or `could not determine zone`):
	- What Caddy does first: before writing the TXT challenge record, Caddy asks DNS for the zone's SOA record to discover the correct parent zone (for example, finding `example.com` from `_acme-challenge.sub.example.com`).
	- Why this can fail: if a configured resolver returns `REFUSED` during that SOA lookup chain, Caddy cannot determine the zone and ACME fails early.
	- Important: this does not necessarily mean the DNSExit API write path is broken; it means the resolver path used for ACME checks is not compatible with this lookup flow.
	- Recommended config: use public recursive resolvers in your Caddy TLS block (for example `1.1.1.1` and `8.8.8.8`) so SOA discovery and propagation checks complete reliably.
	- Good starting config:

```caddyfile
{
	acme_dns dnsexit {
		api_token {$DNSEXIT_API_KEY}
	}
}

notes.megatno.com {
	tls {
		dns dnsexit {
			api_token {$DNSEXIT_API_KEY}
		}
		resolvers 1.1.1.1 8.8.8.8
	}
	reverse_proxy joplin:22300
}
```

	- Avoid using `127.0.0.53` (systemd-resolved stub), split-DNS overlays, or authoritative DNSExit nameservers as ACME resolvers. ACME needs recursive resolution for parent-zone SOA discovery.
	- Why we still call this out: resolver `REFUSED` from DNSExit authoritative nameservers during ACME lookups is not ideal; the workaround is documented here so users can get reliable issuance.
	- Resolver health check (before first issuance):

```bash
# Replace with your real host and zone
HOST="notes.megatno.com"
ZONE="megatno.com"

for R in 1.1.1.1 8.8.8.8; do
  echo "=== @$R ==="
  dig +short @$R SOA "$ZONE"
  dig +short @$R NS "$ZONE"
  dig +noall +authority +comments @$R SOA "_acme-challenge.$HOST"
done
```

	- What healthy output looks like:
		- `SOA <zone>` returns your authoritative zone SOA (for example `ns1.dnsexit.com ...`).
		- `NS <zone>` returns your expected authoritative nameservers.
		- `_acme-challenge.<host>` may be `NXDOMAIN`, but authority should point at your zone (not jump to `com` with an error).

