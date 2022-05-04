package orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/hromov/jevelina/base"
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm/clause"
)

const randomUserEmail = "random@random.org"

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/orders" {
		http.NotFound(w, r)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Only POST is allowed", http.StatusForbidden)
		return
	}

	c := new(models.CreateLeadReq)
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
	if err != nil || user == nil {
		http.Error(w, "Cant find user with email: "+c.UserEmail, http.StatusBadRequest)
		return
	}

	if user.Hash != c.UserHash {
		http.Error(w, "Wrong user-hash values", http.StatusForbidden)
		return
	}
	if user.Email == randomUserEmail {
		user, err = getRandomUser()
		if err != nil {
			log.Println("Can't get random user error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	}
	contact, err := createOrGetContact(c, user)
	if err != nil {
		log.Println("Can't create contact error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	lead, err := createLead(c, contact)
	if err != nil {
		log.Println("Can't create lead error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	if err = createTask(c, lead); err != nil {
		log.Println("Can't create task error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(lead)
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

func createOrGetContact(c *models.CreateLeadReq, user *models.User) (*models.Contact, error) {
	DB := base.GetDB()
	DB.Contacts().List(models.ListFilter{Query: c.ClientPhone})
	var contact *models.Contact
	// var err error
	if c.ClientPhone != "" && len(c.ClientPhone) > 5 {
		contact, _ = DB.Contacts().ByPhone(c.ClientPhone)
		// if err != nil {
		// 	return nil, err
		// }
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
		//if no UID was provided - use phone
		if c.UID == "" {
			item.Analytics.UID = contact.Phone
		} else {
			item.Analytics.UID = c.UID
		}
		item.Analytics.TID = c.TID
		item.Analytics.UtmID = c.UtmID
		item.Analytics.UtmSource = c.UtmSource
		item.Analytics.UtmMedium = c.UtmMedium
		item.Analytics.UtmCampaign = c.UtmCampaign
		item.Analytics.Domain = c.Domain

		if c.Source != "" {
			if source, _ := DB.Misc().SourceByName(c.Source); source != nil {
				contact.SourceID = &source.ID
			}
		}

		contact.ResponsibleID = &user.ID
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

func createLead(c *models.CreateLeadReq, contact *models.Contact) (*models.Lead, error) {
	DB := base.GetDB()
	lead := &models.Lead{
		Name:          c.Name,
		Budget:        uint32(c.Price),
		ResponsibleID: contact.ResponsibleID,
		ContactID:     &contact.ID,
	}
	if step, _ := DB.Misc().DefaultStep(); step != nil {
		lead.StepID = &step.ID
	}
	item := lead
	item.Analytics.CID = c.CID
	item.Analytics.UID = c.UID
	item.Analytics.TID = c.TID
	item.Analytics.UtmID = c.UtmID
	item.Analytics.UtmSource = c.UtmSource
	item.Analytics.UtmMedium = c.UtmMedium
	item.Analytics.UtmCampaign = c.UtmCampaign
	item.Analytics.Domain = c.Domain

	if c.Source != "" {
		if source, _ := DB.Misc().SourceByName(c.Source); source != nil {
			lead.SourceID = &source.ID
		}
	}
	if c.Product != "" {
		if product, _ := DB.Misc().ProductByName(c.Product); product != nil {
			lead.ProductID = &product.ID
		}
	}
	if c.Manufacturer != "" {
		if manuf, _ := DB.Misc().ManufacturerByName(c.Manufacturer); manuf != nil {
			lead.ManufacturerID = &manuf.ID
		}
	}
	if err := DB.Omit(clause.Associations).Create(lead).Error; err != nil {
		return nil, err
	}
	return lead, nil
}

func createTask(c *models.CreateLeadReq, lead *models.Lead) error {
	DB := base.GetDB()
	task := new(models.Task)
	if c.Description != "" {
		task.Description = c.Description
	} else {
		task.Description = "Call me!"
	}
	t := time.Now()
	task.DeadLine = &t
	task.ParentID = lead.ID
	task.ResponsibleID = lead.ResponsibleID
	if err := DB.Omit(clause.Associations).Create(task).Error; err != nil {
		return err
	}
	return nil
}

func getRandomUser() (*models.User, error) {
	users, err := base.GetDB().Misc().Users()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("No users found in base")
	}

	rand.Seed(time.Now().UnixNano())
	var user *models.User
	appropriateSeen := false
	for user == nil {
		for _, u := range users {
			if u.Distribution == 0.0 {
				continue
			}
			appropriateSeen = true
			r := rand.Float32()
			// log.Printf("checking %+v, rand = %.2f, good = %v\n", u, r, u.Distribution >= r)
			if u.Distribution >= r {
				return &u, nil
			}
		}
		if !appropriateSeen {
			return nil, errors.New("No user with distribution more then 0 was found")
		}
	}
	return nil, errors.New("should never be called")

}
