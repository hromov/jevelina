package files

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/hromov/jevelina/base"
	"github.com/hromov/jevelina/cdb/models"
	"github.com/vincent-petithory/dataurl"
)

const bucketName = "jevelina"

func createBucket() error {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "vorota-ua"

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		// log.Fatalf("Failed to create client: %v", err)
		return err
	}
	defer client.Close()

	// Sets the name for the new bucket.
	bucketName := "jevelina"

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	// Creates the new bucket.
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := bucket.Create(ctx, projectID, nil); err != nil {
		// log.Fatalf("Failed to create bucket: %v", err)
		return err
	}

	log.Printf("Bucket %v created.\n", bucketName)
	return nil
}

func DeleteFile(ID uint64) error {
	var file models.File
	if err := base.GetDB().First(&file).Error; err != nil {
		return err
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Bucket(bucketName).Object(file.URL).Delete(ctx); err != nil {
		return err
	}
	return base.GetDB().DB.Delete(&file).Error
}

// uploadFile uploads an object.
func UploadFile(req *models.FileAddReq) (*models.File, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	fileName := fmt.Sprintf("%s_%d", req.Name, time.Now().Unix())

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
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

	if err := base.GetDB().DB.Create(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func GetUrl(ID uint64) (string, error) {
	var file models.File
	if err := base.GetDB().First(&file, ID).Error; err != nil {
		return "", err
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
		// Expires: time.Now().Add(24 * 365 * time.Hour),
	}
	return client.Bucket(bucketName).SignedURL(file.URL, opts)
}

func List(filter models.ListFilter) ([]*models.File, error) {
	var files []*models.File
	err := base.GetDB().DB.Where("parent_id = ?", filter.ParentID).Find(&files).Error
	return files, err
}
