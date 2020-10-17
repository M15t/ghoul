package snsutil

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// New initializes SNS service with default config
func New() *Service {
	return &Service{sns: sns.New(session.New())}
}

// Service represents the snsutil service
type Service struct {
	sns *sns.SNS
}

// FCMNotification represents the user-visible of the notification for FCM platform
type FCMNotification struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Body     string `json:"body,omitempty"`
	Sound    string `json:"sound,omitempty"`
	Badge    *int   `json:"badge,omitempty"`
}

// FCMPayload represents the FCM message payload
type FCMPayload struct {
	Data         map[string]interface{}
	Notification *FCMNotification
	HighPriority bool
}

// MarshalJSON modifies json marshal output for FCMPayload
func (fp FCMPayload) MarshalJSON() ([]byte, error) {
	msg := make(map[string]interface{})
	if fp.Data != nil {
		msg["data"] = fp.Data
	}
	if fp.Notification != nil {
		msg["notification"] = fp.Notification
	}
	if fp.HighPriority {
		msg["priority"] = "high"
	}
	return json.Marshal(msg)
}

// APNSAlert represents the APNS alert object
type APNSAlert struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Body     string `json:"body,omitempty"`
}

// APNSNotification represents the user-visible of the notification for APNS platform
type APNSNotification struct {
	Alert            *APNSAlert `json:"alert,omitempty"`
	Sound            string     `json:"sound,omitempty"`
	Badge            *int       `json:"badge,omitempty"`
	ContentAvailable int        `json:"content-available,omitempty"`
	Category         *string    `json:"category,omitempty"`
	MutableContent   *int       `json:"mutable-content,omitempty"`
}

// APNSPayload represents the APNS message payload
type APNSPayload struct {
	Data         map[string]interface{}
	Notification *APNSNotification
	HighPriority bool
}

// MarshalJSON modifies json marshal output for APNSPayload
func (ap APNSPayload) MarshalJSON() ([]byte, error) {
	msg := make(map[string]interface{})
	if ap.Data != nil {
		for k, v := range ap.Data {
			msg[k] = v
		}
	}
	if ap.Notification == nil {
		msg["aps"] = &APNSNotification{ContentAvailable: 1}
	} else {
		msg["aps"] = ap.Notification
	}
	if ap.HighPriority {
		msg["priority"] = "high"
	}
	return json.Marshal(msg)
}

// Message represents the message structure for SNS message
type Message struct {
	APNS        *APNSPayload
	APNSSandbox *APNSPayload
	FCM         *FCMPayload
}

// MarshalJSON modifies json marshal output for Message
func (m Message) MarshalJSON() ([]byte, error) {
	msg := map[string]string{}
	if m.APNS != nil {
		apns, err := json.Marshal(m.APNS)
		if err != nil {
			return []byte{}, err
		}
		msg["APNS"] = string(apns)
	}
	if m.APNSSandbox != nil {
		apns, err := json.Marshal(m.APNSSandbox)
		if err != nil {
			return []byte{}, err
		}
		msg["APNS_SANDBOX"] = string(apns)
	}
	if m.FCM != nil {
		apns, err := json.Marshal(m.FCM)
		if err != nil {
			return []byte{}, err
		}
		msg["GCM"] = string(apns)
	}
	return json.Marshal(msg)
}
