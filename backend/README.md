# Task Management API 

A minimalistic task management API with CRUD operations.

## 🚀 Quick Overview

This API provides a clean and efficient way to manage tasks through RESTful endpoints. It supports comprehensive task operations with advanced features:

- **Tasks**: Task management with status tracking and priority levels
- **CRUD Operations**: Full create, read, update, delete functionality
- **Advanced Features**: Pagination, filtering, sorting, and rate limiting
- **Performance**: High-performance caching with Ristretto
- **Security**: JWT authentication with role-based access control
- **Monitoring**: Built-in health checks and Swagger documentation

## 🛠 Tech Stack

- **Backend**: Go 1.23
- **Web Framework**: Gin
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Cache**: Ristretto (high-performance)
- **Authentication**: JWT with refresh tokens
- **Testing**: Go testing framework with SQLite
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose

## Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Environment variables configured

### Environment Variables

Create a `.env` file in the backend directory:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskmanager
DB_SSLMODE=disable
```

### Database Setup

1. Make sure PostgreSQL is running
2. Create the database:
   ```sql
   CREATE DATABASE taskmanager;
   ```

3. Run migrations:
   ```bash
   cd backend
   go run scripts/migrate.go up
   ```

### Running the Application

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Start the server:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token

### Tasks (Protected)
- `GET /api/v1/tasks` - Get all tasks (with pagination, filtering, sorting)
- `POST /api/v1/tasks` - Create new task
- `GET /api/v1/tasks/:id` - Get task by ID
- `PUT /api/v1/tasks/:id` - Update task
- `DELETE /api/v1/tasks/:id` - Delete task

### Users (Protected)
- `GET /api/v1/users/profile` - Get current user profile
- `GET /api/v1/users/profile/:user_id` - Get user profile by ID
- `GET /api/v1/users/:user_id/tasks` - Get tasks by user ID
- `GET /api/v1/users` - Get all users (admin only)
- `DELETE /api/v1/users/:user_id` - Delete user (admin only)

### System
- `GET /health` - Health check endpoint
- `GET /swagger/index.html` - API documentation

## Data Models

### Task
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Important Task",
  "description": "Complete the project documentation",
  "status": "pending",
  "priority": "high",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Create Task Request
```json
{
  "title": "New Task",
  "description": "Task description",
  "priority": "high",
  "status": "pending"
}
```

### Update Task Request
```json
{
  "title": "Updated Task",
  "description": "Updated description",
  "priority": "medium",
  "status": "in_progress"
}
```

### Paginated Response
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 100,
    "total_pages": 10,
    "has_next": true,
    "has_prev": false
  }
}
```

## 📋 Prerequisites

### For Local Development (without Docker)
- **Go 1.23** or higher
- **PostgreSQL 15** or higher
- Git

### For Docker Deployment
- **Docker** 20.10+
- **Docker Compose** 2.0+

## 🚀 Getting Started

Choose your preferred setup method:

- [**Local Development**](#-local-development-setup) - Run directly on your machine
- [**Docker Development**](#-docker-development-setup) - Run with Docker Compose

---

## 💻 Local Development Setup

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/task_management_api.git
cd task_management_api/backend
```

### 2. Install Go Dependencies
```bash
go mod download
```

### 3. Environment Configuration
```bash
# Copy environment template
cp .env.example .env

# Edit with your database settings
nano .env
```

**Required environment variables:**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskmanager
DB_SSLMODE=disable
SERVER_PORT=8080
```

### 4. Database Setup
```bash
# Start PostgreSQL service
sudo systemctl start postgresql

# Create database
psql -U postgres -c "CREATE DATABASE taskmanager;"
```

### 5. Run the Application
```bash
# Start the API server
go run main.go
```

✅ **API available at:** `http://localhost:8080`
✅ **Swagger docs at:** `http://localhost:8080/swagger/index.html`

---

## 🐳 Docker Development Setup

### Quick Start (Recommended)

```bash
# Clone the repository
git clone https://github.com/your-username/task_management_api.git
cd task_management_api

# Start all services with Docker Compose
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f backend
```

✅ **API available at:** `http://localhost:8080`
✅ **Database available at:** `localhost:5432`

## 📖 API Usage Examples

### Create a Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "title": "Complete Project",
    "description": "Finish the task management system",
    "priority": "high",
    "status": "pending"
  }'
```

### Get All Tasks (with pagination)
```bash
curl "http://localhost:8080/api/v1/tasks?page=1&page_size=10&sort_by=created_at&sort_order=desc" \
  -H "Authorization: Bearer <your-token>"
```

### Get Tasks with Filtering
```bash
# Filter by status
curl "http://localhost:8080/api/v1/tasks?status=pending" \
  -H "Authorization: Bearer <your-token>"

# Filter by priority
curl "http://localhost:8080/api/v1/tasks?priority=high" \
  -H "Authorization: Bearer <your-token>"

# Search by title
curl "http://localhost:8080/api/v1/tasks?search=project" \
  -H "Authorization: Bearer <your-token>"
```

### Update a Task
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "title": "Updated Task Title",
    "status": "in_progress",
    "priority": "medium"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### Health Check
```bash
curl http://localhost:8080/health
```


## Default Admin User

The system comes with a default admin user:
- Username: `admin`
- Email: `admin@gmail.com`
- Password: `admin123`

## Database Schema

The application uses the following main tables:
- `users` - User accounts
- `roles` - Role definitions
- `user_roles` - User-role associations
- `permissions` - Permission definitions
- `role_permissions` - Role-permission associations
- `tasks` - Task data
- `tokens` - Refresh tokens

## 🧪 Testing

### Run Unit Tests
```bash
go test ./...
```

### Run Integration Tests
```bash
go test ./test/integrations/...
```

### Run with Coverage
```bash
go test -cover ./...
```

### Run Specific Test Package
```bash
go test ./internal/services/...
```

## 📊 Monitoring & Profiling

### Health Check
```bash
curl http://localhost:8080/health
```

### API Documentation
Visit `http://localhost:8080/swagger/index.html` for interactive API documentation.


## 🚀 Production Deployment

### Docker Production Build
```bash
# Build production image
docker build -t taskify-backend .

# Run with production environment
docker run -d \
  --name taskify-backend \
  -p 8080:8080 \
  -e ENV=production \
  -e GIN_MODE=release \
  -e DB_HOST=your-db-host \
  taskify-backend
```

### Environment Variables for Production
```env
ENV=production
GIN_MODE=release
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
DB_NAME=taskmanager
SERVER_PORT=8080
```

### Getting Help

- Check the [API Documentation](http://localhost:8080/swagger/index.html)
- Review the health check endpoint: `GET /health`

