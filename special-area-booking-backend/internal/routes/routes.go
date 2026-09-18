package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ryu111/special-area-booking/internal/config"
	"github.com/ryu111/special-area-booking/internal/controllers"
	"github.com/ryu111/special-area-booking/internal/middleware"
)

// Setup initializes and configures HTTP routes for the application.
func Setup(db *gorm.DB, c config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(c))

	// Health check endpoint
	r.GET("/api/v1/health", func(x *gin.Context) {
		x.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Controllers initialization
	authController := controllers.Auth{DB: db, Cfg: c}
	areasController := controllers.Areas{DB: db}
	bookingsController := controllers.Bookings{DB: db}
	dashboardController := controllers.Dashboard{DB: db}
	trainersController := controllers.Trainers{DB: db}

	v1 := r.Group("/api/v1")

	// Public routes
	v1.POST("/auth/register", authController.Register)
	v1.POST("/auth/login", authController.Login)
	v1.GET("/areas", areasController.List)
	v1.GET("/areas/:id", areasController.Get)
	v1.GET("/slots", areasController.Slots)
	v1.GET("/availability", bookingsController.Availability)

	// Authenticated routes
	secured := v1.Group("")
	secured.Use(middleware.Auth(c))

	// User profile & bookings
	secured.GET("/me", authController.Me)
	secured.POST("/bookings", bookingsController.Create)
	secured.GET("/bookings", bookingsController.Mine)
	secured.GET("/bookings/:id", bookingsController.Get)
	secured.PATCH("/bookings/:id/cancel", bookingsController.Cancel)

	// Trainer browsing & appointments
	secured.GET("/trainers", trainersController.List)
	secured.GET("/trainers/:id", trainersController.Get)
	secured.GET("/trainer-appointments", trainersController.Appointments)
	secured.POST("/trainer-appointments", trainersController.CreateAppointment)
	secured.PATCH("/trainer-appointments/:id/cancel", trainersController.CancelAppointment)

	// Admin / Staff routes
	admin := secured.Group("/admin")
	admin.Use(middleware.Roles("admin", "staff"))

	admin.GET("/dashboard/summary", dashboardController.Summary)
	admin.GET("/bookings", bookingsController.AdminList)
	admin.PATCH("/bookings/:id/status", bookingsController.Status)
	admin.POST("/areas", areasController.Create)
	admin.PUT("/areas/:id", areasController.Update)
	admin.POST("/slots", areasController.CreateSlot)
	admin.POST("/trainers", trainersController.Create)
	admin.PUT("/trainers/:id", trainersController.Update)
	admin.GET("/members", trainersController.Members)
	admin.POST("/trainer-appointments", trainersController.CreateAppointment)
	admin.PATCH("/trainer-appointments/:id/status", trainersController.UpdateAppointmentStatus)

	return r
}
