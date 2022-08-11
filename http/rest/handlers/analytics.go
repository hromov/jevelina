package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hromov/jevelina/useCases/analytics"
)

func LeadsBySource(as analytics.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := parseFilter(r.URL.Query())
		if err != nil {
			log.Println("Can't convert filter: ", err.Error())
			http.Error(w, "Filter error", http.StatusBadRequest)
			return
		}

		sources, err := as.LeadsBySource(r.Context(), filter.toAnalytics())
		if err != nil {
			log.Println("Can't get leads: ", err.Error())
			http.Error(w, "Getting leads error", http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(sources)
	}
}
