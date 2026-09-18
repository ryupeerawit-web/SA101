package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ryu111/special-area-booking/internal/models"
)

type Dashboard struct {
	DB *gorm.DB
}

type AreaBookingSummary struct {
	AreaID   uint   `json:"area_id" gorm:"column:area_id"`
	Code     string `json:"code" gorm:"column:code"`
	Name     string `json:"name" gorm:"column:name"`
	Capacity int    `json:"capacity" gorm:"column:capacity"`
	Booked   int64  `json:"booked" gorm:"column:booked"`
}

// Summary provides aggregate metrics and per-area booking stats for a given date.
func (d Dashboard) Summary(c *gin.Context) {
	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	var (
		total     int64
		confirmed int64
		cancelled int64
		today     int64
	)

	d.DB.Model(&models.Booking{}).Count(&total)

	d.DB.Model(&models.Booking{}).
		Where("status = ?", models.BookingConfirmed).
		Count(&confirmed)

	d.DB.Model(&models.Booking{}).
		Where("status = ?", models.BookingCancelled).
		Count(&cancelled)

	d.DB.Model(&models.Booking{}).
		Where("booking_date = ? AND status IN ?", date, []string{models.BookingPending, models.BookingConfirmed}).
		Count(&today)

	var areaSummaries []AreaBookingSummary
	rawQuery := `
		SELECT 
			a.id AS area_id,
			a.code,
			a.name,
			a.capacity,
			COUNT(b.id) FILTER (WHERE b.status IN ('pending', 'confirmed')) AS booked
		FROM special_areas a
		LEFT JOIN bookings b ON b.area_id = a.id AND b.booking_date = ?
		GROUP BY a.id
		ORDER BY a.id
	`

	if err := d.DB.Raw(rawQuery, date).Scan(&areaSummaries).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "could not fetch dashboard summary")
		return
	}

	OK(c, http.StatusOK, gin.H{
		"date":               date,
		"total_bookings":     total,
		"confirmed_bookings": confirmed,
		"cancelled_bookings": cancelled,
		"today_bookings":     today,
		"by_area":            areaSummaries,
	})
}
