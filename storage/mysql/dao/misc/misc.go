package misc

import (
	"log"

	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
)

type Misc struct {
	db *gorm.DB
}

func NewMisc(db *gorm.DB, automigrate bool) *Misc {
	if automigrate {
		if err := db.AutoMigrate(&models.Manufacturer{}); err != nil {
			log.Printf("misc migration for %s error: %s\n", "manufacturer", err.Error())
		}

		if err := db.AutoMigrate(&models.Product{}); err != nil {
			log.Printf("misc migration for %s error: %s\n", "product", err.Error())
		}

		if err := db.AutoMigrate(&models.Source{}); err != nil {
			log.Printf("misc migration for %s error: %s\n", "source", err.Error())
		}

		if err := db.AutoMigrate(&models.Task{}); err != nil {
			log.Printf("misc migration for %s error: %s\n", "task", err.Error())
		}
	}
	return &Misc{
		db,
	}
}

// func (m *Misc) Tags() ([]models.Tag, error) {
// 	var items []models.Tag
// 	if result := m.DB.Find(&items); result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return items, nil
// }

// func (m *Misc) Tag(ID uint8) (*models.Tag, error) {
// 	var item models.Tag
// 	if result := m.DB.First(&item, ID); result.Error != nil {
// 		return nil, result.Error

// 	}
// 	return &item, nil
// }
