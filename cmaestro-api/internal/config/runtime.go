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
	KeyVal   KeyValRuntimeConfig
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

type KeyValRuntimeConfig struct {
	Endpoint string
	Username string
	Password string
}

// ResolveRuntime turns environment-variable names from Config into actual values.
func (c *StaticConfig) ResolveRuntime() (*RuntimeConfig, error) {
	if c == nil {
		return nil, errors.New("config is nil")
	}

	artifactComponent := c.Components.Artifact
	sqlComponent := c.Components.SQL
	keyvalComponent := c.Components.KeyVal

	// Artifact Database Env Configuration
	artifactEndpoint, err := requiredEnv(artifactComponent.EndpointEnvironmentVariable)
	if err != nil {
		return nil, err
	}

	artifactAccessKey, err := requiredEnv(artifactComponent.AccessKeyEnvironmentVariable)
	if err != nil {
		return nil, err
	}

	artifactSecretKey, err := requiredEnv(artifactComponent.SecretKeyEnvironmentVariable)
	if err != nil {
		return nil, err
	}

	// SQL Database Env Configuration
	sqlEndpoint, err := requiredEnv(sqlComponent.EndpointEnvironmentVariable)
	if err != nil {
		return nil, err
	}

	sqlUsername, err := requiredEnv(sqlComponent.UsernameEnvironmentVariable)
	if err != nil {
		return nil, err
	}

	sqlPassword, err := requiredEnv(sqlComponent.PasswordEnvironmentVariable)
	if err != nil {
		return nil, err
	}

	sqlPort, err := requiredEnv(sqlComponent.PortEnvironmentVariable)
	if err != nil {
		return nil, err
	}
	sqlPortInt, err := strconv.Atoi(sqlPort)
	if err != nil {
		return nil, err
	}

	// KeyVal Database Env Configuration
	keyValEndpoint, err := requiredEnv(keyvalComponent.EndpointEnvironemtnVariable)
	if err != nil {
		return nil, err
	}

	keyValUsername, err := requiredEnv(keyvalComponent.UserNameEnvironmentVariable)
	if err != nil {
		return nil, err
	}

	keyValPassword, err := requiredEnv(keyvalComponent.PasswordEnvironmentVariable)
	if err != nil {
		return nil, err
	}

	return &RuntimeConfig{
		Artifact: ArtifactRuntimeConfig{
			Endpoint:   artifactEndpoint,
			AccessKey:  artifactAccessKey,
			SecretKey:  artifactSecretKey,
			BucketName: artifactComponent.BucketName,
			RootPrefix: artifactComponent.RootPrefix,
		},
		SQL: SQLRuntimeConfig{
			Endpoint: sqlEndpoint,
			Username: sqlUsername,
			Password: sqlPassword,
			Database: sqlComponent.DatabaseName,
			Port:     sqlPortInt,
		},
		KeyVal: KeyValRuntimeConfig{
			Endpoint: keyValEndpoint,
			Username: keyValUsername,
			Password: keyValPassword,
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
