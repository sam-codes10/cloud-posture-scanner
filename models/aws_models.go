package models

type EC2Instance struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Region    string `json:"region"`
	PublicIP  string `json:"public_ip"`
	PrivateIP string `json:"private_ip"`
	State     string `json:"state"`
}

type S3Bucket struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Region           string `json:"region"`
	EncryptionStatus string `json:"encryption_status"`
	AccessPolicy     string `json:"access_policy"`
}

type CISCheck struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	AssociatedResource string `json:"associated_resource"`
	Status             string `json:"status"`
	Remediation        string `json:"remediation"`
}
