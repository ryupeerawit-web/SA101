package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ryu111/special-area-booking/internal/config"
	"github.com/ryu111/special-area-booking/internal/utils"
)

// CORS configures Access-Control headers to handle cross-origin requests.
func CORS(c config.Config) gin.HandlerFunc {
	return func(x *gin.Context) {
		origin := x.GetHeader("Origin")
		for _, v := range c.CORSOrigins {
			if v == "*" || v == origin {
				x.Header("Access-Control-Allow-Origin", origin)
				x.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		x.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		x.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if x.Request.Method == http.MethodOptions {
			x.AbortWithStatus(http.StatusNoContent)
			return
		}

		x.Next()
	}
}

// Auth validates JWT Bearer tokens and extracts user payload to context.
func Auth(c config.Config) gin.HandlerFunc {
	return func(x *gin.Context) {
		parts := strings.SplitN(x.GetHeader("Authorization"), " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			x.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		claims, err := utils.ParseToken(parts[1], c)
		if err != nil {
			x.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			x.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token subject"})
			return
		}

		x.Set("user_id", uint(sub))
		if role, ok := claims["role"].(string); ok {
			x.Set("user_role", role)
		}

		x.Next()
	}
}

// Roles restricts access to specific user roles.
func Roles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("user_role")
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}

// UserID retrieves the authenticated user's ID from context.
func UserID(c *gin.Context) uint {
	v, _ := c.Get("user_id")
	return v.(uint)
}

// ID converts a URL path string parameter into a uint ID.
func ID(s string) (uint, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	return uint(n), err
}
