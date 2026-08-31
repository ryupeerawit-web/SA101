package models

// SlotStatus Represent status of a TimeSlot
type SlotStatus string

const (
	SlotStatusAvailable   SlotStatus = "AVAILABLE"
	SlotStatusBooked      SlotStatus = "BOOKED"
	SlotStatusUnavailable SlotStatus = "UNAVAILABLE"
)

// BookingStatus Represent status of a Booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
)
