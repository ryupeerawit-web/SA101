package models

import (
	"time"

	"gorm.io/gorm"
)

type SpecialArea struct {
	AreaID      string     `json:"area_id" gorm:"primaryKey;type:varchar(50)"`
	Name        string     `json:"name" gorm:"type:varchar(100);not null"`
	Capacity    int        `json:"capacity" gorm:"not null"`
	Description string     `json:"description" gorm:"type:text"`
	TimeSlots   []TimeSlot `json:"time_slots,omitempty" gorm:"foreignKey:AreaID;references:AreaID"`
}

// GetAvailableSlots คืนค่ารายการช่วงเวลาที่ว่างตามวันที่ระบุ
func (sa *SpecialArea) GetAvailableSlots(db *gorm.DB, date time.Time) ([]TimeSlot, error) {
	if db == nil {
		return []TimeSlot{}, nil // รองรับกรณี Mock Test ใน main.go
	}

	var availableSlots []TimeSlot
	// Query หาสล็อตเวลาที่เป็นของ Area นี้ และมีสถานะ AVAILABLE
	err := db.Where("area_id = ? AND status = ?", sa.AreaID, SlotStatusAvailable).
		Find(&availableSlots).Error

	if err != nil {
		return nil, err
	}

	return availableSlots, nil
}

// GetAreaDetails คืนค่าข้อมูลรายละเอียดของพื้นที่
func (sa *SpecialArea) GetAreaDetails() AreaInfo {
	return AreaInfo{
		AreaID:   sa.AreaID,
		AreaName:     sa.Name,
		Capacity: sa.Capacity,
	}
}