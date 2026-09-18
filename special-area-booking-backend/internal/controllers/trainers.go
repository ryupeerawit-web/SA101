package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryu111/special-area-booking/internal/dto"
	"github.com/ryu111/special-area-booking/internal/middleware"
	"github.com/ryu111/special-area-booking/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Trainers struct{ DB *gorm.DB }

func (t Trainers) List(c *gin.Context) {
	q := t.DB.Order("id asc")
	if c.Query("active") != "false" {
		q = q.Where("active = ?", true)
	}
	if v := strings.TrimSpace(c.Query("search")); v != "" {
		like := "%" + v + "%"
		q = q.Where("full_name ILIKE ? OR code ILIKE ? OR specialty ILIKE ?", like, like, like)
	}
	var rows []models.Trainer
	if q.Find(&rows).Error != nil {
		Fail(c, 500, "could not load trainers")
		return
	}
	OK(c, 200, rows)
}

func (t Trainers) Get(c *gin.Context) {
	id, e := middleware.ID(c.Param("id"))
	if e != nil {
		Fail(c, 400, "invalid trainer id")
		return
	}
	var row models.Trainer
	if t.DB.First(&row, id).Error != nil {
		Fail(c, 404, "trainer not found")
		return
	}
	OK(c, 200, row)
}

func (t Trainers) Create(c *gin.Context) {
	var r dto.CreateTrainerRequest
	if c.ShouldBindJSON(&r) != nil {
		Fail(c, 400, "invalid request")
		return
	}
	active := true
	if r.Active != nil {
		active = *r.Active
	}
	row := models.Trainer{
		Code:      strings.ToUpper(strings.TrimSpace(r.Code)),
		FullName:  r.FullName,
		Specialty: r.Specialty,
		Phone:     r.Phone,
		Email:     r.Email,
		Active:    active,
	}
	if t.DB.Create(&row).Error != nil {
		Fail(c, 409, "trainer code already exists or data is invalid")
		return
	}
	OK(c, http.StatusCreated, row)
}

func (t Trainers) Update(c *gin.Context) {
	id, e := middleware.ID(c.Param("id"))
	if e != nil {
		Fail(c, 400, "invalid trainer id")
		return
	}
	var r dto.UpdateTrainerRequest
	if c.ShouldBindJSON(&r) != nil {
		Fail(c, 400, "invalid request")
		return
	}
	var row models.Trainer
	if t.DB.First(&row, id).Error != nil {
		Fail(c, 404, "trainer not found")
		return
	}
	if r.FullName != "" {
		row.FullName = r.FullName
	}
	if r.Specialty != "" {
		row.Specialty = r.Specialty
	}
	if r.Phone != "" {
		row.Phone = r.Phone
	}
	if r.Email != "" {
		row.Email = r.Email
	}
	if r.Active != nil {
		row.Active = *r.Active
	}
	if t.DB.Save(&row).Error != nil {
		Fail(c, 500, "could not update trainer")
		return
	}
	OK(c, 200, row)
}

func (t Trainers) Appointments(c *gin.Context) {
	from := c.DefaultQuery("from", time.Now().Format("2006-01-02"))
	to := c.DefaultQuery("to", from)
	var rows []models.TrainerAppointment
	q := t.DB.Preload("Trainer").Preload("Member").Where("appointment_date BETWEEN ? AND ?", from, to).Order("appointment_date asc,start_time asc,id asc")
	if c.GetString("user_role") != models.RoleAdmin && c.GetString("user_role") != models.RoleStaff {
		q = q.Where("member_id = ?", middleware.UserID(c))
	}
	if v := c.Query("trainer_id"); v != "" {
		q = q.Where("trainer_id = ?", v)
	}
	if q.Find(&rows).Error != nil {
		Fail(c, 500, "could not load trainer appointments")
		return
	}
	OK(c, 200, rows)
}

