package services

import (
	awshelpers "cloud-posture-scanner/aws-helpers"
	"cloud-posture-scanner/models"
	"net/http"
)

func GetEC2Instances() (int, models.APIRes) {
	instances, err := awshelpers.EnlistAllEC2Instances()
	if err != nil {
		return http.StatusInternalServerError, models.APIRes{
			Status:  false,
			Message: "Failed to fetch EC2 instances",
			Data:    err,
		}
	}
	if len(instances) == 0 {
		return http.StatusOK, models.APIRes{
			Status:  true,
			Message: "No EC2 instances found",
			Data:    instances,
		}
	}
	return http.StatusOK, models.APIRes{
		Status:  true,
		Message: "EC2 instances fetched successfully",
		Data:    instances,
	}
}

func GetS3Buckets() (int, models.APIRes) {
	buckets, err := awshelpers.EnlistAllS3Buckets()
	if err != nil {
		return http.StatusInternalServerError, models.APIRes{
			Status:  false,
			Message: "Failed to fetch S3 buckets",
			Data:    err,
		}
	}
	if len(buckets) == 0 {
		return http.StatusOK, models.APIRes{
			Status:  true,
			Message: "No S3 buckets found",
			Data:    buckets,
		}
	}
	return http.StatusOK, models.APIRes{
		Status:  true,
		Message: "S3 buckets fetched successfully",
		Data:    buckets,
	}
}

func RunCisAwsChecks() (int, models.APIRes) {
	buckets, err := awshelpers.EnlistAllS3Buckets()
	if err != nil {
		return http.StatusInternalServerError, models.APIRes{
			Status:  false,
			Message: "Failed to fetch S3 buckets",
			Data:    err,
		}
	}

	instances, err := awshelpers.EnlistAllEC2Instances()
	if err != nil {
		return http.StatusInternalServerError, models.APIRes{
			Status:  false,
			Message: "Failed to fetch EC2 instances",
			Data:    err,
		}
	}

	cisChecks := awshelpers.RunCisAwsChecks(buckets, instances)
	if len(cisChecks) == 0 {
		return http.StatusOK, models.APIRes{
			Status:  true,
			Message: "No CIS AWS checks found",
			Data:    cisChecks,
		}
	}

	go awshelpers.PublishCisChecksResultsToS3Bucket(cisChecks)

	return http.StatusOK, models.APIRes{
		Status:  true,
		Message: "CIS AWS checks run successfully",
		Data:    cisChecks,
	}
}
