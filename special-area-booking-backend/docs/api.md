# API Specification & Technical Contract

**Base URL:** `http://localhost:8080/api/v1`  
**Timezone:** `Asia/Bangkok`  
**Date Format:** `YYYY-MM-DD`  
**Authentication:** HTTP Headers — `Authorization: Bearer <token>`

---

## 1. Public Endpoints

| Method | Endpoint | Description | Query / Body Params |
| :--- | :--- | :--- | :--- |
| **GET** | `/health` | System health check | — |
| **POST** | `/auth/register` | Register new account | **Body:** `{ "full_name", "email", "phone", "password" }` |
| **POST** | `/auth/login` | User authentication | **Body:** `{ "email", "password" }` |
| **GET** | `/areas` | List all areas | — |
| **GET** | `/slots` | List time slots | — |
| **GET** | `/availability` | Check area availability | **Query:** `date={YYYY-MM-DD}&area_id={id}` |

---

## 2. Member Endpoints *(Auth Required)*

| Method | Endpoint | Description | Query / Body Params |
| :--- | :--- | :--- | :--- |
| **GET** | `/me` | Get current user profile | — |
| **POST** | `/bookings` | Create new booking | **Body:** `{ "area_id", "time_slot_id", "booking_date", "note" }` |
| **GET** | `/bookings` | List user's bookings | — |
| **GET** | `/bookings/:id` | Get booking details | **Path:** `id` |
| **PATCH** | `/bookings/:id/cancel` | Cancel booking *(Triggers cooldown)* | **Path:** `id` |

---

## 3. Admin & Staff Endpoints *(Admin Auth Required)*

| Method | Endpoint | Description | Query / Body Params |
| :--- | :--- | :--- | :--- |
| **GET** | `/admin/dashboard/summary` | Get summary statistics | **Query:** `date={YYYY-MM-DD}` |
| **GET** | `/admin/bookings` | List all system bookings | — |
| **PATCH** | `/admin/bookings/:id/status` | Update booking status | **Path:** `id` |
| **POST** | `/admin/areas` | Create new area | **Body:** `{ "name", ... }` |
| **PUT** | `/admin/areas/:id` | Update existing area | **Path:** `id` |
| **POST** | `/admin/slots` | Create new time slot | **Body:** `{ "start_time", "end_time", ... }` |

---

## 4. Business Rules & Logic

* **Cancellation Cooldown:** Calling `PATCH /bookings/:id/cancel` updates the booking status to `cancelled` and returns a `cooldown_until` timestamp. The requesting user is blocked from re-booking the exact same area, date, and time slot for **5 minutes**.
* **Dashboard Data Mapping:**
  * **Summary Cards:** Fetch from `GET /admin/dashboard/summary`
  * **Area Lists / Cards:** Fetch from `GET /areas`
  * **Time Slot Grid:** Fetch from `GET /availability`
* **Database Views:**
  * User profile details originate from the `users` table (`full_name`, `email`, `phone`).
  * Historical records are available via the PostgreSQL view: `v_user_booking_history`.
