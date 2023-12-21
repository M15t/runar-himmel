package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

// Options for loading
type Options struct {
	// Support all options from env package
	env.Options
	// Default to .env at current path if empty
	DotenvPath string
	// Whether to load from AWS Parameter Store. Slashes are required: /path/prefix/
	AwsParameterStorePath string
}

// DefaultOptions for LoadAll function
var DefaultOptions = Options{
	DotenvPath:            ".env",
	AwsParameterStorePath: os.Getenv("AWS_PARAMETER_STORE_PATH"),
}

// Load loads configurations from all sources with default options
func Load(dest interface{}) error {
	return LoadWithOptions(dest, DefaultOptions)
}

// LoadWithOptions fills configuration from various sources into destination struct.
//
// Overwriting order: .env file < environment variables < AWS Parameter Store
func LoadWithOptions(dest interface{}, opts Options) error {
	envvars := map[string]string{}

	// overwrite with values from .env.local file
	if !IsLambda() {
		if err := LoadFromDotenv(envvars, ".env.local"); err != nil {
			return fmt.Errorf("error preloading dotenv local: %w", err)
		}
	}

	// overwrite with values from .env file
	if opts.DotenvPath != "" {
		if err := LoadFromDotenv(envvars, opts.DotenvPath); err != nil {
			return fmt.Errorf("error preloading dotenv: %w", err)
		}
	}

	// overwrite with OS environment variables
	if opts.Environment == nil {
		LoadFromEnv(envvars)
	} else {
		// use values from options if provided
		for k, v := range opts.Environment {
			envvars[k] = v
		}
	}

	// overwrite with params from AWS Parameter Store
	if opts.AwsParameterStorePath != "" {
		if err := LoadFromAwsParameterStore(envvars, opts.AwsParameterStorePath, nil); err != nil {
			return fmt.Errorf("error preloading awsenv: %w", err)
		}
	}

	// put everything together
	opts.Options.Environment = envvars
	if err := env.ParseWithOptions(dest, opts.Options); err != nil {
		return fmt.Errorf("error parsing env: %w", err)
	}

	return nil
}

// LoadFromDotenv parses .env file and sets the values to destination map
func LoadFromDotenv(dest map[string]string, files ...string) error {
	envmap, err := godotenv.Read(files...)
	// ignore file not found error
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for k, v := range envmap {
		dest[k] = v
	}
	return nil
}

// LoadFromEnv parses OS environment variables and set values to destination map
func LoadFromEnv(dest map[string]string) {
	envs := osEnviron()
	for _, e := range envs {
		p := strings.SplitN(e, "=", 2)
		dest[p[0]] = p[1]
	}
}

// LoadFromAwsParameterStore gets params from AWS Parameter Store and sets them to destination map
func LoadFromAwsParameterStore(dest map[string]string, path string, nextToken *string) error {
	if path == "" {
		return fmt.Errorf("path prefix is required")
	}
	if path[0] != '/' {
		path = "/" + path
	}
	if path[len(path)-1] != '/' {
		path = path + "/"
	}

	svc := getSsmService()
	resp, err := svc.GetParametersByPath(&ssm.GetParametersByPathInput{
		NextToken:      nextToken,
		Path:           aws.String(path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	for _, param := range resp.Parameters {
		paramName := strings.Replace(*param.Name, path, "", 1)
		dest[strings.ToUpper(paramName)] = *param.Value
	}

	if resp.NextToken != nil {
		return LoadFromAwsParameterStore(dest, path, resp.NextToken)
	}

	return nil
}

// For testing purpose, to be overwritten in test
var osEnviron = os.Environ

// For testing purpose, to be overwritten in test
var getSsmService = func() SsmSvc {
	return ssm.New(session.Must(session.NewSession()))
}

// SsmSvc to generate the mock functions
// mockgen -destination config/mock_ssm_test.go -package config github.com/alextanhongpin/go-config/config SsmSvc
type SsmSvc interface {
	GetParametersByPath(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error)
}

// IsLambda checks whether the code is running on lambda using predefined Lambda environment variables
func IsLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" && os.Getenv("AWS_LAMBDA_FUNCTION_VERSION") != ""
}
