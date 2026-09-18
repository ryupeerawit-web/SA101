-- ============================================================================
-- 1. TABLES DEFINITIONS
-- ============================================================================

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(120) NOT NULL,
    email VARCHAR(180) UNIQUE NOT NULL,
    phone VARCHAR(30) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'staff', 'member')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Special Areas Table
CREATE TABLE IF NOT EXISTS special_areas (
    id SERIAL PRIMARY KEY,
    code VARCHAR(40) UNIQUE NOT NULL,
    name VARCHAR(120) NOT NULL,
    description TEXT,
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Time Slots Table
CREATE TABLE IF NOT EXISTS time_slots (
    id SERIAL PRIMARY KEY,
    start_time CHAR(5) NOT NULL,
    end_time CHAR(5) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (start_time, end_time),
    CHECK (start_time < end_time)
);

-- Bookings Table
CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    booking_code VARCHAR(24) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    area_id BIGINT NOT NULL REFERENCES special_areas(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    time_slot_id BIGINT NOT NULL REFERENCES time_slots(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    booking_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed' CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed')),
    cooldown_until TIMESTAMPTZ,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Trainers Table
CREATE TABLE IF NOT EXISTS trainers (
    id SERIAL PRIMARY KEY,
    code VARCHAR(30) UNIQUE NOT NULL,
    full_name VARCHAR(120) NOT NULL,
    specialty VARCHAR(160),
    phone VARCHAR(30),
    email VARCHAR(180),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Trainer Appointments Table
CREATE TABLE IF NOT EXISTS trainer_appointments (
    id SERIAL PRIMARY KEY,
    trainer_id BIGINT NOT NULL REFERENCES trainers(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    member_id BIGINT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    appointment_date DATE NOT NULL,
    start_time CHAR(5) NOT NULL,
    end_time CHAR(5) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed' CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed')),
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- 2. INDEXES
-- ============================================================================

CREATE INDEX IF NOT EXISTS users_phone_idx ON users(phone);

CREATE UNIQUE INDEX IF NOT EXISTS bookings_one_active_slot 
    ON bookings(area_id, time_slot_id, booking_date) 
    WHERE status IN ('pending', 'confirmed');

CREATE INDEX IF NOT EXISTS bookings_cooldown_idx 
    ON bookings(user_id, area_id, time_slot_id, booking_date, cooldown_until);

CREATE INDEX IF NOT EXISTS bookings_date_status_idx 
    ON bookings(booking_date, status);

CREATE INDEX IF NOT EXISTS trainer_appointments_date_idx 
    ON trainer_appointments(appointment_date, status);

CREATE UNIQUE INDEX IF NOT EXISTS trainer_appointments_active_slot 
    ON trainer_appointments(trainer_id, appointment_date, start_time, end_time) 
    WHERE status IN ('pending', 'confirmed');

-- ============================================================================
-- 3. FUNCTIONS & TRIGGERS
-- ============================================================================

-- Function: Updated At Timestamp Trigger
CREATE OR REPLACE FUNCTION set_updated_at() 
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN 
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

-- Triggers for Updated At
DROP TRIGGER IF EXISTS users_updated_at ON users;
CREATE TRIGGER users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS areas_updated_at ON special_areas;
CREATE TRIGGER areas_updated_at BEFORE UPDATE ON special_areas FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS slots_updated_at ON time_slots;
CREATE TRIGGER slots_updated_at BEFORE UPDATE ON time_slots FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS bookings_updated_at ON bookings;
CREATE TRIGGER bookings_updated_at BEFORE UPDATE ON bookings FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trainers_updated_at ON trainers;
CREATE TRIGGER trainers_updated_at BEFORE UPDATE ON trainers FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trainer_appointments_updated_at ON trainer_appointments;
CREATE TRIGGER trainer_appointments_updated_at BEFORE UPDATE ON trainer_appointments FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Helper Functions
CREATE OR REPLACE FUNCTION fn_is_slot_available(
    p_area_id BIGINT,
    p_time_slot_id BIGINT,
    p_booking_date DATE
) RETURNS BOOLEAN LANGUAGE SQL STABLE AS $$
    SELECT NOT EXISTS (
        SELECT 1 
        FROM bookings 
        WHERE area_id = p_area_id 
          AND time_slot_id = p_time_slot_id 
          AND booking_date = p_booking_date 
          AND status IN ('pending', 'confirmed')
    );
$$;

CREATE OR REPLACE FUNCTION fn_booking_count(
    p_area_id BIGINT,
    p_booking_date DATE
) RETURNS BIGINT LANGUAGE SQL STABLE AS $$
    SELECT COUNT(*) 
    FROM bookings 
    WHERE area_id = p_area_id 
      AND booking_date = p_booking_date 
      AND status IN ('pending', 'confirmed');
$$;

-- ============================================================================
-- 4. VIEWS
-- ============================================================================

CREATE OR REPLACE VIEW v_user_booking_history AS
SELECT 
    u.id AS user_id,
    u.full_name,
    u.email,
    u.phone,
    b.id AS booking_id,
    b.booking_code,
    b.booking_date,
    b.status,
    b.note,
    a.id AS area_id,
    a.code AS area_code,
    a.name AS area_name,
    ts.id AS time_slot_id,
    ts.start_time,
    ts.end_time
FROM users u
LEFT JOIN bookings b ON b.user_id = u.id
LEFT JOIN special_areas a ON a.id = b.area_id
LEFT JOIN time_slots ts ON ts.id = b.time_slot_id;

-- ============================================================================
-- 5. INITIAL DATA SEEDS
-- ============================================================================

INSERT INTO special_areas (code, name, description, capacity) 
VALUES
    ('YOGA', 'ห้องโยคะ', 'พื้นที่สำหรับคลาสโยคะและการยืดเหยียด', 15),
    ('SAUNA', 'ห้องอบซาวน่า', 'ห้องอบซาวน่าสำหรับสมาชิก', 8),
    ('WEIGHT', 'โซนเวท', 'โซนฝึกเวทเทรนนิ่ง', 6),
    ('VIP-POOL', 'สระว่ายน้ำ VIP', 'สระว่ายน้ำสำหรับสมาชิก VIP', 10)
ON CONFLICT (code) DO NOTHING;

INSERT INTO time_slots (start_time, end_time) 
VALUES
    ('08:00', '09:00'),
    ('09:00', '10:00'),
    ('10:00', '11:00'),
    ('11:00', '12:00'),
    ('13:00', '14:00'),
    ('14:00', '15:00'),
    ('15:00', '16:00'),
    ('16:00', '17:00'),
    ('17:00', '18:00'),
    ('18:00', '19:00'),
    ('19:00', '20:00'),
    ('20:00', '21:00')
ON CONFLICT (start_time, end_time) DO NOTHING;

INSERT INTO trainers (code, full_name, specialty, phone, email) 
VALUES
    ('T001', 'โค้ชกรณ์ วงศ์สุข', 'เวทเทรนนิ่ง / เพิ่มความฟิต', '0810000001', 'coach.korn@example.com'),
    ('T002', 'โค้ชณัฐชา พรพิพัฒน์', 'คาร์ดิโอ / ลดน้ำหนัก', '0810000002', 'coach.nat@example.com'),
    ('T003', 'โค้ชวิรวัฒน์ ใจเย็น', 'ฟิตเนสส่วนบุคคล', '0810000003', 'coach.wirot@example.com')
ON CONFLICT (code) DO NOTHING;