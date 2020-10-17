package snsutil

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"strings"
)

// CreateApplication creates new platform application on SNS service
func (s *Service) CreateApplication(name string, platform string, attr map[string]*string) (string, error) {
	input := &sns.CreatePlatformApplicationInput{
		Name:       aws.String(name),
		Platform:   aws.String(platform),
		Attributes: attr,
	}
	output, err := s.sns.CreatePlatformApplication(input)
	if err != nil {
		return "", err
	}
	return *output.PlatformApplicationArn, nil
}

// CreateAPNSApplication creates new Apple APNS platform application on SNS service
func (s *Service) CreateAPNSApplication(name string, cert string, key string, sandbox bool) (string, error) {
	attr := map[string]*string{
		"PlatformPrincipal":  aws.String(cert),
		"PlatformCredential": aws.String(key),
	}
	platform := "APNS"
	if sandbox {
		platform = "APNS_SANDBOX"
	}
	return s.CreateApplication(name, platform, attr)
}

// CreateFCMApplication creates new FCM platform application on SNS service
func (s *Service) CreateFCMApplication(name string, key string) (string, error) {
	attr := map[string]*string{
		"PlatformCredential": aws.String(key),
	}
	return s.CreateApplication(name, "GCM", attr)
}

// DeleteApplication deletes the platform application from SNS service
func (s *Service) DeleteApplication(appArn string) error {
	_, err := s.sns.DeletePlatformApplication(&sns.DeletePlatformApplicationInput{
		PlatformApplicationArn: aws.String(appArn),
	})
	return err
}

// UpdateApplication updates the platform application attributes
func (s *Service) UpdateApplication(appArn string, attr map[string]*string) error {
	input := &sns.SetPlatformApplicationAttributesInput{
		PlatformApplicationArn: aws.String(appArn),
		Attributes:             attr,
	}
	_, err := s.sns.SetPlatformApplicationAttributes(input)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAPNSApplication updates Apple APNS platform application attributes
func (s *Service) UpdateAPNSApplication(arn string, cert string, key string, sandbox bool) (string, error) {
	attr := map[string]*string{
		"PlatformPrincipal":  aws.String(cert),
		"PlatformCredential": aws.String(key),
	}
	appIsSandBox := strings.Contains(arn, "APNS_SANDBOX")
	if appIsSandBox != sandbox {
		// cannot switch the sandbox flag, so delete & recreate the app
		tmp := strings.Split(arn, "/")
		newarn, err := s.CreateAPNSApplication(tmp[2], cert, key, sandbox)
		if err != nil {
			return "", err
		}

		if err := s.DeleteApplication(arn); err != nil {
			return "", err
		}

		return newarn, nil
	}

	return "", s.UpdateApplication(arn, attr)
}

// UpdateFCMApplication updates FCM platform application attributes
func (s *Service) UpdateFCMApplication(arn string, key string) error {
	attr := map[string]*string{
		"PlatformCredential": aws.String(key),
	}
	return s.UpdateApplication(arn, attr)
}

// RegisterDevice registers a device with SNS application
func (s *Service) RegisterDevice(appArn, deviceToken string) (string, error) {
	output, err := s.sns.CreatePlatformEndpoint(&sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(appArn),
		Token:                  aws.String(deviceToken),
	})
	if err != nil {
		return "", err
	}
	return *output.EndpointArn, nil
}

// DeregisterDevice remove a device from SNS application
func (s *Service) DeregisterDevice(endpointArn string) error {
	_, err := s.sns.DeleteEndpoint(&sns.DeleteEndpointInput{EndpointArn: aws.String(endpointArn)})
	return err
}

// SendToDevice sends push notification to a device. The "msg" should contains correct payload for FCM or APNS plarform, depends on the "target" OS
func (s *Service) SendToDevice(target string, msg Message) (string, error) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	output, err := s.sns.Publish(&sns.PublishInput{
		Message:          aws.String(string(jsonMsg)),
		MessageStructure: aws.String("json"),
		TargetArn:        aws.String(target),
	})
	if err != nil {
		return "", err
	}
	return *output.MessageId, nil
}

// SendToAndroid sends push notification to an Android device using FCM platform
func (s *Service) SendToAndroid(target string, payload FCMPayload) (string, error) {
	msg := Message{
		FCM: &payload,
	}
	return s.SendToDevice(target, msg)
}

// SendToIOS sends push notification to an IOS device using APNS platform
func (s *Service) SendToIOS(target string, payload APNSPayload) (string, error) {
	msg := Message{
		APNS:        &payload,
		APNSSandbox: &payload,
	}
	return s.SendToDevice(target, msg)
}

// SendToTopic sends push notification to a topic. The "msg" should contains payload for both FCM and APNS plarforms
func (s *Service) SendToTopic(topic string, msg Message) (string, error) {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	output, err := s.sns.Publish(&sns.PublishInput{
		Message:          aws.String(string(jsonMsg)),
		MessageStructure: aws.String("json"),
		TopicArn:         aws.String(topic),
	})

	if err != nil {
		return "", err
	}
	return *output.MessageId, nil
}
