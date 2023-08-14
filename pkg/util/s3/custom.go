package s3

import "github.com/aws/aws-sdk-go/aws"

const (
	defaultExpireTime = 10 // minutes

)

// custom variables
var (
	metadata = map[string]*string{
		"x-amz-storage-class": aws.String("STANDARD_IA"),
	}
)