func (t Trainers) CreateAppointment(c *gin.Context) {
	var r dto.CreateTrainerAppointmentRequest
	if c.ShouldBindJSON(&r) != nil {
		Fail(c, 400, "invalid request")
		return
	}
	date, e := time.ParseInLocation("2006-01-02", r.AppointmentDate, time.Local)
	if e != nil {
		Fail(c, 400, "appointment_date must be YYYY-MM-DD")
		return
	}
	memberID := middleware.UserID(c)
	if r.MemberID > 0 && (c.GetString("user_role") == models.RoleAdmin || c.GetString("user_role") == models.RoleStaff) {
		memberID = r.MemberID
	}
	var row models.TrainerAppointment
	e = t.DB.Transaction(func(tx *gorm.DB) error {
		var trainer models.Trainer
		if tx.Where("id = ? AND active = ?", r.TrainerID, true).First(&trainer).Error != nil {
			return fmt.Errorf("trainer not found")
		}
		var existing models.TrainerAppointment
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("trainer_id=? AND appointment_date=? AND start_time=? AND end_time=? AND status IN ?", r.TrainerID, date, r.StartTime, r.EndTime, []string{models.AppointmentPending, models.AppointmentConfirmed}).First(&existing).Error
		if err == nil {
			return fmt.Errorf("trainer slot already booked")
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}
		row = models.TrainerAppointment{TrainerID: r.TrainerID, MemberID: memberID, AppointmentDate: date, StartTime: r.StartTime, EndTime: r.EndTime, Status: models.AppointmentConfirmed, Note: r.Note}
		return tx.Create(&row).Error
	})
	if e != nil {
		if e.Error() == "trainer slot already booked" {
			Fail(c, 409, e.Error())
			return
		}
		if e.Error() == "trainer not found" {
			Fail(c, 404, e.Error())
			return
		}
		Fail(c, 500, "could not create trainer appointment")
		return
	}
	t.DB.Preload("Trainer").Preload("Member").First(&row, row.ID)
	OK(c, http.StatusCreated, row)
}

func (t Trainers) CancelAppointment(c *gin.Context) {
	id, e := middleware.ID(c.Param("id"))
	if e != nil {
		Fail(c, 400, "invalid appointment id")
		return
	}
	var row models.TrainerAppointment
	if t.DB.First(&row, id).Error != nil {
		Fail(c, 404, "trainer appointment not found")
		return
	}
	role := c.GetString("user_role")
	if row.MemberID != middleware.UserID(c) && role != models.RoleAdmin && role != models.RoleStaff {
		Fail(c, 403, "not allowed")
		return
	}
	if row.Status == models.AppointmentCancelled || row.Status == models.AppointmentCompleted {
		Fail(c, 409, "appointment cannot be cancelled")
		return
	}
	row.Status = models.AppointmentCancelled
	if t.DB.Save(&row).Error != nil {
		Fail(c, 500, "could not cancel appointment")
		return
	}
	OK(c, 200, row)
}

func (t Trainers) UpdateAppointmentStatus(c *gin.Context) {
	id, e := middleware.ID(c.Param("id"))
	if e != nil {
		Fail(c, 400, "invalid appointment id")
		return
	}
	var r dto.UpdateTrainerAppointmentStatusRequest
	if c.ShouldBindJSON(&r) != nil {
		Fail(c, 400, "invalid request")
		return
	}
	var row models.TrainerAppointment
	if t.DB.First(&row, id).Error != nil {
		Fail(c, 404, "trainer appointment not found")
		return
	}
	row.Status = r.Status
	if t.DB.Save(&row).Error != nil {
		Fail(c, 500, "could not update appointment")
		return
	}
	OK(c, 200, row)
}

func (t Trainers) Members(c *gin.Context) {
	q := t.DB.Where("role = ?", models.RoleMember).Order("full_name asc")
	if v := strings.TrimSpace(c.Query("search")); v != "" {
		like := "%" + v + "%"
		q = q.Where("full_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", like, like, like)
	}
	var rows []models.User
	if q.Find(&rows).Error != nil {
		Fail(c, 500, "could not load members")
		return
	}
	OK(c, 200, rows)
}
