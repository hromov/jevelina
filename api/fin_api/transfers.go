package fin_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/auth"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
)

func CompleteTransferHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	fin := cdb.Finance()
	//or PUT?
	if r.Method == "GET" {
		user, err := auth.GetCurrentUser(r)
		if err != nil || user == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		if err := fin.CompleteTransfer(ID, user.ID); err != nil {
			log.Printf("Can't save item with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}
	fin := cdb.Finance()
	switch r.Method {
	case "PUT":
		var transfer *models.Transfer
		if err = json.NewDecoder(r.Body).Decode(&transfer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(transfer.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, transfer.ID), http.StatusBadRequest)
			return
		}

		user, err := auth.GetCurrentUser(r)
		if err != nil || user == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		if err = fin.UpdateTransfer(user.ID, transfer); err != nil {
			log.Printf("Can't save item with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	case "DELETE":
		user, err := auth.GetCurrentUser(r)
		if err != nil || user == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		if err := fin.DeleteTransfer(ID, user.ID); err != nil {
			log.Printf("Can't delete item with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	}

}

func TransfersHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/transfers" {
		http.NotFound(w, r)
		return
	}
	fin := cdb.Finance()
	if r.Method == "POST" {
		item := new(models.Transfer)
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := auth.GetCurrentUser(r)
		if err != nil || user == nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item.CreatedBy = user.ID
		item, err = fin.CreateTransfer(item)
		if err != nil {
			log.Printf("Can't create transfer. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		b, err := json.Marshal(item)
		if err != nil {
			log.Println("Can't json.Marshal(transfer) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
	}

	tResponse, err := fin.Transfers(api.FilterFromQuery(r.URL.Query()))
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

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := cdb.Finance().Categories()
	if err != nil {
		http.Error(w, "Can't get transfer categories error: %s"+err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(categories)
	fmt.Fprint(w, string(b))
}

func CategoriesSumHandler(w http.ResponseWriter, r *http.Request) {
	sums, err := cdb.Finance().SumByCategory(api.FilterFromQuery(r.URL.Query()))
	if err != nil {
		http.Error(w, "Can't get sum by category error: %s"+err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(sums)
	fmt.Fprint(w, string(b))
}
