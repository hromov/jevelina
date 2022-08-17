package files

import (
	"context"
	"log"

	"github.com/hromov/jevelina/domain/misc/files"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
)

type Files struct {
	*gorm.DB
}

func NewFiles(db *gorm.DB, automigrate bool) *Files {
	if automigrate {
		if err := db.AutoMigrate(&models.File{}); err != nil {
			log.Printf("misc migration for %s error: %s\n", "files", err.Error())
		}
	}
	return &Files{db}
}

func (fs *Files) CreateFile(ctx context.Context, f files.FileCreateReq) (files.File, error) {
	file := models.FileFromDomain(f)
	if err := fs.DB.WithContext(ctx).Create(&file).Error; err != nil {
		return files.File{}, err
	}
	return file.ToDomain(), nil
}

func (fs *Files) DeleteFile(ctx context.Context, id uint64) error {
	return fs.DB.WithContext(ctx).Delete(&models.File{ID: id}).Error
}

func (fs *Files) GetFile(ctx context.Context, id uint64) (files.File, error) {
	var file models.File
	if err := fs.DB.WithContext(ctx).First(&file, id).Error; err != nil {
		return files.File{}, err
	}
	return file.ToDomain(), nil
}

func (fs *Files) GetFilesByParent(ctx context.Context, parentID uint64) ([]files.File, error) {
	var files []models.File
	if err := fs.DB.WithContext(ctx).Where("parent_id", parentID).Find(&files).Error; err != nil {
		return nil, err
	}
	return models.FilesToDomain(files), nil
}
