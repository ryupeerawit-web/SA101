package database

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ryu111/special-area-booking/internal/config"
	"github.com/ryu111/special-area-booking/internal/models"
)

// Connect initializes and returns a GORM PostgreSQL database connection.
func Connect(c config.Config) (*gorm.DB, error) {
	dsn := c.DatabaseURL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
			c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
		)
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// Migrate executes GORM AutoMigrate and applies database schema objects.
func Migrate(db *gorm.DB) error {
	// 1. Run AutoMigrate using the provided models
	err := db.AutoMigrate(
		&models.User{},
		&models.SpecialArea{},
		&models.TimeSlot{},
		&models.Booking{},
		&models.Trainer{},
		&models.TrainerAppointment{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	// 2. Custom SQL Statements for Constraints, Indexes, Triggers, Functions, and Views
	statements := []string{
		// Unique Constraint for TimeSlots (start_time + end_time)
		`CREATE UNIQUE INDEX IF NOT EXISTS time_slots_start_end_key ON time_slots(start_time, end_time)`,

		// Partial Unique Indexes
		`CREATE UNIQUE INDEX IF NOT EXISTS bookings_one_active_slot 
		 ON bookings(area_id, time_slot_id, booking_date) 
		 WHERE status IN ('pending', 'confirmed')`,

		`CREATE INDEX IF NOT EXISTS bookings_cooldown_idx 
		 ON bookings(user_id, area_id, time_slot_id, booking_date, cooldown_until)`,

		`CREATE UNIQUE INDEX IF NOT EXISTS trainer_appointments_active_slot 
		 ON trainer_appointments(trainer_id, appointment_date, start_time, end_time) 
		 WHERE status IN ('pending', 'confirmed')`,

		// Trigger Function for updated_at
		`CREATE OR REPLACE FUNCTION set_updated_at() 
		 RETURNS TRIGGER LANGUAGE plpgsql AS $$
		 BEGIN 
		     NEW.updated_at = NOW();
		     RETURN NEW;
		 END;
		 $$;`,

		// Create Triggers
		`DROP TRIGGER IF EXISTS users_updated_at ON users;
		 CREATE TRIGGER users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		`DROP TRIGGER IF EXISTS areas_updated_at ON special_areas;
		 CREATE TRIGGER areas_updated_at BEFORE UPDATE ON special_areas FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		`DROP TRIGGER IF EXISTS slots_updated_at ON time_slots;
		 CREATE TRIGGER slots_updated_at BEFORE UPDATE ON time_slots FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		`DROP TRIGGER IF EXISTS bookings_updated_at ON bookings;
		 CREATE TRIGGER bookings_updated_at BEFORE UPDATE ON bookings FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		`DROP TRIGGER IF EXISTS trainers_updated_at ON trainers;
		 CREATE TRIGGER trainers_updated_at BEFORE UPDATE ON trainers FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		`DROP TRIGGER IF EXISTS trainer_appointments_updated_at ON trainer_appointments;
		 CREATE TRIGGER trainer_appointments_updated_at BEFORE UPDATE ON trainer_appointments FOR EACH ROW EXECUTE FUNCTION set_updated_at();`,

		// Helper Functions
		`CREATE OR REPLACE FUNCTION fn_is_slot_available(
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
		 $$;`,

		`CREATE OR REPLACE FUNCTION fn_booking_count(
		     p_area_id BIGINT,
		     p_booking_date DATE
		 ) RETURNS BIGINT LANGUAGE SQL STABLE AS $$
		     SELECT COUNT(*) 
		     FROM bookings 
		     WHERE area_id = p_area_id 
		       AND booking_date = p_booking_date 
		       AND status IN ('pending', 'confirmed');
		 $$;`,

		// Views
		`CREATE OR REPLACE VIEW v_user_booking_history AS
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
		 LEFT JOIN time_slots ts ON ts.id = b.time_slot_id;`,
	}

	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			return fmt.Errorf("failed to execute migration statement: %w", err)
		}
	}

	return nil
}

// Seed inserts initial master data including default Admin, Special Areas, Time Slots, and Trainers.
func Seed(db *gorm.DB, c config.Config) error {
	// 1. Seed initial Admin user
	var count int64
	if err := db.Model(&models.User{}).Where("email = ?", c.AdminEmail).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 && c.AdminEmail != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		adminUser := models.User{
			FullName:     c.AdminName,
			Email:        c.AdminEmail,
			Phone:        c.AdminPhone,
			PasswordHash: string(hashedPassword),
			Role:         models.RoleAdmin,
			Status:       models.UserActive,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return err
		}
		log.Printf("seeded admin: %s", c.AdminEmail)
	}

	// 2. Seed initial Special Areas
	areas := []models.SpecialArea{
		{Code: "YOGA", Name: "ห้องโยคะ", Description: "พื้นที่สำหรับคลาสโยคะและการยืดเหยียด", Capacity: 15, Active: true},
		{Code: "SAUNA", Name: "ห้องอบซาวน่า", Description: "ห้องอบซาวน่าสำหรับสมาชิก", Capacity: 8, Active: true},
		{Code: "WEIGHT", Name: "โซนเวท", Description: "โซนฝึกเวทเทรนนิ่ง", Capacity: 6, Active: true},
		{Code: "VIP-POOL", Name: "สระว่ายน้ำ VIP", Description: "สระว่ายน้ำสำหรับสมาชิก VIP", Capacity: 10, Active: true},
	}
	for _, area := range areas {
		var existing models.SpecialArea
		db.Where("code = ?", area.Code).FirstOrCreate(&existing, area)
	}

	// 3. Seed initial Time Slots
	slots := []models.TimeSlot{
		{StartTime: "08:00", EndTime: "09:00", Active: true},
		{StartTime: "09:00", EndTime: "10:00", Active: true},
		{StartTime: "10:00", EndTime: "11:00", Active: true},
		{StartTime: "11:00", EndTime: "12:00", Active: true},
		{StartTime: "13:00", EndTime: "14:00", Active: true},
		{StartTime: "14:00", EndTime: "15:00", Active: true},
		{StartTime: "15:00", EndTime: "16:00", Active: true},
		{StartTime: "16:00", EndTime: "17:00", Active: true},
		{StartTime: "17:00", EndTime: "18:00", Active: true},
		{StartTime: "18:00", EndTime: "19:00", Active: true},
		{StartTime: "19:00", EndTime: "20:00", Active: true},
		{StartTime: "20:00", EndTime: "21:00", Active: true},
	}
	for _, slot := range slots {
		var existing models.TimeSlot
		db.Where("start_time = ? AND end_time = ?", slot.StartTime, slot.EndTime).FirstOrCreate(&existing, slot)
	}

	// 4. Seed initial Trainers
	trainers := []models.Trainer{
		{Code: "T001", FullName: "โค้ชกรณ์ วงศ์สุข", Specialty: "เวทเทรนนิ่ง / เพิ่มความฟิต", Phone: "0810000001", Email: "coach.korn@example.com", Active: true},
		{Code: "T002", FullName: "โค้ชณัฐชา พรพิพัฒน์", Specialty: "คาร์ดิโอ / ลดน้ำหนัก", Phone: "0810000002", Email: "coach.nat@example.com", Active: true},
		{Code: "T003", FullName: "โค้ชวิรวัฒน์ ใจเย็น", Specialty: "ฟิตเนสส่วนบุคคล", Phone: "0810000003", Email: "coach.wirot@example.com", Active: true},
	}
	for _, trainer := range trainers {
		var existing models.Trainer
		db.Where("code = ?", trainer.Code).FirstOrCreate(&existing, trainer)
	}

	return nil
}
