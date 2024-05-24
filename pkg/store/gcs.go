package store

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)


type GCSStorage struct {
	client *storage.Client
	bucketName string
}

func NewGCSStorage(bucketName string) (*GCSStorage, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	return &GCSStorage{client: client, bucketName: bucketName}, nil
}

func (s *GCSStorage) Close() error {
	return s.client.Close()
}

func (s *GCSStorage) ListObjects() ([]string, error) {
	var objects []string
	it := s.client.Bucket(s.bucketName).Objects(context.Background(), nil)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj.Name)
	}
	return objects, nil
}
