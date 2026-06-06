# 📊 Sales Module — SalesHub

A full-stack Sales Module system with a Go (Gin) REST API backend and Vue 3 admin dashboard frontend, backed by PostgreSQL.

## 🏗 Architecture

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│  Vue 3 SPA  │────▶│  Go (Gin)    │────▶│  PostgreSQL  │
│  Port: 5173 │     │  Port: 8080  │     │  Port: 5432  │
│  Dashboard  │     │  REST API    │     │  salesdb     │
└─────────────┘     └──────┬───────┘     └──────────────┘
                           │
                    ┌──────┴──────┐
                    │ External    │
                    │ Modules     │
                    │ CRM/Inv/Fin │
                    └─────────────┘
```

**Pattern**: MVC (Model-View-Controller) with Service Layer

## 🛠 Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Vue 3, Vite, Pinia, Vue Router, Chart.js |
| Backend | Go 1.22+, Gin, JWT, lib/pq |
| Database | PostgreSQL 16+ |
| DevOps | Docker, Docker Compose, Nginx |

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- Node.js 18+
- PostgreSQL 16+ (with `salesdb` database created)

### Local Development

**1. Database Setup**
```bash
createdb salesdb
psql -U postgres -d salesdb -f backend/migrations/init.sql
```

**2. Backend**
```bash
cd backend
go mod download
go run cmd/server/main.go
# API running at http://localhost:8080
```

**3. Frontend**
```bash
cd frontend
npm install
npm run dev
# Dashboard at http://localhost:5173
```

### Docker (Production)
```bash
docker-compose up --build
# Frontend: http://localhost:80
# Backend API: http://localhost:8080
```

## 🔐 Authentication

All API endpoints (except `/api/health` and `/api/auth/login`) require JWT authentication.

**Default Credentials:**
- Username: `admin`
- Password: `admin123`

**Usage:**
```bash
# Get token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r .token)

# Use token
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/sales
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/auth/login` | Login & get JWT token |
| `GET` | `/api/health` | Health check |
| `GET` | `/api/sales` | List sales (paginated) |
| `POST` | `/api/sales` | Create new sale |
| `GET` | `/api/sales/:id` | Get sale detail with items |
| `PUT` | `/api/sales/:id` | Update sale |
| `PATCH` | `/api/sales/:id/status` | Update sale status |
| `DELETE` | `/api/sales/:id` | Delete pending sale |
| `GET` | `/api/invoices` | List invoices |
| `GET` | `/api/invoices/:id` | Get invoice detail |
| `POST` | `/api/invoices/generate/:saleId` | Generate invoice |
| `GET` | `/api/customers` | List customers |
| `GET` | `/api/customers/:id/history` | Customer sales history |
| `GET` | `/api/reports/dashboard` | Dashboard stats |
| `GET` | `/api/reports/summary` | Sales summary |
| `GET` | `/api/reports/top-products` | Top selling products |
| `GET` | `/api/reports/revenue` | Revenue report |
| `POST` | `/api/reports/export` | Export to MIS |

### Query Parameters

**Sales & Invoices:**
- `page` — Page number (default: 1)
- `limit` — Items per page (default: 20, max: 100)
- `status` — Filter by status
- `customer_id` — Filter by customer
- `date_from`, `date_to` — Date range filter

**Customer History:**
- `status`, `date_from`, `date_to`, `product_id`

**Reports Summary:**
- `period` — `daily`, `monthly`, or `yearly`

## 🔗 Integration Contracts

### Shared IDs
| ID | Format | Source |
|----|--------|--------|
| `CustomerID` | Integer (auto-increment) | CRM Module |
| `OrderID` | `ORD-YYYY-NNNN` | Sales Module |
| `InvoiceNumber` | `INV-YYYYMMDD-NNNN` | Sales Module |
| `ProductID` | Integer (auto-increment) | Inventory Module |

### External Module Endpoints (Stubs)
- **Inventory**: `http://localhost:8081/api/inventory/check`
- **Finance**: `http://localhost:8082/api/revenue`
- **CRM**: `http://localhost:8083/api/customers/sync`

Integration calls include retry logic with exponential backoff (3 retries).

## 🌐 Ports

| Service | Port |
|---------|------|
| Frontend (dev) | 5173 |
| Frontend (Docker) | 80 |
| Backend API | 8080 |
| PostgreSQL | 5432 |

## 📁 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | localhost | Database host |
| `DB_PORT` | 5432 | Database port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | (empty) | Database password |
| `DB_NAME` | salesdb | Database name |
| `DB_SSLMODE` | disable | SSL mode |
| `JWT_SECRET` | sales-module-secret-key-2026 | JWT signing key |
| `PORT` | 8080 | API server port |

## 📊 Features

- ✅ Sales CRUD with order processing
- ✅ Status tracking (pending → completed → cancelled)
- ✅ Discount calculation with threshold validation (max 50%)
- ✅ Auto-invoice generation on sale completion
- ✅ Customer sales history with filters
- ✅ Daily/monthly/yearly sales summaries
- ✅ Revenue analytics with Chart.js visualizations
- ✅ Top-selling products report
- ✅ JWT authentication on all endpoints
- ✅ CORS middleware
- ✅ MIS/Analytics data export
- ✅ Retry logic for external module integration
- ✅ Docker + Docker Compose deployment
- ✅ Premium dark-theme admin dashboard
