# Deployment Guide

This directory contains deployment artifacts for the Go Web Server.

## Files

- **Caddyfile**: Caddy web server configuration for reverse proxy and SSL termination
- **go-web-server.service**: systemd service file for process management
- **README.md**: This deployment guide

## Quick Deployment Guide

### 1. Build the Application

```bash
make build
```

### 2. Deploy to Server

Copy your binary and files to the server:

```bash
# Create application directory
sudo mkdir -p /opt/go-web-server/{bin,data}

# Copy binary
sudo cp bin/server /opt/go-web-server/bin/

# Set permissions
sudo chown -R www-data:www-data /opt/go-web-server
sudo chmod +x /opt/go-web-server/bin/server
```

### 3. Install systemd Service

```bash
# Copy service file
sudo cp deploy/go-web-server.service /etc/systemd/system/

# Reload systemd and enable service
sudo systemctl daemon-reload
sudo systemctl enable go-web-server
sudo systemctl start go-web-server

# Check status
sudo systemctl status go-web-server
```

### 4. Install and Configure Caddy

#### Install Caddy:
```bash
# Debian/Ubuntu
curl -fsSL https://dl.cloudsmith.io/public/caddy/stable/gpg.key | sudo tee /etc/apt/trusted.gpg.d/caddy-stable.asc
echo "deb [signed-by=/etc/apt/trusted.gpg.d/caddy-stable.asc] https://dl.cloudsmith.io/public/caddy/stable/deb/debian any-version main" | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy
```

#### Configure Caddy:
```bash
# Backup default config
sudo mv /etc/caddy/Caddyfile /etc/caddy/Caddyfile.backup

# Copy your configuration (update your-domain.com first!)
sudo cp deploy/Caddyfile /etc/caddy/

# Test configuration
sudo caddy validate --config /etc/caddy/Caddyfile

# Restart Caddy
sudo systemctl restart caddy
sudo systemctl enable caddy
```

### 5. Database Setup

If using SQLite (default):
```bash
# Ensure database directory exists
sudo mkdir -p /opt/go-web-server/data
sudo chown www-data:www-data /opt/go-web-server/data
```

### 6. Environment Configuration

Create environment file (optional):
```bash
sudo tee /opt/go-web-server/.env << 'EOF'
ENVIRONMENT=production
PORT=8080
HOST=127.0.0.1
DATABASE_URL=/opt/go-web-server/data/app.db
LOG_LEVEL=info
LOG_FORMAT=json
EOF
```

### 7. Monitoring

Check application logs:
```bash
# Application logs
sudo journalctl -u go-web-server -f

# Caddy logs
sudo journalctl -u caddy -f

# Or check Caddy access logs
sudo tail -f /var/log/caddy/your-domain.com.log
```

## Security Considerations

1. **Firewall**: Ensure only ports 80 and 443 are open to the public
2. **SSL**: Caddy automatically handles SSL certificates via Let's Encrypt
3. **User permissions**: The service runs as `www-data` with restricted permissions
4. **Database**: Ensure SQLite file has proper permissions if using SQLite
5. **Updates**: Regularly update both your application and system packages

## Backup Strategy

```bash
# Database backup (SQLite)
sudo cp /opt/go-web-server/data/app.db /opt/go-web-server/data/app.db.backup.$(date +%Y%m%d_%H%M%S)

# Full application backup
sudo tar -czf go-web-server-backup-$(date +%Y%m%d).tar.gz -C /opt go-web-server/
```

## Rolling Updates

```bash
# Build new version
make build

# Stop service
sudo systemctl stop go-web-server

# Backup current binary
sudo cp /opt/go-web-server/bin/server /opt/go-web-server/bin/server.backup

# Deploy new binary
sudo cp bin/server /opt/go-web-server/bin/

# Start service
sudo systemctl start go-web-server

# Check status
sudo systemctl status go-web-server
```

## Troubleshooting

- Check service status: `sudo systemctl status go-web-server`
- View logs: `sudo journalctl -u go-web-server -n 50`
- Test binary directly: `sudo -u www-data /opt/go-web-server/bin/server`
- Check port binding: `sudo netstat -tlnp | grep :8080`