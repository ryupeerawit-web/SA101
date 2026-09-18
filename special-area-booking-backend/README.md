# Special Area Booking Backend

Go + Gin + GORM + PostgreSQL backend for the supplied special-area booking UI. Includes authentication, admin/staff/member roles, users with name/email/phone, booking history through user_id, areas, time slots, availability, booking lifecycle, five-minute cancellation cooldown, dashboard summary, foreign keys, duplicate-booking protection, SQL functions, Docker, and pgAdmin.

## Start

```bash
cp .env.example .env
docker compose up --build -d
docker compose logs -f api
```

The API does not run GORM AutoMigrate on startup. PostgreSQL schema creation is owned by `db/migrations/001_init.sql`; the API only applies a safe legacy repair for the `users.phone` column. This avoids the PostgreSQL `uni_users_email` constraint loop.

API: http://localhost:8080/api/v1/health
pgAdmin: http://localhost:5050, login `admin@example.com` / `admin`

In pgAdmin register PostgreSQL with host `postgres`, port `5432`, database `special_booking`, user `booking`, password `booking`. The API seeds an admin from `.env`, then seeds the four areas and twelve one-hour slots shown in the UI.

Change the default admin password, admin phone, and JWT secret before production. Use HTTPS and do not expose pgAdmin to the public internet.

Trainer module: the migration seeds three trainers and exposes `/trainers` plus `/trainer-appointments`; admin/staff can manage trainers and the combined schedule.
