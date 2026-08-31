package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	BookingID   string        `json:"booking_id" gorm:"primaryKey;type:varchar(50)"`
	BookingDate time.Time     `json:"booking_date" gorm:"not null"`
	CreatedAt   time.Time     `json:"created_at"`
	Status      BookingStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`

	// Foreign Keys & Relationships (ตัด tag gorm ด้านหลังออก)
	UserID string `json:"user_id" gorm:"type:varchar(50);not null"`
	User   *User  `json:"user,omitempty"`

	AreaID string       `json:"area_id" gorm:"type:varchar(50);not null"`
	Area   *SpecialArea `json:"area,omitempty"`

	SlotID string    `json:"slot_id" gorm:"type:varchar(50);not null"`
	Slot   *TimeSlot `json:"slot,omitempty"`
}

func (Booking) TableName() string {
	return "bookings"
}

// BeforeCreate สร้าง BookingID อัตโนมัติหากไม่ได้ระบุมา
func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	if b.BookingID == "" {
		b.BookingID = fmt.Sprintf("BK-%d", time.Now().UnixNano())
	}
	return
}

func (b *Booking) CreateBooking(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	return db.Create(b).Error
}

func (b *Booking) CancelBooking(db *gorm.DB) error {
	b.Status = BookingStatusCancelled
	if db == nil {
		return nil
	}
	return db.Model(b).Update("status", BookingStatusCancelled).Error
}
