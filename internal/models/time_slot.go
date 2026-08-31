package models

import (
	"time"

	"gorm.io/gorm"
)

type TimeSlot struct {
	SlotID    string     `json:"slot_id" gorm:"primaryKey;type:varchar(50)"`
	AreaID    string     `json:"area_id" gorm:"type:varchar(50);not null"`
	StartTime time.Time  `json:"start_time" gorm:"not null"`
	EndTime   time.Time  `json:"end_time" gorm:"not null"`
	Status    SlotStatus `json:"status" gorm:"type:varchar(20);default:'AVAILABLE'"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// ✅ เลือกเชื่อมกับ SpecialArea เพียงอย่างเดียว และเอา Tag gorm foreignKey ออก
	Area *SpecialArea `json:"area,omitempty"`
}

// IsAvailable ตรวจสอบว่าช่วงเวลานี้ว่างหรือไม่
func (ts *TimeSlot) IsAvailable() bool {
	return ts.Status == SlotStatusAvailable
}

// UpdateStatus อัปเดตสถานะของช่วงเวลาลง Database จริง
func (ts *TimeSlot) UpdateStatus(db *gorm.DB, newStatus SlotStatus) error {
	ts.Status = newStatus
	// บันทึกเฉพาะ field Status ลง DB
	return db.Model(ts).Update("status", newStatus).Error
}

func (TimeSlot) TableName() string {
	return "time_slots"
}
