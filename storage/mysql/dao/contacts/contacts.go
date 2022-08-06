package contacts

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"unicode"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/storage/mysql/dao/misc"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
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

func (c *Contacts) Contacts(ctx context.Context, filter contacts.Filter) (contacts.ContactsResponse, error) {
	cr := &models.ContactsResponse{}
	q := c.DB.WithContext(ctx).Preload(clause.Associations).Order("created_at desc").Limit(filter.Limit).Offset(filter.Offset)
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

	if err := q.Count(&cr.Total).Error; err != nil {
		return contacts.ContactsResponse{}, err
	}
	resp := contacts.ContactsResponse{
		Contacts: make([]contacts.Contact, len(cr.Contacts)),
		Total:    cr.Total,
	}
	for i, contact := range cr.Contacts {
		resp.Contacts[i] = contact.ToDomain()
	}
	return resp, nil
}

func (c *Contacts) ByID(ctx context.Context, ID uint64) (contacts.Contact, error) {
	var contact models.Contact

	if result := c.DB.Unscoped().Preload(clause.Associations).First(&contact, ID); result.Error != nil {
		return contacts.Contact{}, result.Error
	}

	return contact.ToDomain(), nil
}

func (c *Contacts) ByPhone(ctx context.Context, phone string) (contacts.Contact, error) {
	if phone == "" || len(phone) < 6 {
		return contacts.Contact{}, errors.New("Phone should be at least 6 char length")
	}
	var contact models.Contact
	if err := c.DB.Where(phonesOnly, sql.Named("query", phone)).First(contact).Error; err != nil {
		return contacts.Contact{}, err
	}
	return contact.ToDomain(), nil
}

func (c *Contacts) DeleteContact(ctx context.Context, id uint64) error {
	if err := c.DB.Delete(&models.Contact{ID: id}).Error; err != nil {
		return err
	}
	if err := misc.DeleteTaskByParent(c.DB, id); err != nil {
		log.Printf("Error while deliting tasks for contact: %s", err.Error())
	}
	return nil
}

func (c *Contacts) CreateContact(ctx context.Context, newContact contacts.ContactRequest) (contacts.Contact, error) {
	dbContact := models.ContactFromDomain(newContact)
	if err := c.DB.WithContext(ctx).Omit(clause.Associations).Create(&dbContact).Error; err != nil {
		return contacts.Contact{}, err
	}
	return c.ByID(ctx, dbContact.ID)
}

func (c *Contacts) UpdateContact(ctx context.Context, contact contacts.ContactRequest) error {
	dbContact := models.ContactFromDomain(contact)
	if err := c.DB.Debug().WithContext(ctx).Omit(clause.Associations).Model(&models.Contact{}).Where("id", contact.ID).Updates(&dbContact).Error; err != nil {
		return err
	}
	return nil
}
