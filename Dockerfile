# ---- Build stage ----
FROM caddy:2.8.4-builder AS builder

# Build Caddy with the dnsexit plugin
RUN xcaddy build \
    --with github.com/caddy-dns/dnsexit

# ---- Final runtime image ----
FROM caddy:2.8.4

# Copy the custom-built Caddy binary
COPY --from=builder /usr/bin/caddy /usr/bin/caddy

# Expose standard Caddy ports
EXPOSE 80 443 2019

# Default command
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile"]
