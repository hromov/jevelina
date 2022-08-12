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

func Source(ms misc.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			source, err := ms.GetSource(r.Context(), uint32(id))
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
			encode(w, source)
			return
		case "PUT":
			source := misc.Source{}
			if err = json.NewDecoder(r.Body).Decode(&source); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uint64(source.ID) != id {
				http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", id, source.ID), http.StatusBadRequest)
				return
			}

			if err := ms.UpdateSource(r.Context(), source); err != nil {
				log.Printf("Can't update source with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		case "DELETE":
			if err := ms.DeleteSource(r.Context(), uint32(id)); err != nil {
				log.Printf("Can't delete source with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
	}
}

func Sources(ms misc.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			source := misc.Source{}
			if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			source, err := ms.CreateSource(r.Context(), source)
			if err != nil {
				log.Println("Can't create source error: ", err.Error())
				http.Error(w, "Can't create source error", http.StatusInternalServerError)
			}

			w.WriteHeader(http.StatusCreated)
			encode(w, source)
			return
		}

		sources, err := ms.ListSources(r.Context())
		if err != nil {
			log.Println("Can't get sources error: ", err)
			http.Error(w, "Can't get sources list", http.StatusInternalServerError)
		}
		encode(w, sources)
	}
}
