package controllers

import (
	"net/http"
	"time"

	"SA01/internal/dto"
	"SA01/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingController struct {
	DB *gorm.DB
}

func NewBookingController(db *gorm.DB) *BookingController {
	return &BookingController{DB: db}
}

func (c *BookingController) CreateBooking(ctx *gin.Context) {
	var req dto.CreateBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. ตรวจสอบว่ามี User นี้ในระบบหรือไม่
	var user models.User
	if err := c.DB.First(&user, "user_id = ?", req.UserID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลผู้ใช้นี้"})
		return
	}

	// 2. แปลง String Date เป็น time.Time
	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "รูปแบบวันที่ไม่ถูกต้อง (YYYY-MM-DD)"})
		return
	}

	// 3. ตรวจสอบว่าสล็อตเวลานี้ถูกจองไปแล้วในวันที่เลือกหรือไม่
	var existingBooking models.Booking
	err = c.DB.Where("slot_id = ? AND DATE(booking_date) = ? AND status != ?", req.SlotID, req.BookingDate, models.BookingStatusCancelled).First(&existingBooking).Error
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "สล็อตเวลานี้ถูกจองไปแล้ว"})
		return
	}

	// 4. บันทึกข้อมูลการจองลงฐานข้อมูล
	booking, err := user.BookSlot(c.DB, req.AreaID, req.SlotID, bookingDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "สร้างคำขอจองเรียบร้อยแล้ว",
		"data":    booking,
	})
}
