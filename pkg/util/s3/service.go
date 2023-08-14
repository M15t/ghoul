package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Config represents the configuration
type Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	Debug           bool
}

// Service represents the s3 service
type Service struct {
	cfg *Config
	s3  *s3.S3
}

// New initializes s3 service with default config
func New(cfg *Config) *Service {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	})
	if err != nil {
		panic("Initialize S3 service failed")
	}

	return &Service{
		cfg: cfg,
		s3:  s3.New(sess),
	}
}
