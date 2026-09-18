package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ryu111/special-area-booking/internal/dto"
	"github.com/ryu111/special-area-booking/internal/models"
	"gorm.io/gorm"
)

type Areas struct {
	DB *gorm.DB
}

func (a Areas) List(c *gin.Context) {
	query := a.DB.Order("id asc")
	if c.Query("active") != "false" {
		query = query.Where("active = ?", true)
	}

	var areas []models.SpecialArea
	if err := query.Find(&areas).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "could not load areas")
		return
	}

	OK(c, http.StatusOK, areas)
}

func (a Areas) Get(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid area id")
		return
	}

	var area models.SpecialArea
	if err := a.DB.First(&area, id).Error; err != nil {
		Fail(c, http.StatusNotFound, "area not found")
		return
	}

	OK(c, http.StatusOK, area)
}

func (a Areas) Create(c *gin.Context) {
	var req dto.CreateAreaRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Capacity < 1 {
		Fail(c, http.StatusBadRequest, "invalid request")
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	area := models.SpecialArea{
		Code:        strings.ToUpper(strings.TrimSpace(req.Code)),
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		Active:      active,
	}

	if err := a.DB.Create(&area).Error; err != nil {
		Fail(c, http.StatusConflict, "area code already exists or data is invalid")
		return
	}

	OK(c, http.StatusCreated, area)
}

func (a Areas) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid area id")
		return
	}

	var req dto.UpdateAreaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, "invalid request")
		return
	}

	var area models.SpecialArea
	if err := a.DB.First(&area, id).Error; err != nil {
		Fail(c, http.StatusNotFound, "area not found")
		return
	}

	if req.Name != "" {
		area.Name = req.Name
	}
	if req.Description != "" {
		area.Description = req.Description
	}
	if req.Capacity > 0 {
		area.Capacity = req.Capacity
	}
	if req.Active != nil {
		area.Active = *req.Active
	}

	if err := a.DB.Save(&area).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "could not update area")
		return
	}

	OK(c, http.StatusOK, area)
}

func (a Areas) Slots(c *gin.Context) {
	var slots []models.TimeSlot
	if err := a.DB.Where("active = ?", true).Order("start_time").Find(&slots).Error; err != nil {
		Fail(c, http.StatusInternalServerError, "could not load slots")
		return
	}

	OK(c, http.StatusOK, slots)
}

func (a Areas) CreateSlot(c *gin.Context) {
	var req dto.CreateSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, "invalid request")
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	slot := models.TimeSlot{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Active:    active,
	}

	if err := a.DB.Create(&slot).Error; err != nil {
		Fail(c, http.StatusConflict, "time slot already exists or data is invalid")
		return
	}

	OK(c, http.StatusCreated, slot)
}

func parseID(s string) (uint, error) {
	var id uint
	_, err := fmt.Sscan(s, &id)
	return id, err
}
