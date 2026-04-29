package controllers

import (
	"cloud-posture-scanner/services"

	"github.com/gin-gonic/gin"
)

// GetEC2Instances godoc
// @Summary Get all EC2 instances
// @Description Get all EC2 instances
// @Tags EC2
// @Accept json
// @Produce json
// @Success 200 {object} models.APIRes{data=[]models.EC2Instance}
// @Failure 500 {object} models.APIRes
// @Router /api/instances [get]
func GetEC2Instances(c *gin.Context) {
	statusCode, apiRes := services.GetEC2Instances()
	c.JSON(statusCode, apiRes)
}

// GetS3Buckets godoc
// @Summary Get all S3 buckets
// @Description Get all S3 buckets
// @Tags S3
// @Accept json
// @Produce json
// @Success 200 {object} models.APIRes{data=[]models.S3Bucket}
// @Failure 500 {object} models.APIRes
// @Router /api/buckets [get]
func GetS3Buckets(c *gin.Context) {
	statusCode, apiRes := services.GetS3Buckets()
	c.JSON(statusCode, apiRes)
}

// RunCisAwsChecks godoc
// @Summary Run CIS AWS checks
// @Description Run CIS AWS checks
// @Tags CIS
// @Accept json
// @Produce json
// @Success 200 {object} models.APIRes{data=[]models.CISCheck}
// @Failure 500 {object} models.APIRes
// @Router /api/cis-results [get]
func GetCISResults(c *gin.Context) {
	statusCode, apiRes := services.RunCisAwsChecks()
	c.JSON(statusCode, apiRes)
}
