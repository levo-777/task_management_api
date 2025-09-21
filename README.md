# Taskify - Task Management System

A comprehensive full-stack task management application with modern authentication, role-based access control, and advanced features.

## 🚀 Quick Overview

Taskify is a complete task management solution that provides:

- **Full-Stack Application**: React frontend + Go backend + PostgreSQL database
- **User Authentication**: Secure JWT-based authentication with refresh tokens
- **Role-Based Access Control**: User and Admin roles with granular permissions
- **Advanced Features**: Pagination, filtering, sorting, caching, and rate limiting
- **Production Ready**: Docker containerization and comprehensive documentation
- **High Performance**: Optimized database queries with caching and indexing

## 🛠 Tech Stack & Versions

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin v1.9.1
- **Database**: PostgreSQL 15.4
- **ORM**: GORM v1.25.5
- **Authentication**: JWT with refresh tokens
- **UUID**: github.com/gofrs/uuid v4.4.0
- **Password Hashing**: golang.org/x/crypto/bcrypt

### Frontend
- **Framework**: React 18.2.0
- **Build Tool**: Create React App 5.0.1
- **HTTP Client**: Axios 1.6.0
- **Routing**: React Router DOM 6.20.1
- **Form Handling**: React Hook Form 7.47.0
- **UI Components**: Material-UI (MUI) 5.14.18
- **JWT Decoding**: jwt-decode 4.0.0
- **Date Handling**: Day.js 1.11.10
- **Web Server**: Nginx 1.25.3 (Alpine)

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Database**: PostgreSQL with optimized indexes
- **Cache**: Redis (optional)
- **Reverse Proxy**: Nginx
- **Monitoring**: Health checks and logging

## 📋 Prerequisites

### For Local Development
- **Go** 1.21 or higher
- **Node.js** 18.0 or higher
- **PostgreSQL** 15 or higher
- **Git**

### For Docker Deployment
- **Docker** 20.10+
- **Docker Compose** 2.0+

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

#### Quick Start
```bash
# Clone the repository
git clone https://github.com/your-username/task_management_api.git
cd task_management_api

# Start all services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f
```

#### Docker Services Overview
The `docker-compose.yml` includes:
- **PostgreSQL Database**: Port 5433 (external), 5432 (internal)
- **Backend API**: Port 8080
- **Frontend**: Port 3000
- **Nginx**: Disabled by default (optional)

#### Docker Management Commands
```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Restart services
docker-compose restart

# View logs
docker-compose logs [service-name]

# Rebuild and restart
docker-compose build --no-cache
docker-compose up -d

# Clean up everything (including volumes)
docker-compose down -v --remove-orphans
docker system prune -f
```

#### Service Endpoints
- ✅ **Frontend**: `http://localhost:3000`
- ✅ **Backend API**: `http://localhost:8080`
- ✅ **API Health**: `http://localhost:8080/health`
- ✅ **Database**: `localhost:5433` (external access)

### Option 2: Local Development (Without Docker)

#### Prerequisites Setup
```bash
# Install PostgreSQL (Ubuntu/Debian)
sudo apt update
sudo apt install postgresql postgresql-contrib

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database and user
sudo -u postgres psql
CREATE DATABASE taskmanager;
CREATE USER taskuser WITH PASSWORD 'taskpass';
GRANT ALL PRIVILEGES ON DATABASE taskmanager TO taskuser;
\q
```

#### Backend Setup
```bash
# Navigate to backend directory
cd backend

# Install Go dependencies
go mod download

# Create environment file (optional)
cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=taskuser
DB_PASSWORD=taskpass
DB_NAME=taskmanager
DB_SSLMODE=disable
SERVER_PORT=8080
JWT_SECRET=your-secret-key-here
EOF

# Run database migrations (automatic on startup)
go run main.go
```

**Backend will be available at**: `http://localhost:8080`

#### Frontend Setup
```bash
# Navigate to frontend directory
cd frontend

# Install Node.js dependencies
npm install

# Create environment file (optional)
cat > .env << EOF
REACT_APP_API_URL=http://localhost:8080/api/v1
EOF

# Start development server
npm start
```

**Frontend will be available at**: `http://localhost:3000`

#### Access Points
- ✅ **Frontend**: `http://localhost:3000`
- ✅ **Backend API**: `http://localhost:8080`
- ✅ **API Health**: `http://localhost:8080/health`
- ✅ **Database**: `localhost:5432`

## 🎯 Features

### Core Functionality
- **Task Management**: Create, read, update, delete tasks
- **User Authentication**: Registration, login, password hashing
- **Role-Based Access**: User and Admin roles with different permissions
- **Task Filtering**: Filter by status, priority, date
- **Task Sorting**: Sort by creation date, priority, status
- **Pagination**: Handle large datasets efficiently

### Advanced Features
- **Real-time Updates**: Live task status updates
- **Search Functionality**: Find tasks by title or description
- **Priority Levels**: Low, medium, high priority tasks
- **Status Tracking**: Pending, in progress, completed
- **User Profiles**: Manage user account settings
- **Admin Dashboard**: User management (admin only)

