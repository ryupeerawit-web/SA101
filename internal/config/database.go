package config

import (
	"log"

	"SA01/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// เรียงลำดับการสร้างตาราง: ตารางหลักต้องขึ้นก่อน ตารางที่มี Foreign Key (Booking)
	err = db.AutoMigrate(
		&models.User{},        // 1. สร้างตาราง users ก่อน
		&models.SpecialArea{}, // 2. สร้างตาราง special_areas
		&models.AreaInfo{},
		&models.TimeSlot{}, // 3. สร้างตาราง time_slots
		&models.Booking{},  // 5. สร้างตาราง bookings เป็นลำดับสุดท้าย (เพราะอ้างอิง User และ Slot)
	)
	if err != nil {
		log.Fatal("Auto Migration failed:", err)
	}

	log.Println("Database connection & migration completed successfully!")
	return db
}
