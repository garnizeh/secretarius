#!/bin/bash

# TLS Certificate Generation Script for EngLog Development
# This script generates self-signed certificates for development use

set -e

CERT_DIR="certs"
CERT_FILE="$CERT_DIR/server.crt"
KEY_FILE="$CERT_DIR/server.key"

echo "🔐 Generating TLS certificates for EngLog development..."

# Create certs directory if it doesn't exist
mkdir -p "$CERT_DIR"

# Generate private key and certificate
openssl req -x509 -newkey rsa:4096 \
    -keyout "$KEY_FILE" \
    -out "$CERT_FILE" \
    -days 365 \
    -nodes \
    -subj "/C=BR/ST=SP/L=São Paulo/O=EngLog/OU=Development/CN=localhost/emailAddress=dev@englog.local"

# Set appropriate permissions
chmod 600 "$KEY_FILE"
chmod 644 "$CERT_FILE"

echo "✅ Certificates generated successfully!"
echo "📁 Certificate: $CERT_FILE"
echo "🔑 Private Key: $KEY_FILE"
echo ""

# Display certificate info
echo "📋 Certificate Details:"
openssl x509 -in "$CERT_FILE" -subject -dates -noout

echo ""
echo "⚠️  Note: This is a self-signed certificate for development only."
echo "   For production, use certificates from a trusted CA."
echo ""
echo "🚀 You can now start the gRPC server with TLS enabled!"
