# Special Area Booking Frontend

React + Vite frontend connected to the Go API. Pages now use real URLs instead of one stateful root page:

- `/` main home page for the whole system
- `/special-areas` special-area booking dashboard
- `/special-areas/history` member booking history and detail
- `/reports`, `/staff`, `/members`, and other module routes

The flow is two-step. Clicking an available slot opens confirmation; only `ยืนยันการจอง` calls `POST /bookings` and writes to PostgreSQL. History supports detail view and cancellation. After cancellation, the backend blocks the same account from booking the same area, date, and time slot for 5 minutes.

## Run

```bash
npm install
cp .env.example .env
npm run dev
```

Open `http://localhost:5173/` for the main page, or `http://localhost:5173/special-areas` for booking.

Trainer module routes:

- `/trainers`: trainer directory, search, availability summary, and create an appointment.
- `/trainers/schedule`: combined weekly trainer appointment schedule.

The frontend calls `GET /trainers`, `GET /trainer-appointments`, and `POST /trainer-appointments` from the backend.

Layout fix: the left sidebar is fixed to the viewport height, the user card stays at the bottom, and only the sidebar menu/main content scroll when needed.

Trainer pairing update: admins/staff must search and select a member first, then choose a trainer and confirm an appointment. The create request goes to `POST /api/v1/admin/trainer-appointments` with `member_id` and `trainer_id`; the sidebar user card remains fixed at the bottom.
