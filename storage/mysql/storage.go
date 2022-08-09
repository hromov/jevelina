package mysql

import (
	"fmt"

	"github.com/hromov/jevelina/storage/mysql/dao/contacts"
	"github.com/hromov/jevelina/storage/mysql/dao/files"
	"github.com/hromov/jevelina/storage/mysql/dao/finance"
	"github.com/hromov/jevelina/storage/mysql/dao/leads"
	"github.com/hromov/jevelina/storage/mysql/dao/misc"
	"github.com/hromov/jevelina/storage/mysql/dao/users"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Storage struct {
	*users.Users
	*misc.Misc
	*leads.Leads
	*contacts.Contacts
	*files.Files
	*finance.Finance
}

func NewStorage(dns string) (*Storage, error) {
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database error: %s", err.Error())
	}
	return &Storage{
		users.NewUsers(db),
		misc.NewMisc(db),
		leads.NewLeads(db),
		contacts.NewContacts(db),
		files.NewFiles(db),
		finance.NewFinance(db),
	}, nil
}
