Sentinel

Sentinel is a high-performance, concurrent web scraper and metadata extractor built with Go. It features a robust worker pool architecture to handle high-volume scraping tasks triggered by file uploads or manual input.

🚀 Key Features

Concurrent Processing: Configurable worker pool (default 100 workers) for parallel URL processing.

Multi-format Parser: Extract URLs from .txt, .csv, .json, and .pdf files.

Robust Auth System:

JWT-based authentication.

Email registration with OTP verification via Brevo (formerly Sendinblue).

Google OAuth2 integration.

Comprehensive Metadata Extraction:

HTTP Status Codes and Response Times.

HTML content hashing (SHA256).

Title, H1, and Meta Description extraction.

Outbound link harvesting.

PostgreSQL Persistence: Full job tracking, user management, and result storage using pgx and SQL migrations.

🛠 Tech Stack

Language: Go 1.25.5

Web Framework: Gin Gonic

HTML Parsing: Goquery

Database: PostgreSQL (v15)

Auth: JWT & Google OAuth2

Containerization: Docker & Docker Compose

📂 Project Structure

├── cmd/api/             # Entry point (main.go)
├── internal/
│   ├── database/        # DB connection and CRUD (User, Job, Result)
│   ├── email/           # Brevo SMTP client implementation
│   ├── models/          # Data structures and DB mappings
│   ├── server/          # HTTP Handlers and Auth Middleware
│   ├── utils/           # JWT and OTP generation
│   └── worker/          # Concurrent pool and scraping logic
├── migrations/          # SQL schema migrations
└── uploads/             # Temporary storage for processed files


🏁 Getting Started

1. Prerequisites

Docker & Docker Compose

Go 1.25+

2. Environment Configuration

Create a .env file in the root directory:

DB_USER=admin
DB_PASS=password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=sentinel

JWT_SECRET=your_super_secret_key
EMAIL_APIKEY=your_brevo_api_key

GOOGLE_CLIENT_ID=your_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_secret
GOOGLE_REDIRECT_URL=http://localhost:8081/auth/google/callback


3. Spin up Infrastructure

Use Docker Compose to start the database and run migrations:

docker-compose up -d


4. Run the API

go run cmd/api/main.go


The server will start on :8081.

📡 API Reference

Authentication

Endpoint

Method

Description

/register

POST

Create user and send OTP

/verify

POST

Verify email with OTP

/login

POST

Exchange credentials for JWT

/auth/google/login

GET

Start Google OAuth flow

Protected Endpoints (Requires Bearer Token)

Endpoint

Method

Description

/api/profile

GET

Get current user info

/api/upload

POST

Upload file (document field) to scrape

/api/set-password

POST

Set password for OAuth users

🏗 Database Schema

The system uses several key tables:

users: Core account data and verification status.

user_identities: Mapping for multiple auth providers (local, google).

jobs: Tracking file paths, job types, and processing status.

results: JSONB storage for all scraped metadata.

verifications: Temporary storage for OTP codes and expiry.
