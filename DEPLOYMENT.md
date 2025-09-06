# PokéFactory Server Deployment Guide

## Quick Start for Partners

### 1. Download Files
```bash
# Download these files to your server:
# - docker-compose.prod.yml
# - .env (create from template below)
```

### 2. Create Environment File
```bash
# Create .env file:
DB_NAME=pokefactory_server_1
DB_USER=postgres
DB_PASSWORD=your_secure_password_here
API_PORT=8080
JWT_SECRET=your-super-secret-jwt-key-change-this
```

### 3. Deploy
```bash
# Pull and start the server
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f backend
```

### 4. Test API
```bash
# Health check
curl http://localhost:8080/health

# Should return: {"status":"healthy","time":"..."}
```

## Web Dashboard Deployment

For web dashboard access, use the web-enabled configuration:

```bash
# Deploy with web access enabled
WEB_PORT=8081 docker-compose -f docker-compose.web.yml up -d
```

**Port Configuration:**
- `localhost:8080` - Minecraft server access (secure)
- `your-domain:8081` - Web dashboard access (public)

## Multiple Server Instances

For multiple Minecraft servers on same machine:

### Server 1 (Minecraft only)
```bash
API_PORT=8080 DB_NAME=pokefactory_server_1 docker-compose -f docker-compose.prod.yml up -d
```

### Server 2 (With web dashboard)
```bash
API_PORT=8082 WEB_PORT=8083 DB_NAME=pokefactory_server_2 docker-compose -f docker-compose.web.yml -p server2 up -d
```

## Minecraft Server Configuration

Point your NeoForge mod to:
- Server 1: `http://localhost:8080/api/v1/server`
- Server 2: `http://localhost:8081/api/v1/server`

## DNS & Network Configuration for Web Dashboard

### Step 1: Router/Firewall Setup
```bash
# Open port 8081 in firewall
sudo ufw allow 8081

# Forward port 8081 in router to your machine's local IP
# Router settings: External Port 8081 → Internal IP:8081
```

### Step 2: DNS Configuration Options

**Option A: Use Dynamic DNS (Easiest)**
```bash
# Services like No-IP, DuckDNS, or Cloudflare
# Point yourserver.duckdns.org → your public IP
# Webapp calls: http://yourserver.duckdns.org:8081/api/v1/web/
```

**Option B: Domain with A Record**
```dns
# DNS Provider (Cloudflare, Namecheap, etc.)
api.yourdomain.com    A    YOUR_PUBLIC_IP

# Webapp calls: http://api.yourdomain.com:8081/api/v1/web/
```

**Option C: Reverse Proxy (Professional)**
```nginx
# nginx config for api.yourdomain.com
server {
    listen 80;
    server_name api.yourdomain.com;
    
    location /api/v1/web/ {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Option D: Direct IP Access (Testing)**
```bash
# Find your public IP
curl ifconfig.me

# Webapp calls: http://YOUR_PUBLIC_IP:8081/api/v1/web/
```

### Step 3: Test External Access
```bash
# From external network, test:
curl http://yourserver.duckdns.org:8081/api/v1/web/leaderboards
```

**Web API Endpoints:**
- `GET /api/v1/web/leaderboards`
- `GET /api/v1/web/player/{username}/stats`
- `GET /api/v1/web/server/analytics`
- `GET /api/v1/web/pokemon/{dex}/popularity`

### Security Considerations

**CORS Configuration (if webapp on different domain):**
```bash
# Add to .env file:
CORS_ORIGINS=https://yourwebapp.com,http://localhost:3000
```

**Recommended: Use HTTPS**
```bash
# Install Let's Encrypt SSL certificate
sudo certbot --nginx -d api.yourdomain.com
```

**Rate Limiting (Production)**
```bash
# Limit API calls per IP
# Built into the Go backend - no additional config needed
```

## Troubleshooting

**Check if containers are running:**
```bash
docker ps
```

**View backend logs:**
```bash
docker logs pokefactory_backend
```

**Test web endpoints:**
```bash
curl http://localhost:8081/api/v1/web/leaderboards
```

**Restart services:**
```bash
docker-compose -f docker-compose.web.yml restart
```