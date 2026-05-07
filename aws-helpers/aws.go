package awshelpers

import (
	"bytes"
	"cloud-posture-scanner/models"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var cfg aws.Config

func InitS3Config() {
	var err error
	cfg, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}
	log.Println("AWS Account ID:", cfg.Credentials)
}

func EnlistAllEC2Instances() ([]models.EC2Instance, error) {
	client := ec2.NewFromConfig(cfg)
	instances, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Println("failed to describe instances, %v", err)
		return nil, err
	}

	var ec2Instances []models.EC2Instance

	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			ec2Instances = append(ec2Instances, models.EC2Instance{
				ID:        aws.ToString(instance.InstanceId),
				Name:      aws.ToString(instance.Tags[0].Value),
				Type:      string(instance.InstanceType),
				Region:    aws.ToString(instance.Placement.AvailabilityZone),
				PublicIP:  aws.ToString(instance.PublicIpAddress),
				PrivateIP: aws.ToString(instance.PrivateIpAddress),
				State:     string(instance.State.Name),
			})
		}
	}
	return ec2Instances, nil
}

func EnlistAllS3Buckets() ([]models.S3Bucket, error) {
	client := s3.NewFromConfig(cfg)
	buckets, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Println("failed to list buckets, %v", err)
		return nil, err
	}

	var s3Buckets []models.S3Bucket

	for _, bucket := range buckets.Buckets {
		region := ""
		location, err := client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
			Bucket: bucket.Name,
		})
		if err == nil && location != nil && string(location.LocationConstraint) != "" {
			region = string(location.LocationConstraint)
		} else if err != nil {
			log.Println("failed to get bucket location of name %s, %v", *bucket.Name, err)
		}

		policyStr := ""
		policy, err := client.GetBucketPolicy(context.TODO(), &s3.GetBucketPolicyInput{
			Bucket: bucket.Name,
		})
		if err == nil && policy != nil {
			policyStr = *policy.Policy
		} else if err != nil {
			log.Println("failed to get bucket policy of name %s, %v", *bucket.Name, err)
		}

		encryption, err := client.GetBucketEncryption(context.TODO(), &s3.GetBucketEncryptionInput{
			Bucket: bucket.Name,
		})
		encryptionStr := ""
		if err == nil && encryption != nil && len(encryption.ServerSideEncryptionConfiguration.Rules) > 0 {
			encryptionStr = string(encryption.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)
		} else if err != nil {
			log.Println("failed to get bucket encryption of name %s, %v", *bucket.Name, err)
		}

		s3Buckets = append(s3Buckets, models.S3Bucket{
			ID:               aws.ToString(bucket.Name),
			Name:             aws.ToString(bucket.Name),
			Region:           region,
			EncryptionStatus: encryptionStr,
			AccessPolicy:     policyStr,
		})
	}
	return s3Buckets, nil
}

func RunCisAwsChecks(buckets []models.S3Bucket, instances []models.EC2Instance) []models.CISCheck {
	var cisChecks []models.CISCheck
	for _, bucket := range buckets {
		// check-1 s3 bucket shall not be publicly accessible
		if bucket.AccessPolicy != "" {
			cisChecks = append(cisChecks, models.CISCheck{
				Name:               "Check 1",
				Description:        "S3 bucket shall not be publicly accessible",
				AssociatedResource: bucket.ID,
				Status:             "FAIL",
				Remediation:        "Remove public access from S3 bucket",
			})
		} else {
			cisChecks = append(cisChecks, models.CISCheck{
				Name:               "Check 1",
				Description:        "S3 bucket shall not be publicly accessible",
				AssociatedResource: bucket.ID,
				Status:             "PASS",
				Remediation:        "All good",
			})
		}

		// check-2 s3 bucket shall have encryption enabled
		if bucket.EncryptionStatus == "" {
			cisChecks = append(cisChecks, models.CISCheck{
				Name:               "Check 2",
				Description:        "S3 bucket shall have encryption enabled",
				AssociatedResource: bucket.ID,
				Status:             "FAIL",
				Remediation:        "Enable encryption on S3 bucket",
			})
		} else {
			cisChecks = append(cisChecks, models.CISCheck{
				Name:               "Check 2",
				Description:        "S3 bucket shall have encryption enabled",
				AssociatedResource: bucket.ID,
				Status:             "PASS",
				Remediation:        "All good",
			})
		}
	}

	// check 3 IAM root account should have MFA enabled
	iamClient := iam.NewFromConfig(cfg)
	summary, err := iamClient.GetAccountSummary(context.TODO(), &iam.GetAccountSummaryInput{})
	if err == nil && summary != nil {
		mfaEnabled := false
		if val, ok := summary.SummaryMap["AccountMFAEnabled"]; ok && val == 1 {
			mfaEnabled = true
		}

		status := "FAIL"
		remediation := "Enable MFA for the root account"
		if mfaEnabled {
			status = "PASS"
			remediation = "All good"
		}

		cisChecks = append(cisChecks, models.CISCheck{
			Name:               "Check 3",
			Description:        "IAM root account should have MFA enabled",
			AssociatedResource: "AWS Account Root",
			Status:             status,
			Remediation:        remediation,
		})
	} else if err != nil {
		log.Println("failed to get IAM account summary:", err)
	}

	// check 4 CloudTrail should be enabled
	ctClient := cloudtrail.NewFromConfig(cfg)
	trails, err := ctClient.DescribeTrails(context.TODO(), &cloudtrail.DescribeTrailsInput{})
	if err == nil {
		enabled := false
		if trails != nil && len(trails.TrailList) > 0 {
			enabled = true
		}

		status := "FAIL"
		remediation := "Enable CloudTrail logging for the AWS account"
		if enabled {
			status = "PASS"
			remediation = "All good"
		}

		cisChecks = append(cisChecks, models.CISCheck{
			Name:               "Check 4",
			Description:        "CloudTrail should be enabled",
			AssociatedResource: "AWS Account",
			Status:             status,
			Remediation:        remediation,
		})
	} else {
		log.Println("failed to get CloudTrail info:", err)
	}

	// check 5 Security groups should not open to 0.0.0.0/0 for SSH or RDP
	for _, instance := range instances {
		if instance.PublicIP != "" {
			cisChecks = append(cisChecks, models.CISCheck{
				Name:               "Check 5",
				Description:        "Security groups should not open to 0.0.0.0/0 for SSH or RDP",
				AssociatedResource: instance.ID,
				Status:             "FAIL",
				Remediation:        "Remove public access from EC2 instance",
			})
		} else {
			cisChecks = append(cisChecks, models.CISCheck{
				Name:               "Check 5",
				Description:        "Security groups should not open to 0.0.0.0/0 for SSH or RDP",
				AssociatedResource: instance.ID,
				Status:             "PASS",
				Remediation:        "All good",
			})
		}
	}

	return cisChecks
}

func PublishCisChecksResultsToS3Bucket(cisChecks []models.CISCheck) error {
	file, err := os.Create("cis-results.json")
	if err != nil {
		log.Println("failed to create file:", err)
		return err
	}
	defer file.Close()

	byteData, err := json.MarshalIndent(cisChecks, "", "  ")
	if err != nil {
		log.Println("failed to marshal json:", err)
		return err
	}

	file.Write(byteData)

	client := s3.NewFromConfig(cfg)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("cis-check-reports"),
		Key:    aws.String("cis-results.json"),
		Body:   bytes.NewReader(byteData),
	})
	if err != nil {
		log.Println("failed to put object:", err)
		return err
	}

	return nil
}
