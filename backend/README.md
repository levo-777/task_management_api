# Taskify Backend API

A comprehensive task management system API with CRUD operations, pagination, filtering, sorting, rate limiting, and caching.

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

## 🔧 Advanced Features

### Pagination
- **Page-based pagination** for efficient large dataset handling
- Use `page` parameter to control page number
- Use `page_size` parameter to control page size

### Filtering
- **By status**: `?status=pending|in_progress|completed`
- **By priority**: `?priority=low|medium|high`
- **By search**: `?search=keyword`

### Sorting
- **Sort by**: `title`, `status`, `priority`, `created_at`, `updated_at`
- **Order**: `asc` or `desc`
- Example: `?sort_by=created_at&sort_order=desc`

### Rate Limiting
- **10 requests per minute** for general endpoints
- **3 requests per minute** for authentication endpoints
- Applied to all API endpoints except health and documentation

### Caching
- **High-performance Ristretto cache** for frequently accessed tasks
- Automatic cache invalidation on updates
- 5-minute TTL for cached tasks

## Authentication

All protected endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <access_token>
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

### Performance Monitoring
The application includes built-in performance monitoring:
- Request/response logging
- Error tracking
- Rate limiting metrics
- Cache hit/miss statistics

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

## 🔒 Security Features

- **Password Hashing**: bcrypt with salt rounds
- **JWT Authentication**: Access tokens with 1-hour expiry
- **Refresh Tokens**: Secure token renewal mechanism
- **Role-Based Access Control**: User and Admin roles with granular permissions
- **Rate Limiting**: Protection against brute force attacks
- **SQL Injection Protection**: Parameterized queries
- **CORS Configuration**: Secure cross-origin requests
- **Input Validation**: Request validation and sanitization

## 🏗️ Project Structure

```
backend/
├── docs/                    # Swagger documentation
├── internal/
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # Custom middleware (auth, rate limiting, logging)
│   ├── models/            # Data models
│   ├── repositories/      # Database connection and config
│   ├── services/          # Business logic layer
│   └── utils/             # Utility functions (pagination, etc.)
├── migrations/            # Database migrations
├── test/                  # Test files
│   └── integrations/      # Integration tests
├── Dockerfile            # Docker configuration
├── go.mod               # Go module dependencies
├── main.go              # Application entry point
└── README.md            # This file
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

**Database Connection Failed**
- Check if PostgreSQL is running
- Verify database credentials in environment variables
- Ensure database exists

**Rate Limit Exceeded**
- Wait 1 minute between requests
- Check rate limit headers in response
- Use authentication for higher limits

**JWT Token Expired**
- Use the refresh token endpoint to get a new access token
- Check token expiry time in JWT payload

### Getting Help

- Check the [API Documentation](http://localhost:8080/swagger/index.html)
- Review the health check endpoint: `GET /health`
- Check application logs for detailed error messages
- Create an issue in the repository for bugs or feature requests
