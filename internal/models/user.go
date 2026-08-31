package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID      string    `json:"user_id" gorm:"primaryKey;type:varchar(50)"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(20)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Bookings []Booking `json:"bookings,omitempty" gorm:"foreignKey:UserID"`
}

func (u *User) BookSlot(db *gorm.DB, areaID string, slotID string, date time.Time) (*Booking, error) {
	booking := Booking{
		UserID:      u.UserID,
		AreaID:      areaID,
		SlotID:      slotID,
		BookingDate: date,
		CreatedAt:   time.Now(),
		Status:      BookingStatusPending,
	}

	if err := db.Create(&booking).Error; err != nil {
		return nil, err
	}

	return &booking, nil
}

func (u *User) GetBookingHistory(db *gorm.DB) ([]Booking, error) {
	var history []Booking
	err := db.Where("user_id = ?", u.UserID).
		Order("created_at desc").
		Find(&history).Error

	if err != nil {
		return nil, err
	}

	return history, nil
}

func (User) TableName() string {
	return "users"
}
