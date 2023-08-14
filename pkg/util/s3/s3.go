package s3

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
)

// PreparePresignedURL to get object presign then return a PUT url
func (s *Service) PreparePresignedURL(key string, expireTime int, presignedURL *string) error {
	req, _ := s.s3.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.cfg.BucketName),
		Key:    aws.String(key),
	})

	return processPresignedURL(req, expireTime, presignedURL)
}

// GetPresignedURL to get object presigned || get object with signed URL
func (s *Service) GetPresignedURL(key string, expireTime int, presignedURL *string) error {
	req, _ := s.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.cfg.BucketName),
		Key:    aws.String(key),
	})

	return processPresignedURL(req, expireTime, presignedURL)
}

func processPresignedURL(request *request.Request, expireTime int, presignedURL *string) error {
	if expireTime == 0 {
		expireTime = defaultExpireTime
	}

	url, err := request.Presign(time.Duration(expireTime) * time.Minute)
	if err != nil {
		return err
	}

	*presignedURL = url

	return nil
}
