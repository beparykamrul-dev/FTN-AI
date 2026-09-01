# FTN-AI Development Environment Setup

## Quick Start

### Option 1: Docker Compose (Recommended)
```bash
docker-compose up -d
```

### Option 2: Native Development
```bash
# Terminal 1: Start Backend
bash scripts/dev/start-backend.sh

# Terminal 2: Start Frontend
bash scripts/dev/start-frontend.sh
```

### Option 3: Combined Development Server
```bash
bash scripts/dev/dev-server.sh
```

## Access Points

- **Frontend Control Center**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Dashboard**: http://localhost:3000/overview
- **Diagnostics**: http://localhost:8080/diagnostics

## Development Workflow

### Frontend Development
```bash
cd frontend/control-center
npm install
npm run dev
```

### Backend Development
```bash
cd backend
go run ./internal/diagnostics/service.go
```

### Testing
```bash
# Frontend tests
cd frontend/control-center && npm test

# Backend tests
go test ./backend/internal/...
```

## Environment Variables

### Backend
- `GO_ENV=development` - Enable development mode
- `LOG_LEVEL=debug` - Enable debug logging
- `PORT=8080` - Backend port

### Frontend
- `VITE_API_URL=http://localhost:8080` - Backend API URL
- `VITE_WS_URL=ws://localhost:8080` - WebSocket URL
- `NODE_ENV=development` - Development mode

## Troubleshooting

### Port Already in Use
```bash
# Find and kill process on port 3000 (frontend)
lsof -ti:3000 | xargs kill -9

# Find and kill process on port 8080 (backend)
lsof -ti:8080 | xargs kill -9
```

### Dependencies Not Installing
```bash
# Clear npm cache
npm cache clean --force

# Reinstall dependencies
npm install
```

### Go Build Issues
```bash
# Clear Go cache
go clean -cache

# Download dependencies
go mod download
```

## Architecture

```
FTN-AI Development Stack
├── Frontend (React/TypeScript - Port 3000)
│   ├── Control Center UI
│   ├── Dashboard & Visualization
│   └── Real-time Updates (WebSocket)
├── Backend (Go - Port 8080)
│   ├── Diagnostics Service
│   ├── Policy Enforcement
│   └── Health Monitoring
└── Shared
    ├── Config Contracts (YAML)
    ├── API Contracts
    └── Module Registry
```

## Production Deployment

See `docs/DEPLOYMENT_ACCEPTANCE.md` for production deployment procedures.
