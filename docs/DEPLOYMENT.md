# FTN-AI Deployment Guide

## Prerequisites

- Docker & Docker Compose 2.0+
- Kubernetes 1.24+ (for production)
- PostgreSQL 15+
- Go 1.24+
- Node.js 18+

## Development Deployment

### Using Docker Compose

```bash
# Clone repository
git clone https://github.com/beparykamrul-dev/FTN-AI.git
cd FTN-AI

# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Local Development

```bash
# Backend
cd backend
go run main.go

# Frontend (in another terminal)
cd frontend
npm start
```

## Production Deployment

### Prerequisites

1. Set up environment variables:
```bash
cp configs/v1/.env.example configs/v1/.env.production
# Edit with production values
```

2. Configure database:
```bash
psql -U postgres -d ftn_ai -f internal/database/migrations.sql
```

3. Build Docker images:
```bash
docker build -t ftn-backend:1.0.0 .
docker build -t ftn-frontend:1.0.0 ./frontend
```

### Kubernetes Deployment

```bash
# Create namespace
kubectl create namespace ftn-ai

# Create secrets
kubectl create secret generic ftn-secrets \
  --from-env-file=configs/v1/.env.production \
  -n ftn-ai

# Deploy services
kubectl apply -f k8s/
```

## Monitoring and Logging

### Prometheus Metrics

```bash
# Access at http://localhost:9090
```

### Logging

Logs are collected to:
- `/var/log/ftn-ai/`
- Container logs: `docker logs <container>`

## Backup and Recovery

### Database Backup

```bash
pg_dump -U postgres ftn_ai > backup_$(date +%Y%m%d).sql
```

### Restore

```bash
psql -U postgres ftn_ai < backup_20240818.sql
```

## Health Checks

```bash
# API health
curl http://localhost:8080/health

# Database health
psql -U postgres -d ftn_ai -c "SELECT 1"
```