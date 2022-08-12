package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hromov/jevelina/domain/misc/files"
)

func Files(fs files.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			fileAddRequest := files.FileAddReq{}
			if err := json.NewDecoder(r.Body).Decode(&fileAddRequest); err != nil {
				http.Error(w, "File Decode Error: "+err.Error(), http.StatusBadRequest)
				return
			}

			file, err := fs.Upload(r.Context(), fileAddRequest)
			if err != nil {
				http.Error(w, "File Uploading Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			encode(w, file)
			return
		case "GET":
			parentString := r.URL.Query().Get("parent")
			if parentString == "" {
				http.Error(w, "parent query should not be empty", http.StatusBadRequest)
				return
			}
			parentID, err := strconv.ParseUint(parentString, 10, 64)
			if err != nil {
				log.Println("Can't parse parent id from query error: ", err.Error())
				http.Error(w, "parent parsing error", http.StatusBadRequest)
				return
			}

			files, err := fs.GetByParent(r.Context(), parentID)
			if err != nil {
				log.Println("Can't get files error: " + err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			encode(w, files)
			return
		}
	}
}

func File(fs files.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			url, err := fs.GetUrl(r.Context(), id)
			if err != nil {
				log.Println(err)
				http.Error(w, "Can't get url error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "%q", url)
			return
		case "DELETE":
			if err := fs.Delete(r.Context(), id); err != nil {
				log.Println(err)
				http.Error(w, "Can't delete file error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
}
