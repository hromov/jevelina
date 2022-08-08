package gcloud

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/hromov/jevelina/domain/misc/files"
	"github.com/vincent-petithory/dataurl"
)

type service struct {
	bucketName string
	c          *storage.Client
}

func NewService(ctx context.Context, bucketName string) (*service, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &service{bucketName: bucketName, c: client}, nil
}

func (fs *service) Delete(ctx context.Context, url string) error {
	if err := fs.c.Bucket(fs.bucketName).Object(url).Delete(ctx); err != nil {
		return fmt.Errorf("Can't delete file: %s, error: %s", url, err.Error())
	}
	return nil
}

func (fs *service) Upload(ctx context.Context, req files.FileAddReq) (files.FileCreateReq, error) {
	fileName := uuid.New().String()

	wc := fs.c.Bucket(fs.bucketName).Object(fileName).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.
	wc.ContentType = req.Type
	wc.Metadata = map[string]string{
		"MIMEType": req.Type,
		"FileName": req.Name,
	}

	dataURL, _ := dataurl.DecodeString(req.Value)
	if _, err := wc.Write(dataURL.Data); err != nil {
		return files.FileCreateReq{}, err
	}
	if err := wc.Close(); err != nil {
		return files.FileCreateReq{}, err
	}
	file := files.FileCreateReq{
		ParentID: req.Parent,
		Name:     req.Name,
		URL:      fileName,
	}

	return file, nil
}

func (fs *service) PresignUrl(url string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}
	return fs.c.Bucket(fs.bucketName).SignedURL(url, opts)
}
