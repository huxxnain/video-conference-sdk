# Video Conference SDK

A production-ready video conferencing SDK with WebRTC signaling, user authentication, room management, and queue functionality.

## Features

- **Backend (Go)**
  - JWT-based authentication
  - PostgreSQL database with GORM
  - WebRTC signaling via WebSocket
  - Room creation and management
  - Queue system for joining rooms
  - CORS support
  - Health check endpoint
  
- **Frontend (React)**
  - Basic React application scaffold
  - Nginx reverse proxy for API and WebSocket
  
- **DevOps**
  - Docker Compose orchestration
  - PostgreSQL database container
  - Production-ready Dockerfiles
  - Easy local development setup

## Prerequisites

- Docker and Docker Compose installed
- (Optional) Go 1.21+ for local backend development
- (Optional) Node.js 18+ for local frontend development

## Quick Start with Docker Compose

1. Clone the repository:
```bash
git clone https://github.com/huxxnain/video-conference-sdk.git
cd video-conference-sdk
```

2. Start all services:
```bash
docker-compose up -d
```

3. Access the services:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Health Check: http://localhost:8080/health
   - PostgreSQL: localhost:5432

4. Stop all services:
```bash
docker-compose down
```

5. Clean up (including volumes):
```bash
docker-compose down -v
```

## API Endpoints

### Public Endpoints

**POST /auth/signup**
```json
{
  "org_name": "My Organization",
  "email": "user@example.com",
  "password": "securepassword"
}
```

**POST /auth/login**
```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**GET /health**
Returns service health status.

### Protected Endpoints (Requires JWT)

**POST /room/create**
```json
{
  "org_id": 1,
  "name": "Conference Room 1"
}
```
Headers: `Authorization: Bearer <token>`

**POST /room/join**
```json
{
  "user_id": 1,
  "room_id": 1
}
```
Headers: `Authorization: Bearer <token>`

### WebSocket Endpoint

**GET /ws/signaling?room=<room_name>**
WebRTC signaling endpoint for real-time communication.

## Local Development

### Backend Development

1. Navigate to backend directory:
```bash
cd backend
```

2. Copy and configure environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Install dependencies:
```bash
go mod download
```

4. Run the backend:
```bash
go run cmd/main.go
```

### Frontend Development

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start development server:
```bash
npm start
```

The app will open at http://localhost:3000

## Configuration

### Backend Environment Variables

See `backend/.env.example` for all available configuration options:

- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT tokens
- `ALLOWED_ORIGINS`: CORS allowed origins (comma-separated)
- `ENVIRONMENT`: development/production

### Database Schema

The application automatically migrates the following tables:
- `organizations`: Organization entities
- `users`: User accounts with authentication
- `rooms`: Video conference rooms
- `queue_entries`: Queue for joining rooms

## Architecture

```
┌─────────────┐
│   Frontend  │ (React + Nginx)
│  Port: 3000 │
└──────┬──────┘
       │
       │ HTTP/WebSocket
       │
┌──────▼──────┐
│   Backend   │ (Go + Gin)
│  Port: 8080 │
└──────┬──────┘
       │
       │ PostgreSQL
       │
┌──────▼──────┐
│  PostgreSQL │
│  Port: 5432 │
└─────────────┘
```

## Security Considerations

- Change `JWT_SECRET` in production
- Use strong database passwords
- Configure CORS `ALLOWED_ORIGINS` appropriately
- Enable SSL/TLS in production
- Use environment-specific `.env` files
- Never commit `.env` files to version control

## Testing

### Test Backend Locally

```bash
cd backend
go test ./...
```

### Test with curl

1. Sign up:
```bash
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"org_name":"TestOrg","email":"test@example.com","password":"password123"}'
```

2. Login:
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

3. Create room (with token):
```bash
curl -X POST http://localhost:8080/room/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{"org_id":1,"name":"Room 1"}'
```

## Troubleshooting

### Database Connection Issues
- Ensure PostgreSQL container is running: `docker-compose ps`
- Check database logs: `docker-compose logs postgres`
- Verify DATABASE_URL environment variable

### Backend Not Starting
- Check backend logs: `docker-compose logs backend`
- Ensure database is healthy before backend starts
- Verify all required environment variables are set

### Frontend Build Issues
- Clear npm cache: `npm cache clean --force`
- Remove node_modules: `rm -rf node_modules && npm install`
- Check frontend logs: `docker-compose logs frontend`

## License

MIT

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
