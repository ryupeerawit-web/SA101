package controllers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ryu111/special-area-booking/internal/dto"
	"github.com/ryu111/special-area-booking/internal/middleware"
	"github.com/ryu111/special-area-booking/internal/models"
)

// Custom domain errors for clean control flow in transactions
var (
	errAreaNotFound      = errors.New("area not found")
	errSlotNotFound      = errors.New("slot not found")
	errSlotAlreadyBooked = errors.New("slot already booked")
)

type Bookings struct {
	DB *gorm.DB
}

// Create handles creating a new booking with concurrency locking and cooldown validation.
func (b Bookings) Create(c *gin.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, "invalid request")
		return
	}

	date, err := time.ParseInLocation("2006-01-02", req.BookingDate, time.Local)
	if err != nil {
		Fail(c, http.StatusBadRequest, "booking_date must be YYYY-MM-DD")
		return
	}

	if date.Before(truncateToDate(time.Now())) {
		Fail(c, http.StatusBadRequest, "booking date cannot be in the past")
		return
	}

	userID := middleware.UserID(c)
	var newBooking models.Booking

	err = b.DB.Transaction(func(tx *gorm.DB) error {
		var area models.SpecialArea
		if err := tx.Where("id = ? AND active = ?", req.AreaID, true).First(&area).Error; err != nil {
			return errAreaNotFound
		}

		var slot models.TimeSlot
		if err := tx.Where("id = ? AND active = ?", req.TimeSlotID, true).First(&slot).Error; err != nil {
			return errSlotNotFound
		}

		// Check for active cancel cooldown period
		var cooldown models.Booking
		cooldownErr := tx.Where(
			"user_id = ? AND area_id = ? AND time_slot_id = ? AND booking_date = ? AND status = ? AND cooldown_until > ?",
			userID, req.AreaID, req.TimeSlotID, date, models.BookingCancelled, time.Now(),
		).Order("cooldown_until desc").First(&cooldown).Error

		if cooldownErr == nil && cooldown.CooldownUntil != nil {
			return fmt.Errorf("cooldown_until:%s", cooldown.CooldownUntil.Format(time.RFC3339))
		}
		if cooldownErr != nil && !errors.Is(cooldownErr, gorm.ErrRecordNotFound) {
			return cooldownErr
		}

		// Acquire row-level lock to prevent double booking race conditions
		var existing models.Booking
		lockErr := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("area_id = ? AND time_slot_id = ? AND booking_date = ? AND status IN ?",
				req.AreaID, req.TimeSlotID, date, []string{models.BookingPending, models.BookingConfirmed},
			).First(&existing).Error

		if lockErr == nil {
			return errSlotAlreadyBooked
		}
		if !errors.Is(lockErr, gorm.ErrRecordNotFound) {
			return lockErr
		}

		newBooking = models.Booking{
			BookingCode: generateBookingCode(),
			UserID:      userID,
			AreaID:      req.AreaID,
			TimeSlotID:  req.TimeSlotID,
			BookingDate: date,
			Status:      models.BookingConfirmed,
			Note:        req.Note,
		}

		return tx.Create(&newBooking).Error
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "cooldown_until:") {
			until := strings.TrimPrefix(err.Error(), "cooldown_until:")
			c.JSON(http.StatusConflict, gin.H{
				"error":          "บัญชีนี้ต้องรอ cooldown ก่อนจองสล็อตเดิม",
				"cooldown_until": until,
			})
			return
		}

		switch {
		case errors.Is(err, errSlotAlreadyBooked):
			Fail(c, http.StatusConflict, err.Error())
		case errors.Is(err, errAreaNotFound), errors.Is(err, errSlotNotFound):
			Fail(c, http.StatusNotFound, err.Error())
		default:
			Fail(c, http.StatusInternalServerError, "could not create booking")
		}
		return
	}

	b.DB.Preload("Area").Preload("TimeSlot").Preload("User").First(&newBooking, newBooking.ID)
	OK(c, http.StatusCreated, newBooking)
}

// Mine lists the authenticated user's bookings.
func (b Bookings) Mine(c *gin.Context) {
	var bookings []models.Booking
	query := b.DB.Preload("Area").Preload("TimeSlot").
		Where("user_id = ?", middleware.UserID(c)).
		Order("booking_date desc, id desc")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&bookings).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "could not load bookings")
		return
	}

	OK(c, http.StatusOK, bookings)
}

