package slogger

import (
	"encoding/json"
	"fmt"
	"strings"
)

func secureBody(body interface{}, sensitiveKeys []string) {
	switch v := body.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if containsSensitiveKey(key, sensitiveKeys) {
				v[key] = "******"
			} else {
				secureBody(value, sensitiveKeys)
			}
		}
	case []interface{}:
		for _, value := range v {
			secureBody(value, sensitiveKeys)
		}
	}
}

func containsSensitiveKey(key string, sensitiveKeys []string) bool {
	for _, sensitiveKey := range sensitiveKeys {
		if strings.Contains(strings.ToLower(key), strings.ToLower(sensitiveKey)) {
			return true
		}
	}
	return false
}

func prettyJSON(jsonData []byte) map[string]interface{} {
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	secureBody(data, sensitiveKeys)

	return data
}
