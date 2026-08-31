package controllers

import (
	"net/http"

	"SA01/internal/dto"
	"SA01/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecialAreaController struct {
	DB *gorm.DB
}

func NewSpecialAreaController(db *gorm.DB) *SpecialAreaController {
	return &SpecialAreaController{DB: db}
}

func (c *SpecialAreaController) GetSpecialAreas(ctx *gin.Context) {
	selectedDate := ctx.DefaultQuery("date", "")

	// 1. ดึงข้อมูล SpecialArea พร้อม Preload TimeSlots
	var areas []models.SpecialArea
	if err := c.DB.Preload("TimeSlots").Find(&areas).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. ดึงรายการ Booking ตามวันที่ระบุ
	var bookings []models.Booking
	if selectedDate != "" {
		c.DB.Preload("User").Where("DATE(booking_date) = ?", selectedDate).Find(&bookings)
	}

	// Map หา SlotID เพื่อความเร็วในการเช็คสถานะการจอง
	bookingMap := make(map[string]models.Booking)
	for _, b := range bookings {
		bookingMap[b.SlotID] = b
	}

	// 3. แปลง Data เป็น DTO ให้ React UI
	var response []dto.SpecialAreaResponse
	for _, area := range areas {
		var slots []dto.TimeSlotResponse
		for _, slot := range area.TimeSlots {
			isBooked := false
			var bookedBy *string

			if b, exists := bookingMap[slot.SlotID]; exists && b.Status != models.BookingStatusCancelled {
				isBooked = true
				if b.User != nil {
					name := b.User.Name
					bookedBy = &name
				}
			}

			// แปลง time.Time เป็น Format HH:mm (เช่น "08:00 - 09:00")
			timeRange := slot.StartTime.Format("15:04") + " - " + slot.EndTime.Format("15:04")

			slots = append(slots, dto.TimeSlotResponse{
				TimeRange: timeRange,
				IsBooked:  isBooked,
				BookedBy:  bookedBy,
			})
		}

		response = append(response, dto.SpecialAreaResponse{
			Name:      area.Name,
			Capacity:  area.Capacity,
			TimeSlots: slots,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}
