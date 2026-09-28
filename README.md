# 🔗 Go URL Shortener – Fast & Simple

[![Go](https://img.shields.io/badge/Go-1.22-blue?logo=go)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-Framework-green?logo=go)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-DB-31648c?logo=postgresql)](https://www.postgresql.org/)
[![Clerk](https://img.shields.io/badge/Auth-Clerk-orange)](https://clerk.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

✨ A backend URL Shortener built using **Golang + Gin + PostgreSQL**, with **Clerk authentication**, rate limiting and clean routes.

---

## ⚡ Key Features

- 🔗 URL Shortening
- 🔁 Redirect Handler
- 🔐 Clerk Auth Middleware
- 📦 PostgreSQL Storage
- ⚡ Gin high-performance router
- 🧩 Signup webhook
- 🧱 Rate Limiting middleware
- 🔧 Modular structure

---

## 📦 Tech Stack

- 🦦 **Golang (1.22+)** – Fast, compiled, highly performant backend language
- 🌿 **Gin Framework** – Lightweight, minimal & extremely fast HTTP router
- 🗄️ **PostgreSQL** – Relational storage for persistent URL and user data
- 🛡️ **Clerk Authentication** – Secure user auth + webhooks

---

## 🔗 API Endpoints

Below are all available endpoints explained briefly 👇

- **GET `/health-check`** _(Public)_

  Checks if the server is healthy and running properly.

---

- **POST `/shorten`** _(Private)_

  Creates a short URL for the authenticated user.  
  Requires Clerk authentication & rate limiting.

---

- **GET `/:code`** _(Public)_

  Redirects to the original long URL based on the short code.  
  Returns `404` if the URL does not exist.

---

- **GET `/stats`** _(Private)_

  Returns paginated list of URLs owned by the authenticated user.  
  Useful for analytics / dashboard views.

---

- **DELETE `/:code`** _(Private)_

  Deletes a specific URL owned by the authenticated user.

---

- **POST `/webhook/signup`** _(Public)_

  Clerk signup webhook.  
  Triggers when a new user registers.

---

## 🛠️ Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/LokeshXs/urlbit_url-shortener.git
cd urlbit_url-shortener
```

### 2. Install dependencies

```bash
go mod tidy

```

### 3. Setup Environment variables in .env file from .env.example file

```bash
cp .env.example .env
```

### Then open .env and fill the values:

```bash
POSTGRES_URL=<your-postgres-url>
CLERK_SECRET_KEY=<your-clerk-secret-key>

```

### 6. Run the development server

```bash
go run main.go

```

If everything is setup correctly, you should now see the service running locally 🎉

Database connection or table initialization failures are logged without stopping
the HTTP server. `/health-check` remains available and reports server liveness.
Database-dependent endpoints return HTTP 500 with
`{"message":"Internal server error"}` while PostgreSQL is unavailable. Each
request retries the database check and any unfinished table initialization, so
service resumes when PostgreSQL recovers. Missing or invalid `POSTGRES_URL`
configuration requires correcting the environment and restarting the app.

---

✨ **Why I built this**

Because almost every product needs a clean and reliable short-link system—and building this from scratch gives maximum control & full ownership instead of depending on third-party platforms 👨‍💻

---

## 🧑‍💻 Author

**Lokesh Singh**  
🔗 Portfolio — https://lokesh-singh.vercel.app/
