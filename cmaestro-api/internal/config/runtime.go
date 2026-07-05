package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RuntimeConfig struct {
	Artifact ArtifactRuntimeConfig
	SQL      SQLRuntimeConfig
}

type ArtifactRuntimeConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	RootPrefix string
}

type SQLRuntimeConfig struct {
	Endpoint string
	Username string
	Password string
	Database string
	Port     int
}

// ResolveRuntime turns environment-variable names from Config into actual values.
func (c *StaticConfig) ResolveRuntime() (RuntimeConfig, error) {
	if c == nil {
		return RuntimeConfig{}, errors.New("config is nil")
	}

	artifact := c.Components.Artifact
	sqlComponent := c.Components.SQL

	artifactEndpoint, err := requiredEnv(artifact.EndpointEnvironmentVariable)
	if err != nil {
		return RuntimeConfig{}, err
	}

	artifactAccessKey, err := requiredEnv(artifact.AccessKeyEnvironmentVariable)
	if err != nil {
		return RuntimeConfig{}, err
	}

	artifactSecretKey, err := requiredEnv(artifact.SecretKeyEnvironmentVariable)
	if err != nil {
		return RuntimeConfig{}, err
	}

	sqlEndpoint, err := requiredEnv(sqlComponent.EndpointEnvironmentVariable)
	if err != nil {
		return RuntimeConfig{}, err
	}

	sqlUsername, err := requiredEnv(sqlComponent.UsernameEnvironmentVariable)
	if err != nil {
		return RuntimeConfig{}, err
	}

	sqlPassword, err := requiredEnv(sqlComponent.PasswordEnvironmentVariable)
	if err != nil {
		return RuntimeConfig{}, err
	}

	sqlPort, err := requiredEnv(sqlComponent.PortEnvironmentVariable)
	if err != nil {
		return RuntimeConfig{}, err
	}
	sqlPortInt, err := strconv.Atoi(sqlPort)
	if err != nil {
		return RuntimeConfig{}, err
	}

	return RuntimeConfig{
		Artifact: ArtifactRuntimeConfig{
			Endpoint:   artifactEndpoint,
			AccessKey:  artifactAccessKey,
			SecretKey:  artifactSecretKey,
			BucketName: artifact.BucketName,
			RootPrefix: artifact.RootPrefix,
		},
		SQL: SQLRuntimeConfig{
			Endpoint: sqlEndpoint,
			Username: sqlUsername,
			Password: sqlPassword,
			Database: sqlComponent.DatabaseName,
			Port:     sqlPortInt,
		},
	}, nil
}

func requiredEnv(name string) (string, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return "", errors.New("environment variable name is empty")
	}

	value, exists := os.LookupEnv(name)
	if !exists || strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("required environment variable %q is not set", name)
	}

	return value, nil
}