// Get fetches details of a single booking by ID.
func (b Bookings) Get(c *gin.Context) {
	bookingID, err := middleware.ID(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid booking id")
		return
	}

	var booking models.Booking
	if err := b.DB.Preload("Area").Preload("TimeSlot").Preload("User").First(&booking, bookingID).Error; err != nil {
		Fail(c, http.StatusNotFound, "booking not found")
		return
	}

	userRole := c.GetString("user_role")
	if userRole != models.RoleAdmin && userRole != models.RoleStaff && booking.UserID != middleware.UserID(c) {
		Fail(c, http.StatusForbidden, "not allowed")
		return
	}

	OK(c, http.StatusOK, booking)
}

// Cancel allows a user or staff member to cancel a booking.
func (b Bookings) Cancel(c *gin.Context) {
	bookingID, err := middleware.ID(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid booking id")
		return
	}

	var booking models.Booking
	if err := b.DB.First(&booking, bookingID).Error; err != nil {
		Fail(c, http.StatusNotFound, "booking not found")
		return
	}

	userRole := c.GetString("user_role")
	if booking.UserID != middleware.UserID(c) && userRole != models.RoleAdmin && userRole != models.RoleStaff {
		Fail(c, http.StatusForbidden, "not allowed")
		return
	}

	if booking.Status == models.BookingCancelled || booking.Status == models.BookingCompleted {
		Fail(c, http.StatusConflict, "booking cannot be cancelled")
		return
	}

	cooldownUntil := time.Now().Add(5 * time.Minute)
	booking.Status = models.BookingCancelled
	booking.CooldownUntil = &cooldownUntil

	if err := b.DB.Save(&booking).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "could not cancel booking")
		return
	}

	OK(c, http.StatusOK, booking)
}

// AdminList lists all system bookings with optional filters for admins/staff.
func (b Bookings) AdminList(c *gin.Context) {
	var bookings []models.Booking
	query := b.DB.Preload("Area").Preload("TimeSlot").Preload("User").
		Order("booking_date desc, id desc")

	if date := c.Query("date"); date != "" {
		query = query.Where("booking_date = ?", date)
	}
	if areaID := c.Query("area_id"); areaID != "" {
		query = query.Where("area_id = ?", areaID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&bookings).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "could not load bookings")
		return
	}

	OK(c, http.StatusOK, bookings)
}

// Status updates the status of a booking by an admin/staff.
func (b Bookings) Status(c *gin.Context) {
	bookingID, err := middleware.ID(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid booking id")
		return
	}

	var req dto.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, "invalid request")
		return
	}

	var booking models.Booking
	if err := b.DB.First(&booking, bookingID).Error; err != nil {
		Fail(c, http.StatusNotFound, "booking not found")
		return
	}

	booking.Status = req.Status
	if err := b.DB.Save(&booking).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "could not update booking")
		return
	}

	OK(c, http.StatusOK, booking)
}

// Availability retrieves time slot availability for areas on a specific date.
func (b Bookings) Availability(c *gin.Context) {
	date := c.Query("date")
	if _, err := time.Parse("2006-01-02", date); err != nil {
		Fail(c, http.StatusBadRequest, "date must be YYYY-MM-DD")
		return
	}

	var areas []models.SpecialArea
	areaQuery := b.DB.Where("active = ?", true)
	if areaID := c.Query("area_id"); areaID != "" {
		areaQuery = areaQuery.Where("id = ?", areaID)
	}
	areaQuery.Find(&areas)

	var slots []models.TimeSlot
	b.DB.Where("active = ?", true).Order("start_time").Find(&slots)

	items := make([]gin.H, 0, len(areas))
	for _, area := range areas {
		slotItems := make([]gin.H, 0, len(slots))
		for _, slot := range slots {
			var bookedCount int64
			b.DB.Model(&models.Booking{}).
				Where("area_id = ? AND time_slot_id = ? AND booking_date = ? AND status IN ?",
					area.ID, slot.ID, date, []string{models.BookingPending, models.BookingConfirmed},
				).
				Count(&bookedCount)

			slotItems = append(slotItems, gin.H{
				"time_slot": slot,
				"available": bookedCount == 0,
			})
		}

		items = append(items, gin.H{
			"area":  area,
			"slots": slotItems,
		})
	}

	OK(c, http.StatusOK, gin.H{
		"date":  date,
		"items": items,
	})
}

// generateBookingCode creates a unique random booking identifier.
func generateBookingCode() string {
	bytes := make([]byte, 6)
	_, _ = rand.Read(bytes)
	return fmt.Sprintf("BK-%X", bytes)
}

// truncateToDate resets hours, minutes, and seconds from a timestamp for exact date comparison.
func truncateToDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
