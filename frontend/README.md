# Taskify Frontend

A modern React-based task management application with a clean, intuitive user interface.

## 🚀 Quick Overview

This frontend application provides a comprehensive user interface for the Taskify task management system:

- **Task Management**: Create, view, update, and delete tasks
- **User Authentication**: Secure login and registration
- **Role-Based UI**: Different interfaces for users and admins
- **Responsive Design**: Works on desktop, tablet, and mobile

## 🛠 Tech Stack

- **Frontend**: React 18
- **Build Tool**: Create React App
- **Styling**: CSS3 with modern features
- **HTTP Client**: Axios
- **State Management**: React Hooks
- **Authentication**: JWT token management
- **Containerization**: Docker
- **Web Server**: Nginx

## 📋 Prerequisites

### For Local Development
- **Node.js** 16.0 or higher
- **npm** 8.0 or higher
- Backend API running on `http://localhost:8080`

### For Docker Deployment
- **Docker** 20.10+
- **Docker Compose** 2.0+

## 🚀 Getting Started

Choose your preferred setup method:

- [**Local Development**](#-local-development-setup) - Run directly on your machine
- [**Docker Development**](#-docker-development-setup) - Run with Docker

---

## 💻 Local Development Setup

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/task_management_api.git
cd task_management_api/frontend
```

### 2. Install Dependencies
```bash
npm install
```

### 3. Environment Configuration
```bash
# Create environment file
cp .env.example .env

# Edit environment variables
nano .env
```

**Environment variables:**
```env
REACT_APP_API_URL=http://localhost:8080
REACT_APP_API_VERSION=v1
```

### 4. Start the Development Server
```bash
npm start
```

✅ **Frontend available at:** `http://localhost:3000`

---

## 🐳 Docker Development Setup

### Quick Start

```bash
# From the project root directory
docker-compose up -d frontend

# Or build and run individually
docker build -t taskify-frontend .
docker run -d -p 3000:3000 taskify-frontend
```

✅ **Frontend available at:** `http://localhost:3000`

## 📖 Available Scripts

### `npm start`
Runs the app in development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in your browser.

The page will reload when you make changes.\
You may also see any lint errors in the console.

### `npm test`
Launches the test runner in interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `npm run build`
Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

### `npm run eject`
**Note: this is a one-way operation. Once you `eject`, you can't go back!**

If you aren't satisfied with the build tool and configuration choices, you can `eject` at any time.


## 🔧 Configuration

### API Integration
The frontend communicates with the backend API through:
- **Base URL**: `http://localhost:8080/api/v1`
- **Authentication**: Bearer token in Authorization header
- **Error Handling**: Comprehensive error handling and user feedback

### Environment Variables
```env
# API Configuration
REACT_APP_API_URL=http://localhost:8080
REACT_APP_API_VERSION=v1

# Feature Flags
REACT_APP_ENABLE_ANALYTICS=false
REACT_APP_DEBUG_MODE=false
```

## 🚀 Production Deployment

### Docker Production Build
```bash
# Build production image
docker build -t taskify-frontend .

# Run with production environment
docker run -d \
  --name taskify-frontend \
  -p 3000:3000 \
  -e REACT_APP_API_URL=https://your-api-domain.com \
  taskify-frontend
```

### Nginx Configuration
The Docker image includes Nginx configuration for:
- Static file serving
- Gzip compression
- Security headers
- SPA routing support

## 🧪 Testing

### Run Tests
```bash
npm test
```

### Run Tests with Coverage
```bash
npm test -- --coverage
```

### Test Types
- **Unit Tests**: Component and utility function tests
- **Integration Tests**: API integration and user flow tests
- **E2E Tests**: End-to-end user journey tests


