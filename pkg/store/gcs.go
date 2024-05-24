package store

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GCSStorage struct {
	client     *storage.Client
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

func (s *GCSStorage) ListObjects(prefix string) ([]string, error) {
	var objects []string
	it := s.client.
		Bucket(s.bucketName).
		Objects(
			context.Background(),
			&storage.Query{
				Prefix:    prefix,
				Delimiter: "/",
			},
		)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj.Prefix+obj.Name)
	}
	return objects, nil
}

func (s *GCSStorage) DownloadObject(objectPath string, targetDir string) error {
	if objectPath[len(objectPath)-1] == '/' {
		return fmt.Errorf("cannot download directory")
	}

	obj := s.client.Bucket(s.bucketName).Object(objectPath)
	rc, err := obj.NewReader(context.Background())
	if err != nil {
		return err
	}
	defer rc.Close()

	index := strings.LastIndex(objectPath, "/")
	var objectName string
	if index == -1 {
		objectName = objectPath
	} else {
		objectName = objectPath[index+1:]
	}

	outFile, err := os.Create(targetDir + "/" + objectName)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, rc); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	return nil
}
