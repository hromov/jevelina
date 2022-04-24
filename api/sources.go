package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/cdb/models"
	"github.com/hromov/jevelina/base"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SourceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := base.GetDB().Misc()
	var source *models.Source

	switch r.Method {
	case "GET":
		source, err = c.Source(uint8(ID))
		if err != nil {
			log.Println("Can't get source error: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		b, err := json.Marshal(source)
		if err != nil {
			log.Println("Can't json.Marchal(source) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&source); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(source.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, source.ID), http.StatusBadRequest)
			return
		}
		//channge to base.DB?
		if err = c.DB.Save(source).Error; err != nil {
			log.Printf("Can't update source with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	case "DELETE":

		if err = c.DB.Delete(&models.Source{ID: uint8(ID)}).Error; err != nil {
			log.Printf("Can't delete source with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	}

}

func SourcesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/sources" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		source := new(models.Source)
		if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := base.GetDB()
		//channge to base.DB?
		if err := c.DB.Omit(clause.Associations).Create(source).Error; err != nil {
			log.Printf("Can't create source. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(source)
		if err != nil {
			log.Println("Can't json.Marchal(source) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
		// it said that its already ok now
		// w.WriteHeader(http.StatusOK)
		return
	}

	c := base.GetDB().Misc()
	sourcesResponse, err := c.Sources()
	if err != nil {
		log.Println("Can't get sources error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(sourcesResponse)
	if err != nil {
		log.Println("Can't json.Marchal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	total := strconv.Itoa(len(sourcesResponse))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", total)
	fmt.Fprintf(w, string(b))
}
