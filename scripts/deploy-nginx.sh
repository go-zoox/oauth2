#!/bin/bash

# Nginx deployment script for acme-go
# This script deploys SSL certificates to nginx configuration

set -e

# Environment variables passed by acme-go:
# ACME_DOMAIN - The domain name
# ACME_CERT_PATH - Path to certificate file
# ACME_KEY_PATH - Path to private key file
# ACME_FULLCHAIN_PATH - Path to full certificate chain
# ACME_CA_PATH - Path to CA certificate

DOMAIN="${ACME_DOMAIN}"
CERT_PATH="${ACME_CERT_PATH}"
KEY_PATH="${ACME_KEY_PATH}"
FULLCHAIN_PATH="${ACME_FULLCHAIN_PATH}"

# Configuration
NGINX_SSL_DIR="/etc/nginx/ssl"
NGINX_CONFIG_DIR="/etc/nginx/sites-available"
NGINX_ENABLED_DIR="/etc/nginx/sites-enabled"

# Ensure SSL directory exists
mkdir -p "${NGINX_SSL_DIR}"

# Copy certificate files
echo "Deploying certificate for ${DOMAIN}..."
cp "${CERT_PATH}" "${NGINX_SSL_DIR}/${DOMAIN}.crt"
cp "${KEY_PATH}" "${NGINX_SSL_DIR}/${DOMAIN}.key"
cp "${FULLCHAIN_PATH}" "${NGINX_SSL_DIR}/${DOMAIN}-fullchain.crt"

# Set proper permissions
chmod 644 "${NGINX_SSL_DIR}/${DOMAIN}.crt"
chmod 644 "${NGINX_SSL_DIR}/${DOMAIN}-fullchain.crt"
chmod 600 "${NGINX_SSL_DIR}/${DOMAIN}.key"

# Create nginx configuration if it doesn't exist
NGINX_CONF="${NGINX_CONFIG_DIR}/${DOMAIN}"
if [ ! -f "${NGINX_CONF}" ]; then
    echo "Creating nginx configuration for ${DOMAIN}..."
    cat > "${NGINX_CONF}" << EOF
server {
    listen 80;
    server_name ${DOMAIN};
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name ${DOMAIN};

    ssl_certificate ${NGINX_SSL_DIR}/${DOMAIN}-fullchain.crt;
    ssl_certificate_key ${NGINX_SSL_DIR}/${DOMAIN}.key;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-SHA256:ECDHE-RSA-AES256-SHA384;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Content-Type-Options nosniff;
    add_header X-Frame-Options DENY;
    add_header X-XSS-Protection "1; mode=block";

    root /var/www/${DOMAIN};
    index index.html index.htm index.php;

    location / {
        try_files \$uri \$uri/ =404;
    }

    # PHP support (uncomment if needed)
    # location ~ \.php$ {
    #     include snippets/fastcgi-php.conf;
    #     fastcgi_pass unix:/var/run/php/php7.4-fpm.sock;
    # }
}
EOF

    # Enable the site
    ln -sf "${NGINX_CONF}" "${NGINX_ENABLED_DIR}/${DOMAIN}"
fi

# Test nginx configuration
echo "Testing nginx configuration..."
nginx -t

# Reload nginx
echo "Reloading nginx..."
systemctl reload nginx

echo "Certificate deployed successfully for ${DOMAIN}!"