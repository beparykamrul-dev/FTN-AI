# FTN-AI Architecture

## System Overview

```
┌────────────────────────────────────────────────────────────────────────┐
│                    Client Applications                       │
│              (Web UI, Mobile, CLI, APIs)                    │
└────────────────────────┬─────────────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────────────────────────────┐
│                  API Gateway & Load Balancer                │
│              (Authentication, Rate Limiting)                │
└────────────────────────┬─────────────────────────────────────┘
                   │
      ┌────────────┼─────────────┐
      ▼            ▼            ▼
 ┌──────────┐ ┌──────────┐ ┌──────────┐
 │ Account  │ │   AI    │ │ Storage │
 │ Service  │ │ Service │ │ Service │
 └────┬─────┘ └────┬─────┘ └────┬─────┘
      │           │           │
      └───────────┼───────────┘
              │
              ▼
    ┌─────────────────────────┐
    │ Shared Data Layer       │
    │  (PostgreSQL, Redis)   │
    └─────────────────────────┘
```

## Service Architecture

### Backend Services

#### Account Service
- User management
- Authentication/Authorization
- Profile management
- Role-based access control (RBAC)

#### AI Service
- Model inference
- AI governance enforcement
- Tool registration and execution
- Approval workflow management

#### Storage Service
- File management
- Media handling
- Object storage integration

#### Notification Service
- Real-time updates (WebSocket)
- Email notifications
- Audit logging

### Data Persistence

#### PostgreSQL
- Primary database
- User data
- Service configurations
- Audit logs
- Transaction safety

#### Redis
- Caching layer
- Session management
- Real-time subscriptions
- Rate limiting

## Security Architecture

### Authentication Flow
1. User submits credentials
2. Backend validates against database
3. JWT token generated
4. Token stored in Redis with TTL
5. Token included in subsequent requests

### Authorization
- Role-based access control (RBAC)
- Resource-level permissions
- Audit trail of all access

### AI Governance
- All AI actions tracked
- High-risk operations require approval
- Automatic logging of decisions
- Secure credential handling

## Deployment Architecture

### Development
- Docker Compose for local development
- All services containerized
- Volume mounting for hot reload

### Production
- Kubernetes orchestration
- Auto-scaling based on metrics
- High availability setup
- Blue-green deployments

## Data Flow

```
Client Request
    │
API Gateway (Auth/Rate Limit)
    │
Route to appropriate service
    │
Service processes request
    │
Query data layer if needed
    │
Cache layer checks
    │
Database query execution
    │
Response construction
    │
Client receives response
```