package events_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/events/") {
		http.NotFound(w, r)
		return
	}
	events := cdb.Events()
	//TODO: write normal filter or use listFilter
	listFilter := api.FilterFromQuery(r.URL.Query())
	eventFilter := models.EventFilter{
		UserID:   listFilter.ResponsibleID,
		ParentID: listFilter.ParentID,
		Limit:    listFilter.Limit,
		Offset:   listFilter.Offset,
	}
	switch {
	case strings.HasSuffix(r.URL.Path, "transfers"):
		eventFilter.EventParentType = models.TransferEvent
	default:
		http.NotFound(w, r)
		return
	}
	eResponse, err := events.List(eventFilter)
	if err != nil {
		log.Println("Can't get events error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	b, err := json.Marshal(eResponse)
	if err != nil {
		log.Println("Can't json.Marshal(events) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	total := strconv.Itoa(int(eResponse.Total))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", total)
	fmt.Fprint(w, string(b))
}
