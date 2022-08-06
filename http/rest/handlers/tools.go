package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func getID(r *http.Request) (uint64, error) {
	vars := mux.Vars(r)
	return strconv.ParseUint(vars["id"], 10, 32)
}

func timeOrNull(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func encode(w http.ResponseWriter, val interface{}) {
	if err := json.NewEncoder(w).Encode(val); err != nil {
		log.Println("can't encode value error: ", err.Error())
	}
}
