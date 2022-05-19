package contacts

import (
	"database/sql"
	"unicode"

	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const fullSearch = "name LIKE @query OR second_name LIKE @query OR phone LIKE @query OR second_phone LIKE @query OR email LIKE @query OR second_email LIKE @query OR url LIKE @query OR city LIKE @query OR address LIKE @query OR position LIKE @query"
const phonesOnly = "phone LIKE @query OR second_phone LIKE @query"

type Contacts struct {
	*gorm.DB
}

func digitsOnly(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func (c *Contacts) List(filter models.ListFilter) (*models.ContactsResponse, error) {
	cr := &models.ContactsResponse{}
	q := c.DB.Preload(clause.Associations).Order("created_at desc").Limit(filter.Limit).Offset(filter.Offset)
	if filter.Query != "" {
		searchType := ""
		if digitsOnly(filter.Query) {
			searchType = phonesOnly
		} else {
			searchType = fullSearch
		}
		q = q.Where(searchType, sql.Named("query", "%"+filter.Query+"%"))
	}
	if filter.TagID != 0 {
		IDs := []uint{}
		c.DB.Raw("select contact_id from contacts_tags WHERE tag_id = ?", filter.TagID).Scan(&IDs)
		q.Find(&cr.Contacts, IDs)
	} else {
		q.Find(&cr.Contacts)
	}

	if result := q.Count(&cr.Total); result.Error != nil {
		return nil, result.Error
	}
	return cr, nil
}

func (c *Contacts) ByID(ID uint64) (*models.Contact, error) {
	// log.Println(limit, offset, query, query == "")
	var contact models.Contact

	if result := c.DB.Preload(clause.Associations).First(&contact, ID); result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

func (c *Contacts) ByPhone(phone string) (*models.Contact, error) {
	contact := new(models.Contact)
	if err := c.DB.Where(phonesOnly, sql.Named("query", phone)).First(contact).Error; err != nil {
		return nil, err
	}
	return contact, nil
}
