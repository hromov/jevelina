package mysql

import (
	"fmt"

	"github.com/hromov/jevelina/storage/mysql/dao/misc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Storage struct {
	*misc.Misc
}

func NewStorage(dns string) (*Storage, error) {
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database error: %s", err.Error())
	}
	return &Storage{
		misc.NewMisc(db),
	}, nil
}
