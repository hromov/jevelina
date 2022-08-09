package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/hromov/jevelina/domain/finances"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/http/rest/auth"
	"github.com/hromov/jevelina/utils/events"
)

func CompleteTransferHandler(f finances.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		//TODO: change to PUT
		if r.Method == "GET" {
			userValue := r.Context().Value(auth.KeyUser{})
			user, ok := userValue.(users.User)
			if !ok {
				http.Error(w, "Not a user", http.StatusForbidden)
				return
			}

			if err := f.CompleteTransfer(r.Context(), id, user.ID); err != nil {
				log.Printf("Can't save item with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}
	}
}

func TransferHandler(f finances.Service, es events.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "PUT":
			var transfer finances.Transfer
			if err = json.NewDecoder(r.Body).Decode(&transfer); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uint64(transfer.ID) != id {
				http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", id, transfer.ID), http.StatusBadRequest)
				return
			}

			oldTransfer, err := f.GetTransfer(r.Context(), id)
			if err != nil {
				log.Println("Can't get old transfer error: ", err.Error())
				http.Error(w, "Can't get old transfer", http.StatusInternalServerError)
				return
			}

			userValue := r.Context().Value(auth.KeyUser{})
			user, ok := userValue.(users.User)
			if !ok {
				http.Error(w, "Not a user", http.StatusForbidden)
				return
			}

			if needCheck(oldTransfer) && oldTransfer.Category != transfer.Category {
				event := events.NewEvent{
					UserID:          user.ID,
					ParentID:        oldTransfer.ID,
					Message:         fmt.Sprintf("%s > %s", oldTransfer.Category, transfer.Category),
					EventType:       events.CategoryChange,
					EventParentType: events.TransferEvent,
				}

				if err := es.Save(r.Context(), event); err != nil {
					log.Println("events save error: ", err)
					log.Println("Event to save: ", event.Message)
				}

			}

			//TODO: save event here at least on cat change
			if err = f.UpdateTransfer(r.Context(), transfer); err != nil {
				log.Printf("Can't save item with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		case "DELETE":
			userValue := r.Context().Value(auth.KeyUser{})
			user, ok := userValue.(users.User)
			if !ok {
				http.Error(w, "Not a user", http.StatusForbidden)
				return
			}
			if err := f.DeleteTransfer(r.Context(), id, user.ID); err != nil {
				log.Printf("Can't delete item with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}

	}
}

func TransfersHandler(f finances.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			transfer := finances.Transfer{}
			if err := json.NewDecoder(r.Body).Decode(&transfer); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			userValue := r.Context().Value(auth.KeyUser{})
			user, ok := userValue.(users.User)
			if !ok {
				http.Error(w, "Not a user", http.StatusForbidden)
				return
			}
			transfer.CreatedBy = user.ID
			transfer, err := f.CreateTransfer(r.Context(), transfer)
			if err != nil {
				log.Printf("Can't create transfer. Error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}

			_ = json.NewEncoder(w).Encode(transfer)
			return
		}

		tResponse, err := f.Transfers(r.Context(), FinFilter(r.URL.Query()))
		if err != nil {
			log.Println("Can't get transfer error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		b, err := json.Marshal(tResponse.Transfers)
		if err != nil {
			log.Println("Can't json.Marshal(transfers) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		total := strconv.Itoa(int(tResponse.Total))
		w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
		w.Header().Set("X-Total-Count", total)
		fmt.Fprint(w, string(b))
	}
}
func CategoriesHandler(f finances.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := f.TransferCategories(r.Context())
		if err != nil {
			http.Error(w, "Can't get transfer categories error: %s"+err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(categories)
	}
}

func CategoriesSumHandler(f finances.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sums, err := f.SumByCategory(r.Context(), FinFilter(r.URL.Query()))
		if err != nil {
			http.Error(w, "Can't get sum by category error: %s"+err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(sums)
	}
}

func needCheck(t finances.Transfer) bool {
	return t.Completed || !t.DeletedAt.IsZero()
}
