package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ManufacturerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := cdb.Misc()
	var manufacturer *models.Manufacturer

	switch r.Method {
	case "GET":
		manufacturer, err = c.Manufacturer(uint16(ID))
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
		b, err := json.Marshal(manufacturer)
		if err != nil {
			log.Println("Can't json.Marshal(manufacturer) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&manufacturer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(manufacturer.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, manufacturer.ID), http.StatusBadRequest)
			return
		}

		if err = c.DB.Omit(clause.Associations).Save(manufacturer).Error; err != nil {
			log.Printf("Can't update manufacturer with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	case "DELETE":

		if err = c.DB.Delete(&models.Manufacturer{ID: uint16(ID)}).Error; err != nil {
			log.Printf("Can't delete manufacturer with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	}

}

func ManufacturersHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/manufacturers" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		manufacturer := new(models.Manufacturer)
		if err := json.NewDecoder(r.Body).Decode(&manufacturer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := cdb.GetDB()
		if err := c.DB.Omit(clause.Associations).Create(manufacturer).Error; err != nil {
			log.Printf("Can't create manufacturer. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		b, err := json.Marshal(manufacturer)
		if err != nil {
			log.Println("Can't json.Marshal(manufacturer) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
		return
	}

	c := cdb.Misc()
	manufacturersResponse, err := c.Manufacturers()
	if err != nil {
		log.Println("Can't get manufacturers error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	b, err := json.Marshal(manufacturersResponse)
	if err != nil {
		log.Println("Can't json.Marshal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	total := strconv.Itoa(len(manufacturersResponse))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", total)
	fmt.Fprint(w, string(b))
}
