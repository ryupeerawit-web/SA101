package main

import (
	"log"
	"time"

	"SA01/internal/config"
	"SA01/internal/models"
	"SA01/internal/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	db := config.ConnectDatabase(cfg)

	// Seed ข้อมูลพื้นที่และสล็อตเวลาตาม UI
	seedUIRequirements(db)

	r := gin.Default()
	routes.SetupRoutes(r, db)

	log.Printf("Server running on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}

func seedUIRequirements(db *gorm.DB) {
	// 1. Seed User
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		users := []models.User{
			{UserID: "U001", Name: "คุณสมชาย ใจดี"},
			{UserID: "U002", Name: "คุณมานี มีจิต"},
		}
		if err := db.Create(&users).Error; err != nil {
			log.Printf("Failed to seed users: %v", err)
		} else {
			log.Println("Seeded users successfully!")
		}
	}

	// 2. Seed Special Areas ตาม UI
	var areaCount int64
	db.Model(&models.SpecialArea{}).Count(&areaCount)
	if areaCount == 0 {
		areas := []models.SpecialArea{
			{AreaID: "SA01", Name: "ห้องโยคะ", Capacity: 15},
			{AreaID: "SA02", Name: "ห้องซาวน่า", Capacity: 8},
			{AreaID: "SA03", Name: "โซนมวย", Capacity: 6},
			{AreaID: "SA04", Name: "สระว่ายน้ำ VIP", Capacity: 10},
		}
		if err := db.Create(&areas).Error; err != nil {
			log.Printf("Failed to seed areas: %v", err)
			return
		}

		// 3. Seed Time Slots (08:00 - 21:00) ให้กับทุกโซน
		now := time.Now()
		timeRanges := []struct {
			startHour int
			endHour   int
			slotNum   string
		}{
			{8, 9, "01"}, {9, 10, "02"}, {10, 11, "03"}, {11, 12, "04"},
			{13, 14, "05"}, {14, 15, "06"}, {15, 16, "07"}, {16, 17, "08"},
			{17, 18, "09"}, {18, 19, "10"}, {19, 20, "11"}, {20, 21, "12"},
		}

		for _, area := range areas {
			var slots []models.TimeSlot
			for _, tr := range timeRanges {
				slots = append(slots, models.TimeSlot{
					SlotID:    area.AreaID + "-SLOT-" + tr.slotNum,
					AreaID:    area.AreaID,
					StartTime: time.Date(now.Year(), now.Month(), now.Day(), tr.startHour, 0, 0, 0, time.Local),
					EndTime:   time.Date(now.Year(), now.Month(), now.Day(), tr.endHour, 0, 0, 0, time.Local),
				})
			}
			if err := db.Create(&slots).Error; err != nil {
				log.Printf("Failed to seed slots for area %s: %v", area.AreaID, err)
			}
		}
		log.Println("Seeded all 4 areas with 12 time slots each successfully!")

		// 4. Seed ตัวอย่างรายการจอง (คุณสมชาย & คุณมานี)
		booking1 := models.Booking{
			BookingID:   "BK001",
			UserID:      "U001",
			AreaID:      "SA01",         // ใส่ AreaID ให้ครอบคลุมถ้าใน Model มี
			SlotID:      "SA01-SLOT-02", // 09:00 - 10:00
			BookingDate: now,
			Status:      models.BookingStatusConfirmed,
		}
		booking2 := models.Booking{
			BookingID:   "BK002",
			UserID:      "U002",
			AreaID:      "SA01",         // ใส่ AreaID ให้ครอบคลุมถ้าใน Model มี
			SlotID:      "SA01-SLOT-06", // 14:00 - 15:00
			BookingDate: now,
			Status:      models.BookingStatusConfirmed,
		}

		if err := booking1.CreateBooking(db); err != nil {
			log.Printf("Failed to create seed booking1: %v", err)
		}
		if err := booking2.CreateBooking(db); err != nil {
			log.Printf("Failed to create seed booking2: %v", err)
		}
	}
}
