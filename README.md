# CulTour Backend

## Project Overview

CulTour is a robust backend application designed to provide comprehensive cultural and tourism services. It offers APIs for managing events, local stories, discussions, and user interactions, primarily focusing on preserving and promoting Indonesian cultural heritage.

## Technology Stack

-   **Language**: Go (Golang) 1.20+
-   **Web Framework**: Gin
-   **Database**: PostgreSQL
-   **Authentication**: Supabase
-   **AI Integration**: Google Generative AI
-   **API Documentation**: Swagger

## Prerequisites

Ensure the following are installed on your system:

-   [Go (Golang)](https://go.dev/doc/install) version 1.20+.
-   [PostgreSQL](https://www.postgresql.org/download/) database system.
-   [Git](https://git-scm.com/downloads) for repository cloning.

## Installation and Local Setup

Follow these steps to set up and run the CulTour Backend locally.

### 1. Clone the Repository

Open your terminal or command prompt and execute:

```bash
git clone https://github.com/holycann/cultour-backend.git
cd cultour-backend
```

### 2. Configure Environment Variables

Create a `.env` file in the project root and populate it with the following. These variables are crucial for database, Supabase, and AI service connections.

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_NAME=cultour_database

# Supabase Configuration
SUPABASE_URL=your_supabase_project_url
SUPABASE_KEY=your_supabase_project_key
SUPABASE_PROJECT_ID=your_supabase_project_id

# Gemini AI Configuration
GEMINI_API_KEY=your_gemini_api_key
GEMINI_AI_MODEL=gemini-pro
GEMINI_TEMPERATURE=0.7
GEMINI_TOP_K=40
GEMINI_TOP_P=0.95

# Application Settings
APP_ENV=development # Options: 'development', 'production'
APP_PORT=8181       # Port for the application
SERVER_HOST=0.0.0.0 # Host address (0.0.0.0 for all interfaces)
PRODUCTION_DOMAIN=your.production.domain # Domain for Swagger in production
LOG_LEVEL=info      # Logging level
LOG_FILE_PATH=logs/app.log # Path for logs
```

### 3. Install Dependencies

Navigate to the project root and install Go modules:

```bash
go mod tidy
```

### 4. Database Preparation

Ensure your Supabase database is configured. The database will be automatically set up through Supabase configuration:

1. Verify your Supabase project is initialized
2. Ensure the connection details in the `.env` file are correct
3. Supabase will handle database schema and migrations

### 5. Run the Application

#### Development Mode

```bash
go run cmd/main.go
```

The server will be accessible at `http://localhost:8181`.

#### Production Build

```bash
# Build the executable
go build -o cultour-backend ./cmd

# Run the compiled application
./cultour-backend
```

## Deployment Guide

This guide covers deploying the CulTour Backend to a production environment, specifically on AWS EC2.

### AWS EC2 Deployment

When deploying to an AWS EC2 instance, consider these best practices:

1.  **Environment Variables**: Securely configure all `.env` variables on your EC2 instance. Set `APP_ENV=production` and `PRODUCTION_DOMAIN` to your live domain.
2.  **Security Groups**: Allow inbound traffic on `APP_PORT` (default 8181).
3.  **Process Management**: Use `systemd` or `Supervisor` for continuous background operation and automatic restarts.
4.  **Reverse Proxy**: Implement Nginx or Apache for SSL/TLS termination, request forwarding, and load balancing.

### Test Deployment Link (AWS EC2 Example)

Access the deployed application (replace with your actual EC2 public IP or domain name):

[http://cultour.holyycan.com/docs/index.html](http://cultour.holyycan.com/docs/index.html)

## Project Structure

```
cultour-backend/
├── cmd/            # Application entry points
├── configs/        # Configuration files
├── internal/       # Core application logic (e.g., cultural, discussion, place, users modules)
├── pkg/            # Shared utilities
└── docs/           # API documentation
```

## API Documentation

Swagger documentation is available at:

-   **Local Development**: `http://localhost:8181/swagger/index.html`
-   **Production**: `http://cultour.holyycan.com/docs/index.html`

## Security Features

-   Supabase authentication
-   JWT based authorization
-   Environment based configuration for sensitive data

## Acknowledgments

-   Developed for **Garuda Hack 6.0**

**CulTour: Bridging Culture, Inspiring Exploration**