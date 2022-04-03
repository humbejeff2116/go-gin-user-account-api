
package routes

import (
	"go-gin-user-account-api/controllers"
	"github.com/gin-gonic/gin"
)

func UsersRoutes(routerVersion *gin.RouterGroup) {
	
	{

		routerVersion.POST("/signup", controllers.SignupUser)
		routerVersion.GET("/login", controllers.LoginUser)
		routerVersion.GET("/user/:userId", controllers.GetUser)
		routerVersion.PUT("/update-user", controllers.UpdateUser);
		routerVersion.DELETE("/user/:userId", controllers.RemoveUser);
		
	}

}