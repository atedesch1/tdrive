package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	configName = ".tdriveconfig"
	configType = "yaml"

	GCSStorageProvider StorageProvider = "GCS"
	// S3StorageProvider    StorageProvider = "s3"
	// LocalStorageProvider StorageProvider = "local"
)

var StorageProviders = []StorageProvider{GCSStorageProvider}

type StorageProvider string

type Config struct {
	StorageProvider StorageProvider `yaml:"storage_provider"`
	StorageConfig   StorageConfig   `yaml:"storage_config"`
	DownloadPath    string          `yaml:"download_path"`
}

type StorageConfig struct {
	GCSStorageConfig *GCSStorageConfig `yaml:"gcs"`
}

type GCSStorageConfig struct {
	BucketName string `yaml:"bucket_name"`
}

func GetConfigPath() string {
	configDir, ok := os.LookupEnv("TDRIVE_CONFIG_DIR")
	if !ok {
		configDir = os.Getenv("HOME")
	}

	return configDir
}

func WriteConfig(config *Config) error {
	path := GetConfigPath()
	configFileName := path + "/" + configName + "." + configType

	file, err := os.Create(configFileName)
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}
	defer file.Close()

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("unable to marshal config to YAML: %v", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("unable to write data to file: %v", err)
	}

	return nil
}

func LoadConfig(path string) (*Config, error) {
	var config Config

	configFileName := path + "/" + configName + "." + configType

	file, err := os.Open(configFileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %v", err)
	}

	return &config, nil
}
