package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/base"
	"github.com/hromov/jevelina/cdb/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := base.GetDB().Misc()
	var product *models.Product

	switch r.Method {
	case "GET":
		product, err = c.Product(uint32(ID))
		if err != nil {
			log.Println("Can't get product error: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		b, err := json.Marshal(product)
		if err != nil {
			log.Println("Can't json.Marshal(product) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(product.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, product.ID), http.StatusBadRequest)
			return
		}

		//channge to base.DB?
		if err = c.DB.Omit(clause.Associations).Save(product).Error; err != nil {
			log.Printf("Can't update product with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	case "DELETE":

		if err = c.DB.Delete(&models.Product{ID: uint32(ID)}).Error; err != nil {
			log.Printf("Can't delete product with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	}

}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/products" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		product := new(models.Product)
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := base.GetDB()
		//channge to base.DB?
		if err := c.DB.Omit(clause.Associations).Create(product).Error; err != nil {
			log.Printf("Can't create product. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(product)
		if err != nil {
			log.Println("Can't json.Marshal(product) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
		// it said that its already ok now
		// w.WriteHeader(http.StatusOK)
		return
	}

	c := base.GetDB().Misc()
	productsResponse, err := c.Products()
	if err != nil {
		log.Println("Can't get products error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(productsResponse)
	if err != nil {
		log.Println("Can't json.Marshal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	total := strconv.Itoa(len(productsResponse))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", total)
	fmt.Fprintf(w, string(b))
}
