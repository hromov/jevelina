package models

import (
	"time"

	"github.com/hromov/jevelina/domain/misc/files"
)

type File struct {
	ID        uint64 `gorm:"primaryKey"`
	ParentID  uint64 `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"size:32"`
	URL       string `gorm:"size:128"`
}

func FileFromDomain(f files.FileCreateReq) File {
	return File{
		ParentID: f.ParentID,
		Name:     f.Name,
		URL:      f.URL,
	}
}

func (f *File) ToDomain() files.File {
	return files.File{
		ID:        f.ID,
		ParentID:  f.ParentID,
		CreatedAt: f.CreatedAt,
		Name:      f.Name,
		URL:       f.URL,
	}
}

func FilesToDomain(items []File) []files.File {
	converted := make([]files.File, len(items))
	for i, f := range items {
		converted[i] = f.ToDomain()
	}
	return converted
}

type FileAddReq struct {
	Parent uint64
	Name   string
	Type   string
	Value  string
}
