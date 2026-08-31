package dto

// CreateBookingRequest รับข้อมูลสำหรับสร้างรายการจองใหม่จาก Frontend
type CreateBookingRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	SlotID      string `json:"slot_id" binding:"required"`
	AreaID      string `json:"area_id" binding:"required"`
	BookingDate string `json:"booking_date" binding:"required"` // รูปแบบ "YYYY-MM-DD"
}

// BookingResponse รูปแบบข้อมูลการจองที่จะส่งกลับไปให้ Frontend แสดงผล
type BookingResponse struct {
	BookingID   string `json:"booking_id"`
	UserID      string `json:"user_id"`
	AreaID      string `json:"area_id"`
	SlotID      string `json:"slot_id"`
	BookingDate string `json:"booking_date"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}