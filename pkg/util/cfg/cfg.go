package cfgutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/caarlos0/env/v5"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

// Load loads configuration from local .env file
func Load(out interface{}, stage string) error {
	if err := PreloadLocalENV(stage); err != nil {
		return err
	}

	if err := env.Parse(out); err != nil {
		return err
	}

	return nil
}

// LoadWithAPS loads configuration from local .env file and AWS Parameter Store as well
func LoadWithAPS(out interface{}, appName, stage string) error {
	if appName != "" && stage != "development" {
		if err := PreloadAWSENV(appName, stage, nil); err != nil {
			return err
		}
	}
	return Load(out, stage)
}

// PreloadLocalENV reads .env* files and sets the values to os ENV
func PreloadLocalENV(stage string) error {
	basePath := ""
	if stage == "test" {
		basePath = "testdata/"
	}
	// // local config per stage
	// if stage != "" {
	// 	godotenv.Load(basePath + ".env." + stage + ".local")
	// }

	// local config
	godotenv.Load(basePath + ".env.local")

	// // per stage config
	// if stage != "" {
	// 	godotenv.Load(basePath + ".env." + stage)
	// }

	// default config
	return godotenv.Load(basePath + ".env")
}

// PreloadAWSENV reads AWS Parameter Store and set them to os ENV
func PreloadAWSENV(appName, stage string, nextToken *string) error {
	basePath := "/" + appName + "/"
	if stage != "" {
		basePath += stage + "/"
	}

	svc := ssm.New(session.New())
	resp, err := svc.GetParametersByPath(&ssm.GetParametersByPathInput{
		NextToken:      nextToken,
		Path:           aws.String(basePath),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	for _, param := range resp.Parameters {
		paramName := strings.Replace(*param.Name, basePath, "", 1)
		os.Setenv(strings.ToUpper(paramName), *param.Value)
	}

	if resp.NextToken != nil {
		return PreloadAWSENV(appName, stage, resp.NextToken)
	}

	return nil
}
