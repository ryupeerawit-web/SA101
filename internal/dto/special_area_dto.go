package dto

type TimeSlotResponse struct {
	ID        uint    `json:"id"`
	TimeRange string  `json:"time_range"` // เช่น "08:00 - 09:00"
	IsBooked  bool    `json:"is_booked"`
	BookedBy  *string `json:"booked_by"`  // null ถ้ายังไม่จอง
}

type SpecialAreaResponse struct {
	ID        uint               `json:"id"`
	Name      string             `json:"name"`
	Capacity  int                `json:"capacity"`
	TimeSlots []TimeSlotResponse `json:"time_slots"`
}