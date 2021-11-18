package email

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Config represents the configuration
type Config struct {
	Sender string
	Region string
	WebURL string
}

// New initializes SNS service with default config
func New(cfg Config) *Email {
	return &Email{
		ses: ses.New(session.New(&aws.Config{Region: aws.String(cfg.Region)})),
		cfg: cfg,
	}
}

// Email represents the sesutil service
type Email struct {
	cfg Config
	ses *ses.SES
}
