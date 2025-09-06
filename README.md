# PokéFactory Server

A secure, scalable REST API backend for managing Pokémon progression data in Minecraft NeoForge servers. Tracks player statistics, Pokédex completion, and provides web dashboard access for community features.

## Overview

PokéFactory Server bridges Minecraft gameplay with persistent Pokémon data storage, enabling:
- **Player Management**: UUID-based authentication, statistics tracking, flexible data storage
- **Pokédex System**: National/regional completion tracking with real-time leaderboards
- **Multi-Interface Access**: Secure Minecraft server endpoints + public web dashboard APIs
- **Scalable Deployment**: Docker-based with support for multiple isolated server instances

## Architecture

```
Minecraft Players ↔ NeoForge Server ↔ API (localhost:8080) ↔ PostgreSQL Database
                                    ↔ Web Dashboard (public:8081) ↔ Community Features
```

**Technology Stack:**
- **Backend**: Go 1.21 + Gin web framework
- **Database**: PostgreSQL with automated migrations
- **Authentication**: JWT tokens for server-to-server communication
- **Deployment**: Docker Compose with multi-environment support

## Quick Start

### 1. Production Deployment
```bash
# Create environment file
cat > .env << EOF
DB_NAME=pokefactory_server_1
DB_USER=postgres
DB_PASSWORD=your_secure_password_here
API_PORT=8080
JWT_SECRET=your-super-secret-jwt-key-change-this
EOF

# Deploy secure API (Minecraft server access only)
docker-compose -f docker-compose.prod.yml up -d

# Verify deployment
curl http://localhost:8080/health
```

### 2. Web Dashboard (Optional)
```bash
# Deploy with public web access
WEB_PORT=8081 docker-compose -f docker-compose.web.yml up -d

# Test web endpoints
curl http://localhost:8081/api/v1/web/leaderboards
```

### 3. Minecraft Server Integration
Configure your NeoForge mod to connect to:
```
API Base URL: http://localhost:8080/api/v1/server
Server Key: your-jwt-secret (same as JWT_SECRET)
```

## API Endpoints

### Minecraft Server Endpoints (Authenticated)
- `POST /api/v1/server/auth` - Server authentication
- `POST /api/v1/server/player/create` - Player registration
- `POST /api/v1/server/pokedex/update` - Pokémon catch/seen updates
- `POST /api/v1/server/pokedex/summary` - Player progress retrieval

### Web Dashboard Endpoints (Public)
- `GET /api/v1/web/leaderboards` - Community leaderboards
- `GET /api/v1/web/player/{username}/stats` - Public player stats
- `GET /api/v1/web/server/analytics` - Server-wide analytics
- `GET /api/v1/web/pokemon/{dex}/popularity` - Pokémon popularity data

## Development & Testing

### Local Development
```bash
# Start development environment
docker-compose up -d

# Run test suite
chmod +x test/curl_tests.sh
./test/curl_tests.sh
```

### Testing Without Minecraft
The included test suite simulates Minecraft mod behavior:
- Server authentication flow
- Player creation and Pokémon catching
- Pokédex progress verification
- Web dashboard data validation

## Security Features

- **Database Isolation**: Never exposed to external networks
- **JWT Authentication**: Secure server-to-server communication
- **Localhost Binding**: Production API only accessible via localhost
- **Input Validation**: Comprehensive request validation and sanitization

## Scaling & Multi-Server Support

Deploy multiple isolated instances:
```bash
# Server 1
API_PORT=8080 DB_NAME=pokefactory_server_1 docker-compose -f docker-compose.prod.yml up -d

# Server 2 with web dashboard
API_PORT=8082 WEB_PORT=8083 DB_NAME=pokefactory_server_2 docker-compose -f docker-compose.web.yml -p server2 up -d
```

Each instance maintains completely isolated player data and statistics.

## Documentation

- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Complete deployment guide with DNS/networking setup
- **[minecraft-server-setup.md](minecraft-server-setup.md)** - Minecraft server integration guide
- **[test/README.md](test/README.md)** - Testing and development guide

## Project Structure

```
pokefactory_server/
├── cmd/server/           # Application entry point
├── internal/
│   ├── api/             # HTTP handlers and routing
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and migrations
│   ├── middleware/      # Authentication middleware
│   └── models/          # Data structures
├── migrations/          # Database schema migrations
├── test/               # Test suite and mock clients
└── docker-compose.*.yml # Deployment configurations
```

## License

This project is designed for Minecraft server administrators and mod developers building Pokémon-themed gameplay experiences.