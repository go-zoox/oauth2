#!/bin/bash

# Docker deployment script for acme-go
# This script deploys SSL certificates to docker containers

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

# Configuration (can be overridden by environment variables)
CONTAINER_NAME="${CONTAINER_NAME:-web-server}"
DOCKER_SSL_DIR="${DOCKER_SSL_DIR:-/etc/ssl/certs}"
RESTART_CONTAINER="${RESTART_CONTAINER:-true}"

echo "Deploying certificate for ${DOMAIN} to container ${CONTAINER_NAME}..."

# Check if container exists
if ! docker ps -a --format "table {{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
    echo "Error: Container ${CONTAINER_NAME} not found!"
    exit 1
fi

# Copy certificate files to container
echo "Copying certificate files to container..."
docker cp "${CERT_PATH}" "${CONTAINER_NAME}:${DOCKER_SSL_DIR}/${DOMAIN}.crt"
docker cp "${KEY_PATH}" "${CONTAINER_NAME}:${DOCKER_SSL_DIR}/${DOMAIN}.key"
docker cp "${FULLCHAIN_PATH}" "${CONTAINER_NAME}:${DOCKER_SSL_DIR}/${DOMAIN}-fullchain.crt"

# Set proper permissions inside container
echo "Setting file permissions..."
docker exec "${CONTAINER_NAME}" chmod 644 "${DOCKER_SSL_DIR}/${DOMAIN}.crt"
docker exec "${CONTAINER_NAME}" chmod 644 "${DOCKER_SSL_DIR}/${DOMAIN}-fullchain.crt"
docker exec "${CONTAINER_NAME}" chmod 600 "${DOCKER_SSL_DIR}/${DOMAIN}.key"

# Restart container if requested
if [ "${RESTART_CONTAINER}" = "true" ]; then
    echo "Restarting container ${CONTAINER_NAME}..."
    docker restart "${CONTAINER_NAME}"
    
    # Wait for container to be ready
    echo "Waiting for container to be ready..."
    sleep 5
    
    # Check if container is running
    if docker ps --format "table {{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
        echo "Container restarted successfully!"
    else
        echo "Warning: Container may not have started properly"
        exit 1
    fi
else
    # Just reload the service inside container (if supported)
    echo "Reloading service inside container..."
    docker exec "${CONTAINER_NAME}" nginx -s reload 2>/dev/null || \
    docker exec "${CONTAINER_NAME}" apache2ctl graceful 2>/dev/null || \
    echo "Service reload not supported, consider restarting the container"
fi

echo "Certificate deployed successfully for ${DOMAIN}!"