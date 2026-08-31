package routes

import (
	"SA01/internal/controllers"
	"SA01/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// เปิดใช้งาน CORS Middleware ให้ทุก Endpoint
	r.Use(middleware.CORSMiddleware())

	// Initialize Controllers
	specialAreaCtrl := controllers.NewSpecialAreaController(db)
	bookingCtrl := controllers.NewBookingController(db)

	// Route Group สำหรับ API Version 1
	api := r.Group("/api/v1")
	{
		api.GET("/special-areas", specialAreaCtrl.GetSpecialAreas)
		api.POST("/bookings", bookingCtrl.CreateBooking)
	}
}
