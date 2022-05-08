package files_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
)

func FilesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fileAddRequest := new(models.FileAddReq)
		if err := json.NewDecoder(r.Body).Decode(fileAddRequest); err != nil {
			http.Error(w, "File Decode Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		file, err := cdb.Files().Upload(fileAddRequest)
		if err != nil {
			http.Error(w, "File Uploading Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		b, _ := json.Marshal(file)
		fmt.Fprintf(w, string(b))
		return
	case "GET":
		files, err := cdb.Files().List(api.FilterFromQuery(r.URL.Query()))
		if err != nil {
			log.Println("Can't get transfer error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		b, _ := json.Marshal(files)
		fmt.Fprintf(w, string(b))
		return
	}
	return
}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		url, err := cdb.Files().GetUrl(ID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Can't get url error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		//%q is the key
		fmt.Fprintf(w, "%q", url)
		return
	case "DELETE":
		if err := cdb.Files().Delete(ID); err != nil {
			http.Error(w, "Can't delete file error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	return
}
