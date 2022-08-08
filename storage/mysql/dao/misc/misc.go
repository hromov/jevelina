package misc

import (
	"log"

	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"gorm.io/gorm"
)

type Misc struct {
	db *gorm.DB
}

func NewMisc(db *gorm.DB) *Misc {

	if err := db.AutoMigrate(&models.Manufacturer{}); err != nil {
		log.Printf("misc migration for %s error: %s\n", "manufacturer", err.Error())
	}

	if err := db.AutoMigrate(&models.Product{}); err != nil {
		log.Printf("misc migration for %s error: %s\n", "product", err.Error())
	}

	if err := db.AutoMigrate(&models.Role{}); err != nil {
		log.Printf("misc migration for %s error: %s\n", "role", err.Error())
	}

	if err := db.AutoMigrate(&models.Source{}); err != nil {
		log.Printf("misc migration for %s error: %s\n", "source", err.Error())
	}

	if err := db.AutoMigrate(&models.Task{}); err != nil {
		log.Printf("misc migration for %s error: %s\n", "task", err.Error())
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Printf("misc migration for %s error: %s\n", "user", err.Error())
	}

	user := models.User{ID: 1}
	if err := db.First(&user).Error; err != nil {
		if err := InitUsers(db); err != nil {
			log.Printf("Can't create base roles error: %s", err.Error())
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
