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

func TagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := cdb.Misc()
	var tag *models.Tag

	switch r.Method {
	case "GET":
		tag, err = c.Tag(uint8(ID))
		if err != nil {
			log.Println("Can't get tag error: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		b, err := json.Marshal(tag)
		if err != nil {
			log.Println("Can't json.Marshal(tag) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(b)
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&tag); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(tag.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, tag.ID), http.StatusBadRequest)
			return
		}
		//channge to base.DB?
		if err = c.DB.Save(tag).Error; err != nil {
			log.Printf("Can't update tag with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	case "DELETE":

		if err = c.DB.Delete(&models.Tag{ID: uint8(ID)}).Error; err != nil {
			log.Printf("Can't delete tag with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	}

}

func TagsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tags" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		tag := new(models.Tag)
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := cdb.GetDB()
		//channge to base.DB?
		if err := c.DB.Omit(clause.Associations).Create(tag).Error; err != nil {
			log.Printf("Can't create tag. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(tag)
		if err != nil {
			log.Println("Can't json.Marshal(tag) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(b)
		// it said that its already ok now
		// w.WriteHeader(http.StatusOK)
		return
	}

	c := cdb.Misc()
	tagsResponse, err := c.Tags()
	if err != nil {
		log.Println("Can't get tags error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(tagsResponse)
	if err != nil {
		log.Println("Can't json.Marshal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	total := strconv.Itoa(len(tagsResponse))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", total)
	_, _ = w.Write(b)
}
