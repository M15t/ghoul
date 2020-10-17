# AWS Simple Notification Service (SNS) wrapper

This package contains wrapper functions for SNS service to simplify the normal usages

## Example

```go
import (
  "fmt"
	snsutil "github.com/M15t/saas-api/internal/util/sns"
)

appArn := "put application arn here"
deviceToken := "put device token here"

// initialize the service with default configuration, using env vars for credentials
snsSvc := snsutil.New()

// register a new device
resp, err := snsSvc.RegisterDevice(appArn, deviceToken)
if err != nil {
	panic(fmt.Errorf("error creating sns endpoint: %+v", err))
}
// grab the arn
deviceArn := *resp.EndpointArn
fmt.Printf("generated sns arn: %s\n", deviceArn)

// sample data for push notification
datapayload := map[string]interface{}{
	"ghoul": map[string]interface{}{
		"type":   "message type",
		"field1": "blah blah",
		"field2": "ok",
	},
}

// send the push notification to Android devices
output, err := snsSvc.SendToAndroid(deviceArn, snsutil.FCMPayload{Data: datapayload})
// or iOS devices
output, err = snsSvc.SendToIOS(deviceArn, snsutil.APNSPayload{Data: datapayload})

// or to all devices without knowing the OS
output, err = snsSvc.SendToDevice(deviceArn, snsutil.Message{
	APNS:        &snsutil.APNSPayload{Data: datapayload},
	APNSSandbox: &snsutil.APNSPayload{Data: datapayload},
	FCM:         &snsutil.FCMPayload{Data: datapayload},
})
if err != nil {
	panic(fmt.Errorf("error publishing message: %+v", err))
}
fmt.Printf("publish output: %+v\n", output)

// all done, to remove the device from SNS
_, err = snsSvc.DeregisterDevice(deviceArn)
if err != nil {
	panic(fmt.Errorf("error deregistering device : %+v", err))
}
```
