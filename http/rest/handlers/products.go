package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hromov/jevelina/domain/misc"
	"gorm.io/gorm"
)

func Product(ms misc.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			product, err := ms.GetProduct(r.Context(), uint32(id))
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
			_ = json.NewEncoder(w).Encode(product)
			return
		case "PUT":
			product := misc.Product{}
			if err = json.NewDecoder(r.Body).Decode(&product); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uint64(product.ID) != id {
				http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", id, product.ID), http.StatusBadRequest)
				return
			}

			if err := ms.UpdateProduct(r.Context(), product); err != nil {
				log.Printf("Can't update product with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		case "DELETE":
			if err := ms.DeleteProduct(r.Context(), uint32(id)); err != nil {
				log.Printf("Can't delete product with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
	}
}

func Products(ms misc.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			product := misc.Product{}
			if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			product, err := ms.CreateProduct(r.Context(), product)
			if err != nil {
				log.Println("Can't create product error: ", err.Error())
				http.Error(w, "Can't create product error", http.StatusInternalServerError)
			}
			return
		}

		products, err := ms.ListProducts(r.Context())
		if err != nil {
			log.Println("Can't get products error: ", err)
			http.Error(w, "Can't get products list", http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(products)
	}
}
