package routers

import (
	"cloud-posture-scanner/controllers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouters() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/cloud-posture-scanner/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/instances", controllers.GetEC2Instances)
		api.GET("/buckets", controllers.GetS3Buckets)
		api.GET("/cis-results", controllers.GetCISResults)
	}

	return r
}
