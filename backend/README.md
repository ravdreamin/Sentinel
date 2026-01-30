# Sentinel

**Sentinel** is a high-performance, concurrent web scraper and metadata extraction engine built with **Go**.  
It is designed to handle **high-volume URL processing** using a configurable worker pool, triggered via **file uploads or manual input**.

This project focuses on **concurrency, I/O efficiency, authentication flows, and backend reliability** — not toy CRUD APIs.

---

## 🚀 Features

### Concurrent Processing
- Configurable **worker pool** (default: 100 workers)
- Parallel URL scraping with controlled concurrency

### Multi-format URL Ingestion
- Extracts URLs from:
  - `.txt`
  - `.csv`
  - `.json`
  - `.pdf`

### Authentication & Access Control
- JWT-based authentication
- Email registration with **OTP verification** (Brevo / Sendinblue)
- Google OAuth2 integration
- Supports multiple auth providers per user

### Metadata Extraction
For each URL:
- HTTP status code
- Response time
- HTML content hash (SHA-256)
- Page title
- `<h1>` tags
- Meta description
- Outbound link extraction

### Persistence & Job Tracking
- PostgreSQL-backed storage
- Full job lifecycle tracking
- User management and result storage
- JSONB-based metadata persistence

---

## 🛠 Tech Stack

- **Language:** Go 1.25+
- **Web Framework:** Gin Gonic
- **HTML Parsing:** Goquery
- **Database:** PostgreSQL 15
- **Auth:** JWT, Google OAuth2
- **Email:** Brevo (SMTP)
- **DB Driver:** pgx
- **Containerization:** Docker, Docker Compose

---

## 📂 Project Structure

├── cmd/api/ # Application entry point
│ └── main.go
├── internal/
│ ├── database/ # DB connection & CRUD (User, Job, Result)
│ ├── email/ # Brevo SMTP client
│ ├── models/ # Data models & mappings
│ ├── server/ # HTTP handlers & auth middleware
│ ├── utils/ # JWT & OTP utilities
│ └── worker/ # Worker pool & scraping logic
├── migrations/ # SQL schema migrations
├── uploads/ # Temporary storage for uploaded files
└── docker-compose.yml


---

## 🏁 Getting Started

### 1. Prerequisites
- Docker & Docker Compose
- Go 1.25+

---

### 2. Environment Configuration

Create a `.env` file in the project root:

```env
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
