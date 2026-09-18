package dto

// Authentication Requests

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=120"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,min=7,max=30"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Special Area Requests

type CreateAreaRequest struct {
	Code        string `json:"code" binding:"required,max=40"`
	Name        string `json:"name" binding:"required,max=120"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity" binding:"required,min=1"`
	Active      *bool  `json:"active"`
}

type UpdateAreaRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
	Active      *bool  `json:"active"`
}

// Time Slot Requests

type CreateSlotRequest struct {
	StartTime string `json:"start_time" binding:"required,len=5"`
	EndTime   string `json:"end_time" binding:"required,len=5"`
	Active    *bool  `json:"active"`
}

// Booking Requests

type CreateBookingRequest struct {
	AreaID      uint   `json:"area_id" binding:"required"`
	TimeSlotID  uint   `json:"time_slot_id" binding:"required"`
	BookingDate string `json:"booking_date" binding:"required,len=10"`
	Note        string `json:"note" binding:"max=500"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=confirmed cancelled completed"`
}

// Trainer Requests

type CreateTrainerRequest struct {
	Code      string `json:"code" binding:"required,max=30"`
	FullName  string `json:"full_name" binding:"required,min=2,max=120"`
	Specialty string `json:"specialty"`
	Phone     string `json:"phone"`
	Email     string `json:"email" binding:"omitempty,email"`
	Active    *bool  `json:"active"`
}

type UpdateTrainerRequest struct {
	FullName  string `json:"full_name"`
	Specialty string `json:"specialty"`
	Phone     string `json:"phone"`
	Email     string `json:"email" binding:"omitempty,email"`
	Active    *bool  `json:"active"`
}

// Trainer Appointment Requests

type CreateTrainerAppointmentRequest struct {
	TrainerID       uint   `json:"trainer_id" binding:"required"`
	MemberID        uint   `json:"member_id"`
	AppointmentDate string `json:"appointment_date" binding:"required,len=10"`
	StartTime       string `json:"start_time" binding:"required,len=5"`
	EndTime         string `json:"end_time" binding:"required,len=5"`
	Note            string `json:"note" binding:"max=500"`
}

type UpdateTrainerAppointmentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=confirmed cancelled completed"`
}
