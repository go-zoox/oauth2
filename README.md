# acme-go

A Go implementation of the popular acme.sh shell script - an ACME client for obtaining SSL/TLS certificates from Let's Encrypt and other ACME-compatible Certificate Authorities.

## Features

- 🔒 **Multiple ACME CAs**: Support for Let's Encrypt, ZeroSSL, SSL.com, and other ACME-compatible CAs
- 🌐 **Multiple Validation Methods**: HTTP-01, DNS-01, and TLS-ALPN-01 challenge types
- 🚀 **DNS Provider Integration**: Support for popular DNS providers with API automation
- 🔄 **Automatic Renewal**: Built-in certificate renewal with configurable schedules
- 📦 **Deployment Hooks**: Flexible deployment system for various services
- 🛠️ **Cross-Platform**: Works on Linux, macOS, and Windows
- ⚡ **High Performance**: Written in Go for better performance and reliability

## Installation

### From Source

```bash
git clone https://github.com/acme-go/acme-client.git
cd acme-client
go build -o acme-go
sudo mv acme-go /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/acme-go/acme-client@latest
```

## Quick Start

### 1. Issue a Certificate (HTTP Validation)

```bash
# Issue certificate using webroot validation
acme-go issue -d example.com -w /var/www/html --email your@email.com

# Issue certificate for multiple domains
acme-go issue -d example.com -d www.example.com -w /var/www/html --email your@email.com
```

### 2. Issue a Certificate (DNS Validation)

```bash
# Issue wildcard certificate using DNS validation
acme-go issue -d "*.example.com" --dns dns_cf --email your@email.com

# Issue certificate using manual DNS
acme-go issue -d example.com --dns manual --email your@email.com
```

### 3. Install Certificate

```bash
# Install certificate for nginx
acme-go install-cert -d example.com \
  --cert-file /etc/nginx/ssl/cert.pem \
  --key-file /etc/nginx/ssl/key.pem \
  --fullchain-file /etc/nginx/ssl/fullchain.pem \
  --reload-cmd "systemctl reload nginx"
```

### 4. List Certificates

```bash
# List all certificates
acme-go list

# List in raw format
acme-go list --raw
```

### 5. Renew Certificates

```bash
# Renew specific certificate
acme-go renew -d example.com

# Renew all certificates
acme-go renew --all

# Force renewal
acme-go renew -d example.com --force
```

## Configuration

acme-go uses a YAML configuration file located at `~/.acme-go.yaml`. Here's an example:

```yaml
email: your@email.com
server: https://acme-v02.api.letsencrypt.org/directory
staging: false
cert_dir: /home/user/.acme-go
key_type: ec-256
renew_days: 30

dns_providers:
  dns_cf:
    CF_API_KEY: your_cloudflare_api_key
    CF_EMAIL: your_cloudflare_email
  dns_ali:
    ALICLOUD_ACCESS_KEY: your_access_key
    ALICLOUD_SECRET_KEY: your_secret_key

deploy_hooks:
  nginx:
    type: script
    script: /usr/local/bin/deploy-nginx.sh
  docker:
    type: script
    script: /usr/local/bin/deploy-docker.sh

notifications:
  enabled: true
  type: email
  settings:
    smtp_server: smtp.gmail.com
    smtp_port: "587"
    username: your@email.com
    password: your_app_password
```

## Supported DNS Providers

- Cloudflare (`dns_cf`)
- Alibaba Cloud (`dns_ali`)
- Amazon Route53 (`dns_aws`)
- Google Cloud DNS (`dns_gcloud`)
- And many more...

## Validation Methods

### HTTP-01 Challenge

```bash
# Using webroot
acme-go issue -d example.com -w /var/www/html

# Using standalone mode (requires port 80)
acme-go issue -d example.com --standalone
```

### DNS-01 Challenge

```bash
# Using DNS provider API
acme-go issue -d example.com --dns dns_cf

# Using manual DNS (you add TXT records manually)
acme-go issue -d example.com --dns manual

# Wildcard certificates (requires DNS validation)
acme-go issue -d "*.example.com" --dns dns_cf
```

### TLS-ALPN-01 Challenge

```bash
# Using TLS-ALPN (requires port 443)
acme-go issue -d example.com --alpn
```

## Key Types

Supported key types:
- `ec-256` (default) - ECDSA P-256
- `ec-384` - ECDSA P-384
- `rsa-2048` - RSA 2048-bit
- `rsa-3072` - RSA 3072-bit
- `rsa-4096` - RSA 4096-bit

```bash
# Issue certificate with RSA 4096-bit key
acme-go issue -d example.com -w /var/www/html --key-type rsa-4096
```

## Deployment Hooks

Deploy certificates automatically to various services:

```bash
# Deploy to nginx
acme-go deploy -d example.com --deploy-hook nginx

# Deploy to docker container
acme-go deploy -d example.com --deploy-hook docker
```

## Automatic Renewal

Set up automatic renewal using cron:

```bash
# Add to crontab (runs daily at 2 AM)
0 2 * * * /usr/local/bin/acme-go renew --all
```

## Commands Reference

### Global Flags

- `--config`: Configuration file path
- `--verbose`: Enable verbose output

### issue

Issue a new certificate.

```bash
acme-go issue [flags]
```

**Flags:**
- `-d, --domain`: Domain name(s) (required)
- `-w, --webroot`: Webroot path for HTTP-01 validation
- `--dns`: DNS provider for DNS-01 validation
- `--key-type`: Key type (default: ec-256)
- `--email`: Email for account registration
- `--server`: ACME server URL
- `--staging`: Use staging environment
- `--force`: Force issue even if certificate exists

### renew

Renew existing certificates.

```bash
acme-go renew [flags]
```

**Flags:**
- `-d, --domain`: Domain name to renew
- `--all`: Renew all certificates
- `--force`: Force renewal even if not due

### install-cert

Install certificate to specified locations.

```bash
acme-go install-cert [flags]
```

**Flags:**
- `-d, --domain`: Domain name (required)
- `--cert-file`: Certificate file path
- `--key-file`: Private key file path
- `--ca-file`: CA certificate file path
- `--fullchain-file`: Full certificate chain file path
- `--reload-cmd`: Command to reload service

### deploy

Deploy certificate using deployment hooks.

```bash
acme-go deploy [flags]
```

**Flags:**
- `-d, --domain`: Domain name (required)
- `--deploy-hook`: Deployment hook name (required)

### list

List all certificates.

```bash
acme-go list [flags]
```

**Flags:**
- `--raw`: Output in raw format

## Comparison with acme.sh

| Feature | acme.sh | acme-go |
|---------|---------|---------|
| Language | Shell Script | Go |
| Performance | Good | Excellent |
| Memory Usage | Low | Low |
| Cross-platform | Yes | Yes |
| DNS Providers | 100+ | Growing |
| Deployment Hooks | 50+ | Growing |
| Configuration | Shell variables | YAML file |
| Error Handling | Basic | Advanced |
| Logging | Basic | Structured |
| Testing | Limited | Comprehensive |

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the GPL v3 License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [acme.sh](https://github.com/acmesh-official/acme.sh) - The original shell script implementation
- [lego](https://github.com/go-acme/lego) - Go ACME client library
- [Let's Encrypt](https://letsencrypt.org/) - Free SSL/TLS certificates for everyone

## Support

- 📖 [Documentation](https://github.com/acme-go/acme-client/wiki)
- 🐛 [Issue Tracker](https://github.com/acme-go/acme-client/issues)
- 💬 [Discussions](https://github.com/acme-go/acme-client/discussions)