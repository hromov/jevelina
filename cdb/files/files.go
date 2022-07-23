package files

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/hromov/jevelina/cdb/models"
	"github.com/vincent-petithory/dataurl"
	"gorm.io/gorm"
)

type FilesService struct {
	BucketName string
	*gorm.DB
}

func (fs *FilesService) Delete(ID uint64) error {
	var file models.File
	if err := fs.DB.First(&file, ID).Error; err != nil {
		return err
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	if err = client.Bucket(fs.BucketName).Object(file.URL).Delete(ctx); err != nil {
		return fmt.Errorf("Can't delete file: %+v, error: %s", file, err.Error())
	}
	return fs.DB.Delete(&file).Error
}

// uploadFile uploads an object.
func (fs *FilesService) Upload(req *models.FileAddReq) (*models.File, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	fileName := uuid.New().String()

	// Upload an object with storage.Writer.
	wc := client.Bucket(fs.BucketName).Object(fileName).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.
	wc.ContentType = req.Type
	wc.Metadata = map[string]string{
		"MIMEType": req.Type,
		"FileName": req.Name,
	}

	dataURL, _ := dataurl.DecodeString(req.Value)
	if _, err := wc.Write(dataURL.Data); err != nil {
		return nil, err
	}
	if err := wc.Close(); err != nil {
		return nil, err
	}
	file := &models.File{
		ParentID: req.Parent,
		Name:     req.Name,
		URL:      fileName,
	}

	if err := fs.DB.Create(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func (fs *FilesService) GetUrl(ID uint64) (string, error) {
	var file models.File
	if err := fs.DB.First(&file, ID).Error; err != nil {
		return "", err
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
		// Expires: time.Now().Add(24 * 365 * time.Hour),
	}
	return client.Bucket(fs.BucketName).SignedURL(file.URL, opts)
}

func (fs *FilesService) List(filter models.ListFilter) ([]*models.File, error) {
	var files []*models.File
	err := fs.DB.Where("parent_id = ?", filter.ParentID).Find(&files).Error
	return files, err
}
