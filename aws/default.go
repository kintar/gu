package go_utils_aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"io"
)

// ReadParameterStoreString attempts to read a string value from AWS parameter store, performing decryption if necessary
// Returns the value and nil on success, or an empty string and error if there was a problem
func ReadParameterStoreString(cfg aws.Config, name string) (string, error) {
	ssmClient := ssm.NewFromConfig(cfg)
	getParam := &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: aws.Bool(true),
	}
	paramResults, err := ssmClient.GetParameter(context.TODO(), getParam)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve SSM parameter '%s': %w", name, err)
	}

	return *paramResults.Parameter.Value, nil
}

// WriteParameterStoreString does what's on the tin
func WriteParameterStoreString(cfg aws.Config, name string, value string, isSecure bool) error {
	ssmClient := ssm.NewFromConfig(cfg)
	paramType := types.ParameterTypeString
	if isSecure {
		paramType = types.ParameterTypeSecureString
	}
	putParam := &ssm.PutParameterInput{
		Name:      &name,
		Value:     &value,
		Type:      paramType,
		Overwrite: aws.Bool(true),
	}
	_, err := ssmClient.PutParameter(context.TODO(), putParam)
	return err
}

// LoadAwsConfig is a light wrapper around config.LoadDefaultConfig
func LoadAwsConfig() (aws.Config, error) {
	return config.LoadDefaultConfig(context.Background())
}

// MustLoadAwsConfig is a wrapper around LoadAwsConfig that panics with a useful message if the AWS configuration fails
// to load.
func MustLoadAwsConfig() aws.Config {
	cfg, err := LoadAwsConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load default AWS configuration: %w", err))
	}
	return cfg
}

// WriteToS3 pushes the content of the supplied io.Reader to the specified bucket and key
func WriteToS3(awsCfg aws.Config, bucket string, key string, data io.Reader) error {
	client := s3.NewFromConfig(awsCfg)
	put := s3.PutObjectInput{
		Body:   data,
		Bucket: &bucket,
		Key:    &key,
	}

	_, err := client.PutObject(context.TODO(), &put)
	return err
}
