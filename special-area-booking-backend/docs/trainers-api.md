# Trainer API Specification & Technical Contract

**Base URL:** `/api/v1`  
**Authentication:** HTTP Headers — `Authorization: Bearer <token>`

---

## 1. Public / Authenticated Endpoints

| Method | Endpoint | Access | Description | Query / Body Params |
| :--- | :--- | :--- | :--- | :--- |
| **GET** | `/trainers` | Public / Member | Search and list trainers | **Query:** `search={query}` |
| **GET** | `/trainer-appointments` | Member / Staff | List appointments within date range | **Query:** `from={YYYY-MM-DD}&to={YYYY-MM-DD}&trainer_id={id}` |
| **PATCH** | `/trainer-appointments/:id/cancel` | Member / Staff | Cancel an appointment | **Path:** `id` |

---

## 2. Admin & Staff Endpoints *(Admin/Staff Auth Required)*

| Method | Endpoint | Description | Query / Body Params |
| :--- | :--- | :--- | :--- |
| **GET** | `/admin/members` | Search members | **Query:** `search={query}` |
| **POST** | `/admin/trainers` | Create new trainer profile | **Body:** Trainer details |
| **PUT** | `/admin/trainers/:id` | Update trainer profile | **Path:** `id`, **Body:** Trainer details |
| **POST** | `/admin/trainer-appointments` | Book trainer appointment | **Body:** `{ "trainer_id", "appointment_date", "start_time", "end_time", "note" }` |
| **PATCH** | `/admin/trainer-appointments/:id/status` | Update appointment status | **Path:** `id` |

---

## 3. Access Control & System Logic

* **Visibility Scope:**
  * **Admin & Staff:** Access the combined schedule across all trainers and members.
  * **Members:** Scoped to see only their own appointment history and details.
* **Concurrency & Overlap Control:**
  * Active trainer time slots are strictly protected against double-booking using a PostgreSQL partial unique index on non-canceled slots.
