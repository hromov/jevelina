package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hromov/jevelina/services/events"
)

func EventsHandler(es events.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		//TODO: write normal filter or use listFilter
		listFilter := FilterFromQuery(r.URL.Query())
		eventFilter := events.EventFilter{
			UserID:   listFilter.ResponsibleID,
			ParentID: listFilter.ParentID,
			Limit:    listFilter.Limit,
			Offset:   listFilter.Offset,
		}
		switch {
		case strings.HasSuffix(r.URL.Path, "transfers"):
			eventFilter.EventParentType = events.TransferEvent
		default:
			http.NotFound(w, r)
			return
		}
		eResponse, err := es.List(r.Context(), eventFilter)
		if err != nil {
			log.Println("Can't get events error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		total := strconv.Itoa(int(eResponse.Total))
		w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
		w.Header().Set("X-Total-Count", total)
		_ = json.NewEncoder(w).Encode(eResponse.Events)
	}
}
