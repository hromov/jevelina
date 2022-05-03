package orders

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hromov/jevelina/base"
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm/clause"
)

type CreateLeadReq struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description,omitempty"`

	ClientName  string `json:"clientname"`
	ClientEmail string `json:"clientemail,omitempty"`
	ClientPhone string `json:"clientphone,omitempty"`

	Source       string `json:"source,omitempty"`
	Product      string `json:"product,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`

	UserEmail string `json:"user_email,omitempty"`
	UserHash  string `json:"user_hash,omitempty"`

	CID string `gorm:"size:64"`
	UID string `gorm:"size:64"`
	TID string `gorm:"size:64"`

	UtmID       string `gorm:"size:64"`
	UtmSource   string `gorm:"size:64"`
	UtmMedium   string `gorm:"size:64"`
	UtmCampaign string `gorm:"size:64"`

	Domain string `gorm:"size:128"`
}

type LeadOrContact interface {
	models.Lead | models.Contact
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/neworder" {
		http.NotFound(w, r)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Only POST is allowed", http.StatusForbidden)
		return
	}

	c := new(CreateLeadReq)
	if err := json.NewDecoder(r.Body).Decode(c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if c.UserEmail == "" || c.UserHash == "" {
		http.Error(w, "User Email and Hash are required", http.StatusBadRequest)
		return
	}
	DB := base.GetDB()
	user, err := DB.Misc().UserByEmail(c.UserEmail)
	if err != nil {
		http.Error(w, "Cant find user with email: "+user.Email, http.StatusNotFound)
		return
	}

	if user.Hash != c.UserHash {
		http.Error(w, "Wrong user-hash values", http.StatusForbidden)
		return
	}
	contact, err := createOrGetContact(c, user)
	if err != nil {
		log.Println("Can't create contact error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	//TODO: lead goes here

	b, err := json.Marshal(contact)
	if err != nil {
		log.Println("Can't json.Marshal(user) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(b))
	// it said that its already ok now
	// w.WriteHeader(http.StatusOK)
	return
}

func createOrGetContact(c *CreateLeadReq, user *models.User) (*models.Contact, error) {
	DB := base.GetDB()
	DB.Contacts().List(models.ListFilter{Query: c.ClientPhone})
	var contact *models.Contact
	var err error
	if c.ClientPhone != "" && len(c.ClientPhone) > 5 {
		contact, err = DB.Contacts().ByPhone(c.ClientPhone)
		if err != nil {
			return nil, err
		}
	}
	if contact == nil {
		contact = &models.Contact{
			Name:  c.ClientName,
			Phone: c.ClientPhone,
			Email: c.ClientEmail,
		}
		//TODO: move to generics in 1.19, if possible
		item := contact
		item.Analytics.CID = c.CID
		item.Analytics.UID = c.UID
		item.Analytics.TID = c.TID
		item.Analytics.UtmID = c.UtmID
		item.Analytics.UtmSource = c.UtmSource
		item.Analytics.UtmMedium = c.UtmMedium
		item.Analytics.UtmCampaign = c.UtmCampaign
		item.Analytics.Domain = c.Domain

		if c.Source != "" {
			if source, err := DB.Misc().SourceByName(c.Source); err == nil && source != nil {
				contact.SourceID = &source.ID
			}
		}
	} else {
		//check for updated fields
		if c.ClientName != "" && strings.Compare(contact.Name, c.ClientName) != 0 {
			if contact.Name == "" {
				contact.Name = c.ClientName
			} else if contact.SecondName == "" {
				contact.SecondName = c.ClientName
			}
		}
		if c.ClientEmail != "" && strings.Compare(contact.Email, c.ClientEmail) != 0 {
			if contact.Email == "" {
				contact.Email = c.ClientEmail
			} else if contact.SecondEmail == "" {
				contact.SecondEmail = c.ClientEmail
			}
		}
	}
	contact.ResponsibleID = &user.ID

	if contact.ID == 0 {
		if err := DB.Omit(clause.Associations).Create(contact).Error; err != nil {
			return nil, err
		}
	} else {
		if err := DB.Omit(clause.Associations).Save(contact).Error; err != nil {
			return nil, err
		}
	}
	return contact, nil
}
