package store

import (
	"fmt"

	"github.com/atedesch1/tdrive/pkg/config"
)

type Store interface {
	ListObjects(prefix string) ([]string, error)
	DownloadObject(objectPath string, targetDir string) error
}

func NewStore(cfg *config.Config) (Store, error) {
	switch cfg.StorageProvider {
	case config.GCSStorageProvider:
		if cfg.StorageConfig.GCSStorageConfig == nil {
			return nil, fmt.Errorf("missing GCS storage config")
		}
		return NewGCSStorage(cfg.StorageConfig.GCSStorageConfig.BucketName)
	default:
		return nil, fmt.Errorf("unsupported storage provider")
	}
}
