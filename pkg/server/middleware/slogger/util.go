package slogger

import (
	"encoding/json"
	"fmt"
	"strings"
)

func censorSecerts(uri string, body map[string]interface{}, secrets []string) map[string]interface{} {
	for key, val := range body {
		found := false
		lowerkey := strings.ToLower(key)
		for _, secretKey := range secrets {
			if secretKey == lowerkey || strings.Contains(lowerkey, secretKey) {
				found = true
				break
			}
		}
		if found {
			body[key] = "***"
			continue
		}

		switch v := val.(type) {
		case map[string]interface{}:
			body[key] = censorSecerts(uri, v, secrets)
			continue
		}
	}

	return body
}

func prettyJSON(jsonData []byte) map[string]interface{} {
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	return data
}
