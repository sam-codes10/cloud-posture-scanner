# Cloud Posture Scanner

A comprehensive cloud security scanner written in Go. This application uses various AWS APIs, including CloudTrail, IAM, S3, and EC2, to monitor AWS resources and detect misconfigurations.

## Features

- **On-Demand Scanning**: Scans AWS resources in real-time when triggered via API endpoints.
- **Multi-Resource Support**:
  - **EC2 Instances**: Monitors instance status, type, IP addresses, and region.
  - **S3 Buckets**: Monitors public accessibility, encryption status, and access policies.
- **CIS Benchmarks**: Automatically runs CIS (Center for Internet Security) checks on detected resources.
- **Automatic Remediation**: Detected misconfigurations are automatically logged to a specified S3 bucket for review.
- **RESTful API**: Exposes endpoints to retrieve resource lists and CIS scan results.
- **Swagger Documentation**: Interactive API documentation generated using Swagger.

## Prerequisites

- **Go**: Version 1.20 or higher.
- **AWS Account**: Configured with appropriate IAM permissions for the scanner (EC2, S3, CloudTrail, IAM).
- **AWS Credentials**: Configured in your environment (e.g., `~/.aws/credentials` or environment variables).

## Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/sam-codes10/cloud-posture-scanner.git
    cd cloud-posture-scanner
    ```

2.  **Install dependencies**:
    ```bash
    go mod tidy
    ```

3.  **Configure Environment Variables** (optional but required if not using ~/.aws/credentials):
    Create a `.env` file in the root directory if needed:
    ```env
    AWS_ACCESS_KEY_ID=your-aws-access-key-id
    AWS_SECRET_ACCESS_KEY=your-aws-secret-access-key
    ```

## Usage

### Running the Application

Start the server using `go run`:

```bash
go run .
```

The server will start on port `8080`.

### API Endpoints

Once the server is running, you can interact with the following endpoints:

#### Get EC2 Instances
Retrieves a list of all EC2 instances in the configured AWS account.
- **Method**: `GET`
- **Endpoint**: `/api/instances`
- **Response**:
  ```json
  {
    "status": true,
    "message": "EC2 instances fetched successfully",
    "data": [
      {
        "id": "i-0123456789abcdef0",
        "name": "web-server",
        "type": "t3.medium",
        "region": "us-east-1a",
        "public_ip": "52.x.x.x",
        "private_ip": "172.31.x.x",
        "state": "running"
      }
    ]
  }
  ```

#### Get S3 Buckets
Retrieves a list of all S3 buckets in the configured AWS account.
- **Method**: `GET`
- **Endpoint**: `/api/buckets`
- **Response**:
  ```json
  {
    "status": true,
    "message": "S3 buckets fetched successfully",
    "data": [
      {
        "id": "my-bucket",
        "name": "my-bucket",
        "region": "us-east-1",
        "encryption_status": "AES256",
        "access_policy": "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":\"*\",\"Action\":\"s3:GetObject\",\"Resource\":\"arn:aws:s3:::my-bucket/*\"}]}"
      }
    ]
  }
  ```

#### Get CIS Results
Triggers an immediate CIS benchmark scan and retrieves the results. This endpoint also starts the background monitoring process.
- **Method**: `GET`
- **Endpoint**: `/api/cis-results`
- **Response**:
  ```json
  {
    "status": true,
    "message": "CIS AWS checks run successfully",
    "data": [
      {
        "name": "Check 1",
        "description": "S3 bucket shall not be publicly accessible",
        "associated_resource": "my-bucket",
        "status": "FAIL",
        "remediation": "Remove public access from S3 bucket"
      }
    ]
  }
  ```

#### Swagger Documentation / Frontend Output
Access interactive API documentation.
- **Endpoint**: `http://localhost:8080/swagger/cloud-posture-scanner/index.html#/`

## Output Format

### CIS Check Results

The scan generates detailed reports for each resource, including:
- **Check Name & Description**: Standardized CIS benchmark names.
- **Associated Resource**: The ID of the resource being checked.
- **Status**: `PASS`, `FAIL`.
- **Remediation**: Specific steps to resolve the misconfiguration.

### Output to S3

When a scan is triggered, the results are automatically uploaded to a dedicated S3 bucket.
- **Bucket**: Named `cis-check-reports`.
