package handlers

import (
	"encoding/json"
	"log"

	"net/http"

	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/useCases/orders"
)

const randomUserEmail = "random@random.org"

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

	CID string
	UID string
	TID string

	UtmID       string
	UtmSource   string
	UtmMedium   string
	UtmCampaign string

	Domain string
}

func (c *CreateLeadReq) ToOrder() orders.Order {
	//TODO: some how get prodID, ManufID, SourceID by name
	return orders.Order{
		Name:        c.Name,
		Price:       c.Price,
		Description: c.Description,

		ClientName:  c.ClientName,
		ClientEmail: c.ClientEmail,
		ClientPhone: c.ClientPhone,

		ProductID:      0,
		ManufacturerID: 0,
		SourceID:       0,

		CID: c.CID,
		UID: c.UID,
		TID: c.TID,

		UtmID:       c.UtmID,
		UtmSource:   c.UtmSource,
		UtmMedium:   c.UtmMedium,
		UtmCampaign: c.UtmCampaign,

		Domain: c.Domain,
	}
}

func Order(us users.Service, os orders.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		c := CreateLeadReq{}

		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if c.UserEmail == "" || c.UserHash == "" {
			http.Error(w, "User Email and Hash are required", http.StatusBadRequest)
			return
		}
		// TODO: some kind of auth for services
		user, err := us.GetByEmail(ctx, c.UserEmail)
		if err != nil || user.ID == 0 {
			http.Error(w, "Cant find user with email: "+c.UserEmail, http.StatusBadRequest)
			return
		}

		if user.Hash != c.UserHash {
			http.Error(w, "Wrong user-hash values", http.StatusForbidden)
			return
		}

		// end
		if user.Email == randomUserEmail {
			if err := os.Create(ctx, c.ToOrder()); err != nil {
				log.Println("Can't create order error: ", err.Error())
				http.Error(w, "Can't create order", http.StatusInternalServerError)
				return
			}
		} else {
			if err := os.CreateForUser(ctx, c.ToOrder(), user); err != nil {
				log.Println("Can't create order error: ", err.Error())
				http.Error(w, "Can't create order", http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
