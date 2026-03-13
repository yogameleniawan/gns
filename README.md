<p align="center">
  <img src="frontend/public/gns.png" alt="GNS" width="120" />
  <img src="frontend/public/logo.png" alt="GNS" width="120" />
</p>

<h1 align="center">GNS — Go + Next.js + shadcn/ui</h1>

<p align="center">
  A production-ready full-stack starter kit with authentication, role-based access control, i18n, and clean architecture.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Next.js-16-000000?style=flat&logo=nextdotjs&logoColor=white" />
  <img src="https://img.shields.io/badge/Tailwind_CSS-v4-38BDF8?style=flat&logo=tailwindcss&logoColor=white" />
  <img src="https://img.shields.io/badge/shadcn/ui-Components-000?style=flat&logo=shadcnui&logoColor=white" />
  <img src="https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql&logoColor=white" />
  <img src="https://img.shields.io/badge/Redis-7-DC382D?style=flat&logo=redis&logoColor=white" />
</p>

---

## ✨ Features

- **JWT Authentication** — Login, register, Google OAuth, token refresh, session management
- **Role-Based Access Control** — Roles, permissions, and module-level access
- **Clean Architecture** — Consistent structure on both backend and frontend
- **Dependency Injection** — uber/dig for automatic wiring on the backend
- **Internationalization** — English & Indonesian out of the box (next-intl)
- **Theme System** — Light & dark mode with next-themes + shadcn/ui
- **Database Migrations** — Auto-run on startup with golang-migrate
- **Docker Ready** — Development and production Docker Compose configs
- **API Rate Limiting** — Built-in rate limiter middleware
- **Security Headers** — HSTS, X-Frame-Options, Content-Type sniffing protection

## 🏗️ Tech Stack

| Layer        | Technology                                       |
| ------------ | ------------------------------------------------ |
| **Backend**  | Go, Chi Router, uber/dig, PostgreSQL, Redis      |
| **Frontend** | Next.js 16, React 19, Tailwind CSS v4, shadcn/ui |
| **Auth**     | JWT (access + refresh tokens), Google OAuth      |
| **State**    | TanStack Query, Zustand                          |
| **i18n**     | next-intl (EN, ID)                               |
| **DevOps**   | Docker, Docker Compose, Makefile                 |

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+ or Bun
- Docker & Docker Compose

### 1. Clone

```bash
git clone https://github.com/yogameleniawan/gns.git
cd gns
```

### 2. Backend

```bash
cd backend
cp .env.example .env
cp config/config.template.yaml config/config.development.yaml

# Option A: Docker (recommended)
make dev

# Option B: Local (requires PostgreSQL & Redis)
go mod download
make init
make run
```

### 3. Frontend

```bash
cd frontend
bun install    # or: npm install
bun dev        # or: npm run dev
```

### 4. Open

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080

### Default Admin Account

| Field    | Value           |
| -------- | --------------- |
| Email    | `admin@gns.com` |
| Password | `admin123`     |

## 📁 Project Structure

```
gns/
├── backend/
│   ├── cmd/api/          # Entry point (main, config, migration, server)
│   ├── config/           # YAML config per environment
│   ├── container/        # Dependency injection (uber/dig)
│   ├── internal/         # Business modules
│   │   ├── auth/         # Authentication (dto, repo, service, handler)
│   │   └── rbac/         # Roles & permissions
│   ├── migrations/       # SQL migrations & seeders
│   └── pkg/              # Shared packages (middleware, router, utils, etc.)
│
├── frontend/
│   ├── app/              # Next.js App Router pages
│   ├── src/
│   │   ├── domain/       # Types & interfaces
│   │   ├── application/  # Hooks & state management
│   │   ├── infrastructure/ # API clients & stores
│   │   └── presentation/ # Components & pages
│   └── locales/          # i18n translations
```

## 🔧 Backend Commands

```bash
make dev                    # Start dev environment (Docker)
make dev-down               # Stop dev containers
make run                    # Run locally
make build                  # Build binary
make create-migration name=X  # Create migration
make migrate-up             # Run migrations
make migrate-down           # Rollback migrations
make db-shell               # PostgreSQL shell
```

## 📖 Adding a New Module

Each backend module follows a 4-file pattern inside `internal/`:

```
internal/your_module/
├── dto.go          # Request/response structs
├── repository.go   # Database queries
├── service.go      # Business logic
└── handler.go      # HTTP handlers
```

Then wire it in `container/container.go` and add routes in `pkg/router/router.go`.

See the [Documentation page](/docs) for detailed guides.

## 🌐 API Routes

All routes are prefixed with `/v1`.

| Method | Path                | Description                 |
| ------ | ------------------- | --------------------------- |
| POST   | `/auth/register`    | Register                    |
| POST   | `/auth/login`       | Login                       |
| POST   | `/auth/refresh`     | Refresh token               |
| GET    | `/auth/profile`     | Get profile (🔒)            |
| GET    | `/users`            | List users (🔒 Admin)       |
| POST   | `/users`            | Create user (🔒 Admin)      |
| GET    | `/rbac/roles`       | List roles (🔒 Admin)       |
| GET    | `/rbac/permissions` | List permissions (🔒 Admin) |

🔒 = Requires authentication

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
