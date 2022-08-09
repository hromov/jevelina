package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hromov/jevelina/domain/finances"
)

// func CloseWalletHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	ID, err := strconv.ParseUint(vars["id"], 10, 32)
// 	if err != nil {
// 		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	fin := mysql.Finance()
// 	if r.Method == "GET" {
// 		if err := fin.ChangeWalletState(uint16(ID), true); err != nil {
// 			log.Printf("Can't save item with ID = %d. Error: %s", ID, err.Error())
// 			http.Error(w, http.StatusText(http.StatusInternalServerError),
// 				http.StatusInternalServerError)
// 		}
// 	}
// }

// func OpenWalletHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	ID, err := strconv.ParseUint(vars["id"], 10, 32)
// 	if err != nil {
// 		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	fin := mysql.Finance()
// 	if r.Method == "GET" {
// 		if err := fin.ChangeWalletState(uint16(ID), false); err != nil {
// 			log.Printf("Can't save item with ID = %d. Error: %s", ID, err.Error())
// 			http.Error(w, http.StatusText(http.StatusInternalServerError),
// 				http.StatusInternalServerError)
// 		}
// 	}
// }

func ChangeWalletState(f finances.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}
		state := struct {
			Closed bool
		}{}
		if err = json.NewDecoder(r.Body).Decode(&state); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := f.ChangeWalletState(r.Context(), uint16(id), state.Closed); err != nil {
			log.Printf("Can't change wallet with ID = %d. Error: %s", id, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func Wallet(f finances.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "PUT":
			wallet := finances.Wallet{}
			if err = json.NewDecoder(r.Body).Decode(&wallet); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uint64(wallet.ID) != id {
				http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", id, wallet.ID), http.StatusBadRequest)
				return
			}

			if err = f.ChangeWalletName(r.Context(), wallet.ID, wallet.Name); err != nil {
				log.Printf("Can't save item with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		case "DELETE":
			if err := f.DeleteWallet(r.Context(), uint16(id)); err != nil {
				log.Printf("Can't delete item with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
	}
}

func Wallets(f finances.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			wallet := finances.Wallet{}
			if err := json.NewDecoder(r.Body).Decode(&wallet); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			wallet, err := f.CreateWallet(r.Context(), wallet)
			if err != nil {
				log.Printf("Can't create wallet. Error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}

			json.NewEncoder(w).Encode(wallet)
			return
		}

		wallets, err := f.ListWallets(r.Context())
		if err != nil {
			log.Println("Can't get wallets error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(wallets)
	}
}
