package mysql

import (
	"fmt"

	"github.com/hromov/jevelina/storage/mysql/dao/contacts"
	"github.com/hromov/jevelina/storage/mysql/dao/leads"
	"github.com/hromov/jevelina/storage/mysql/dao/misc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Storage struct {
	*misc.Misc
	*leads.Leads
	*contacts.Contacts
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
		leads.NewLeads(db),
		contacts.NewContacts(db),
	}, nil
}
