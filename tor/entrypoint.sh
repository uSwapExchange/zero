#!/bin/bash
set -e

HS_DIR=/var/lib/tor/hidden_service
mkdir -p "$HS_DIR"
chmod 700 "$HS_DIR"

# Inject key material from env vars (base64-encoded)
if [ -n "$TOR_SECRET_KEY_B64" ]; then
    printf '%s' "$TOR_SECRET_KEY_B64" | base64 -d > "$HS_DIR/hs_ed25519_secret_key"
    printf '%s' "$TOR_PUBLIC_KEY_B64" | base64 -d > "$HS_DIR/hs_ed25519_public_key"
    printf '%s\n' "$TOR_HOSTNAME"     > "$HS_DIR/hostname"
    chmod 600 "$HS_DIR/"*
    echo "Tor: loaded hidden service key for $TOR_HOSTNAME"
else
    echo "Tor: no key provided — Tor will generate a fresh address"
fi

exec tor -f /etc/tor/torrc
