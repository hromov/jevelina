package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hromov/cdb"
)

const pageSize = 50

func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contacts" {
		http.NotFound(w, r)
		return
	}

	page := r.URL.Query().Get("page")
	limit, offset := pageSize, 0
	if page != "" {
		p, err := strconv.Atoi(page)
		if err == nil {
			limit = pageSize
			offset = p * limit
		}
	}
	query := r.URL.Query().Get("query")

	contacts, err := cdb.Contacts(limit, offset, query)
	if err != nil {
		log.Println("Can't get contacts error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(contacts)
	if err != nil {
		log.Println("Can't json.Marchal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(b))
}
