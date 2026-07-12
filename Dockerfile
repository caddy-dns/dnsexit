ARG CADDY_VERSION=2.11.4

# ---- Build stage ----
FROM caddy:${CADDY_VERSION}-builder AS builder

# Build Caddy with the dnsexit plugin
RUN xcaddy build \
    --with github.com/caddy-dns/dnsexit/v2

# ---- Final runtime image ----
FROM caddy:${CADDY_VERSION}

# Copy the custom-built Caddy binary
COPY --from=builder /usr/bin/caddy /usr/bin/caddy

# Expose standard Caddy ports
EXPOSE 80 443 2019

# Run as non-root user for security
USER caddy

# Default command
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile"]
