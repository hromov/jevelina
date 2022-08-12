package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hromov/jevelina/domain/misc"
	"gorm.io/gorm"
)

func Manufacturer(ms misc.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			manufacturer, err := ms.GetManufacturer(r.Context(), uint32(id))
			if err != nil {
				log.Println("Can't get manufacturer error: " + err.Error())
				if errors.Is(err, gorm.ErrRecordNotFound) {
					http.NotFound(w, r)
				} else {
					http.Error(w, http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
				return
			}
			encode(w, manufacturer)
			return
		case "PUT":
			manufacturer := misc.Manufacturer{}
			if err = json.NewDecoder(r.Body).Decode(&manufacturer); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uint64(manufacturer.ID) != id {
				http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", id, manufacturer.ID), http.StatusBadRequest)
				return
			}

			if err := ms.UpdateManufacturer(r.Context(), manufacturer); err != nil {
				log.Printf("Can't update manufacturer with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		case "DELETE":
			if err := ms.DeleteManufacturer(r.Context(), uint32(id)); err != nil {
				log.Printf("Can't delete manufacturer with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
	}
}

func Manufacturers(ms misc.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			manufacturer := misc.Manufacturer{}
			if err := json.NewDecoder(r.Body).Decode(&manufacturer); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			manufacturer, err := ms.CreateManufacturer(r.Context(), manufacturer)
			if err != nil {
				log.Println("Can't create manufacturer error: ", err.Error())
				http.Error(w, "Can't create manufacturer error", http.StatusInternalServerError)
			}
			return
		}

		manufacturers, err := ms.ListManufacturers(r.Context())
		if err != nil {
			log.Println("Can't get manufacturers error: ", err)
			http.Error(w, "Can't get manufacturers list", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
		encode(w, manufacturers)
	}
}
