# Taskify Frontend

A modern React-based task management application with a clean, intuitive user interface.

## 🚀 Quick Overview

This frontend application provides a comprehensive user interface for the Taskify task management system:

- **Task Management**: Create, view, update, and delete tasks
- **User Authentication**: Secure login and registration
- **Role-Based UI**: Different interfaces for users and admins
- **Responsive Design**: Works on desktop, tablet, and mobile
- **Real-time Updates**: Live task status updates
- **Advanced Filtering**: Sort and filter tasks by status, priority, and date

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

## 🎨 Features

### User Interface
- **Dashboard**: Overview of tasks and statistics
- **Task List**: Paginated list with filtering and sorting
- **Task Form**: Create and edit tasks with validation
- **User Profile**: Manage user account settings
- **Admin Panel**: User management (admin only)

### Authentication
- **Login Form**: Secure user authentication
- **Registration**: New user account creation
- **Token Management**: Automatic token refresh
- **Protected Routes**: Route-level authentication

### Task Management
- **Create Tasks**: Add new tasks with title, description, priority
- **Update Tasks**: Modify existing tasks
- **Delete Tasks**: Remove tasks (with confirmation)
- **Status Updates**: Change task status (pending, in progress, completed)
- **Priority Levels**: Set task priority (low, medium, high)

### Advanced Features
- **Search**: Find tasks by title or description
- **Filtering**: Filter by status, priority, or date range
- **Sorting**: Sort by creation date, priority, or status
- **Pagination**: Navigate through large task lists
- **Responsive Design**: Optimized for all screen sizes

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

## 📱 Responsive Design

The application is fully responsive and optimized for:
- **Desktop**: 1200px and above
- **Tablet**: 768px to 1199px
- **Mobile**: 320px to 767px

## 🔒 Security Features

- **JWT Token Management**: Secure token storage and refresh
- **Input Validation**: Client-side validation with server-side verification
- **XSS Protection**: Content sanitization and CSP headers
- **CSRF Protection**: Token-based CSRF protection
- **Secure Headers**: Security headers via Nginx

## 🏗️ Project Structure

```
frontend/
├── public/                 # Static assets
├── src/
│   ├── components/        # Reusable UI components
│   ├── pages/            # Page components
│   ├── hooks/            # Custom React hooks
│   ├── services/         # API service functions
│   ├── utils/            # Utility functions
│   ├── contexts/         # React contexts
│   ├── styles/           # CSS styles
│   └── App.js            # Main application component
├── Dockerfile            # Docker configuration
├── nginx.conf            # Nginx configuration
├── package.json          # Dependencies and scripts
└── README.md             # This file
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

**API Connection Failed**
- Check if backend API is running
- Verify API URL in environment variables
- Check browser console for CORS errors

**Authentication Issues**
- Clear browser localStorage
- Check token expiry
- Verify login credentials

**Build Failures**
- Clear node_modules and package-lock.json
- Run `npm install` again
- Check Node.js version compatibility

### Getting Help

- Check the browser console for error messages
- Review the [Backend API Documentation](http://localhost:8080/swagger/index.html)
- Create an issue in the repository for bugs or feature requests
