package main

import (
	"go-gin-user-account-api/configs"
	"go-gin-user-account-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	app := gin.Default();

	serverConfigs := configs.SetServerConfigurations();

	// connect to database
	configs.ConnectToMongoDb();

	// disable trust all proxies for now as no proxy client has been used to make request to server
	app.SetTrustedProxies(nil)

	// enable cors
	app.Use(configs.SetCors([]string{"http://loaclhost:3000"}))

	// use logger middleware
	app.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	app.Use(gin.Recovery())

	// api/v1 routes group
	apiVersionOne := app.Group("/api/v1")

	// users routes
	routes.UsersRoutes(apiVersionOne)

	// start  the API server
	app.Run(serverConfigs.Port)	

}