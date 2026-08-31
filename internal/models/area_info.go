package models

import "gorm.io/gorm"

type AreaInfo struct {
	gorm.Model
	AreaID    string     `json:"area_id" gorm:"type:varchar(50);uniqueIndex;not null"`
	AreaName  string     `json:"name" gorm:"type:varchar(100);not null"`
	Capacity  int        `json:"capacity" gorm:"not null"`
	Description string `json:"description"`
	TimeSlots []TimeSlot `json:"time_slots,omitempty" gorm:"foreignKey:AreaID;references:AreaID"`
}