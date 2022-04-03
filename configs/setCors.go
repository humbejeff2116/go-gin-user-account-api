
package configs

import (
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func SetCors(origins[]string) gin.HandlerFunc {

	return cors.New(cors.Config {
		AllowOrigins:     origins,
		AllowMethods:     []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	})

}