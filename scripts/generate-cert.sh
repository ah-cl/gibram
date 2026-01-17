#!/bin/sh
# Generate self-signed TLS certificate for GibRAM server

set -e

CERT_DIR="${CERT_DIR:-/etc/gibram/certs}"
DOMAIN="${DOMAIN:-localhost}"
DAYS="${DAYS:-365}"

echo "Generating self-signed TLS certificate..."
echo "  Domain: $DOMAIN"
echo "  Valid for: $DAYS days"
echo "  Output directory: $CERT_DIR"
echo ""

# Create directory if it doesn't exist
if [ -w "/etc/gibram" ]; then
    mkdir -p "$CERT_DIR"
else
    sudo mkdir -p "$CERT_DIR"
fi

# Generate certificate
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

openssl req -x509 -newkey rsa:4096 -nodes \
    -keyout "$TEMP_DIR/server.key" \
    -out "$TEMP_DIR/server.crt" \
    -days "$DAYS" \
    -subj "/CN=$DOMAIN" \
    -addext "subjectAltName=DNS:$DOMAIN,DNS:localhost,IP:127.0.0.1"

# Install certificate
if [ -w "$CERT_DIR" ]; then
    mv "$TEMP_DIR/server.key" "$CERT_DIR/server.key"
    mv "$TEMP_DIR/server.crt" "$CERT_DIR/server.crt"
    chmod 600 "$CERT_DIR/server.key"
    chmod 644 "$CERT_DIR/server.crt"
else
    sudo mv "$TEMP_DIR/server.key" "$CERT_DIR/server.key"
    sudo mv "$TEMP_DIR/server.crt" "$CERT_DIR/server.crt"
    sudo chmod 600 "$CERT_DIR/server.key"
    sudo chmod 644 "$CERT_DIR/server.crt"
fi

echo ""
echo "✓ Certificate generated successfully!"
echo ""
echo "Certificate: $CERT_DIR/server.crt"
echo "Private key: $CERT_DIR/server.key"
echo ""
echo "Update your config.yaml:"
echo "  tls:"
echo "    cert_file: \"$CERT_DIR/server.crt\""
echo "    key_file: \"$CERT_DIR/server.key\""
echo "    auto_cert: false"
echo ""
echo "Then restart gibram-server"