### Performance & Security
- **High-Performance Caching**: Ristretto cache for frequently accessed data
- **Database Indexing**: Optimized queries with strategic indexes
- **Rate Limiting**: Protection against abuse and brute force attacks
- **SQL Injection Protection**: Parameterized queries
- **JWT Security**: Secure token management with refresh mechanism
- **CORS Configuration**: Secure cross-origin requests

## 📖 API Documentation

The backend provides comprehensive API documentation via Swagger:

- **Interactive Documentation**: `http://localhost:8080/swagger/index.html`
- **Health Check**: `http://localhost:8080/health`
- **API Base URL**: `http://localhost:8080/api/v1`

### Key Endpoints

#### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token

#### Tasks (Protected)
- `GET /api/v1/tasks` - Get all tasks (with pagination/filtering)
- `POST /api/v1/tasks` - Create new task
- `GET /api/v1/tasks/:id` - Get task by ID
- `PUT /api/v1/tasks/:id` - Update task
- `DELETE /api/v1/tasks/:id` - Delete task

#### Users (Protected)
- `GET /api/v1/users/profile` - Get current user profile
- `GET /api/v1/users` - Get all users (admin only)
- `DELETE /api/v1/users/:user_id` - Delete user (admin only)

## 🔧 Configuration

### Environment Variables

#### Backend (.env)
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskmanager
DB_SSLMODE=disable
SERVER_PORT=8080
```

#### Frontend (.env)
```env
REACT_APP_API_URL=http://localhost:8080
REACT_APP_API_VERSION=v1
```

## 🧪 Testing

### Backend Tests
```bash
cd backend

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests
go test ./test/integrations/...
```

### Frontend Tests
```bash
cd frontend

# Run tests
npm test

# Run with coverage
npm test -- --coverage
```

## 🚀 Production Deployment

### Docker Production
```bash
# Build and deploy with Docker Compose
docker-compose -f docker-compose.prod.yml up -d

# Or build individual services
docker build -t taskify-backend ./backend
docker build -t taskify-frontend ./frontend
```

### Environment Variables for Production
```env
# Backend
ENV=production
GIN_MODE=release
DB_HOST=your-db-host
DB_PASSWORD=your-secure-password

# Frontend
REACT_APP_API_URL=https://your-api-domain.com
```

## 📊 Monitoring & Health Checks

### Health Endpoints
- **Backend Health**: `GET /health`
- **Database Status**: Automatic health checks
- **Service Status**: Docker Compose health checks

### Monitoring Features
- Request/response logging
- Error tracking and reporting
- Rate limiting metrics
- Cache hit/miss statistics
- Performance profiling

## 🏗️ Project Structure

```
task_management_api/
├── backend/                 # Go backend API
│   ├── internal/
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── middleware/     # Custom middleware
│   │   ├── models/         # Data models
│   │   ├── services/       # Business logic
│   │   └── utils/          # Utility functions
│   ├── migrations/         # Database migrations
│   ├── test/              # Test files
│   ├── docs/              # Swagger documentation
│   ├── Dockerfile         # Backend Docker config
│   └── README.md          # Backend documentation
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── pages/         # Page components
│   │   ├── services/      # API services
│   │   └── utils/         # Utility functions
│   ├── public/            # Static assets
│   ├── Dockerfile         # Frontend Docker config
│   └── README.md          # Frontend documentation
├── docker-compose.yml      # Full stack orchestration
└── README.md              # This file
```

## 🔒 Security Features

### Authentication & Authorization
- **JWT Tokens**: Secure access tokens with 1-hour expiry
- **Refresh Tokens**: Automatic token renewal
- **Password Hashing**: bcrypt with salt rounds
- **Role-Based Access**: Granular permission system

### API Security
- **Rate Limiting**: 10 req/min general, 3 req/min auth
- **Input Validation**: Request validation and sanitization
- **SQL Injection Protection**: Parameterized queries
- **CORS Configuration**: Secure cross-origin requests

### Infrastructure Security
- **Container Security**: Non-root user in containers
- **Network Security**: Isolated Docker networks
- **Security Headers**: Nginx security headers
- **Health Checks**: Service monitoring

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

**Services Won't Start**
- Check Docker is running: `docker --version`
- Check ports are available: `netstat -tulpn | grep :8080`
- Check logs: `docker-compose logs`

**Database Connection Issues**
- Verify PostgreSQL is running: `sudo systemctl status postgresql`
- Check database exists: `psql -U postgres -l`
- Verify connection string in environment variables

**API Connection Failed**
- Check backend is running: `curl http://localhost:8080/health`
- Verify CORS configuration
- Check browser console for errors

**Authentication Issues**
- Clear browser localStorage
- Check token expiry in JWT payload
- Verify login credentials

### Getting Help

- Check service logs: `docker-compose logs [service-name]`
- Review API documentation: `http://localhost:8080/swagger/index.html`
- Check health endpoints: `http://localhost:8080/health`
- Create an issue in the repository

## 🎉 Default Credentials

The system comes with a default admin user:
- **Username**: `admin`
- **Email**: `admin@gmail.com`
- **Password**: `admin123`

---

**Happy Task Managing! 🚀**
