package models

import "time"

// User Roles & Statuses
const (
	RoleAdmin     = "admin"
	RoleStaff     = "staff"
	RoleMember    = "member"
	UserActive    = "active"
	UserSuspended = "suspended"
)

// Booking Statuses
const (
	BookingPending   = "pending"
	BookingConfirmed = "confirmed"
	BookingCancelled = "cancelled"
	BookingCompleted = "completed"
)

// Trainer Appointment Statuses
const (
	AppointmentPending   = "pending"
	AppointmentConfirmed = "confirmed"
	AppointmentCancelled = "cancelled"
	AppointmentCompleted = "completed"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FullName     string    `gorm:"size:120;not null" json:"full_name"`
	Email        string    `gorm:"size:180;not null" json:"email"`
	Phone        string    `gorm:"size:30;not null;index" json:"phone"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:20;not null;default:member" json:"role"`
	Status       string    `gorm:"size:20;not null;default:active" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Bookings     []Booking `json:"bookings,omitempty"`
}

type SpecialArea struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"size:40;uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"size:120;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Capacity    int       `gorm:"not null;default:1" json:"capacity"`
	Active      bool      `gorm:"not null;default:true;index" json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TimeSlot struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StartTime string    `gorm:"size:5;not null" json:"start_time"`
	EndTime   string    `gorm:"size:5;not null" json:"end_time"`
	Active    bool      `gorm:"not null;default:true;index" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Booking struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	BookingCode   string      `gorm:"size:24;uniqueIndex;not null" json:"booking_code"`
	UserID        uint        `gorm:"not null;index" json:"user_id"`
	AreaID        uint        `gorm:"not null;index" json:"area_id"`
	TimeSlotID    uint        `gorm:"not null;index" json:"time_slot_id"`
	BookingDate   time.Time   `gorm:"type:date;not null;index" json:"booking_date"`
	Status        string      `gorm:"size:20;not null;default:confirmed;index" json:"status"`
	CooldownUntil *time.Time  `json:"cooldown_until,omitempty"`
	Note          string      `gorm:"type:text" json:"note"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	User          User        `json:"user,omitempty"`
	Area          SpecialArea `json:"area,omitempty"`
	TimeSlot      TimeSlot    `json:"time_slot,omitempty"`
}

type Trainer struct {
	ID           uint                 `gorm:"primaryKey" json:"id"`
	Code         string               `gorm:"size:30;uniqueIndex;not null" json:"code"`
	FullName     string               `gorm:"size:120;not null" json:"full_name"`
	Specialty    string               `gorm:"size:160" json:"specialty"`
	Phone        string               `gorm:"size:30" json:"phone"`
	Email        string               `gorm:"size:180" json:"email"`
	Active       bool                 `gorm:"not null;default:true;index" json:"active"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	Appointments []TrainerAppointment `json:"appointments,omitempty"`
}

type TrainerAppointment struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TrainerID       uint      `gorm:"not null;index" json:"trainer_id"`
	MemberID        uint      `gorm:"not null;index" json:"member_id"`
	AppointmentDate time.Time `gorm:"type:date;not null;index" json:"appointment_date"`
	StartTime       string    `gorm:"size:5;not null" json:"start_time"`
	EndTime         string    `gorm:"size:5;not null" json:"end_time"`
	Status          string    `gorm:"size:20;not null;default:confirmed;index" json:"status"`
	Note            string    `gorm:"type:text" json:"note"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Trainer         Trainer   `json:"trainer,omitempty"`
	Member          User      `json:"member,omitempty"`
}
