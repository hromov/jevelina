package fin_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
)

func CloseWalletHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	fin := cdb.Finance()
	if r.Method == "GET" {
		if err := fin.ChangeWalletState(uint16(ID), true); err != nil {
			log.Printf("Can't save item with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func OpenWalletHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	fin := cdb.Finance()
	if r.Method == "GET" {
		if err := fin.ChangeWalletState(uint16(ID), false); err != nil {
			log.Printf("Can't save item with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func WalletHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Println(ID)
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}
	fin := cdb.Finance()
	switch r.Method {
	case "PUT":
		var wallet *models.Wallet
		if err = json.NewDecoder(r.Body).Decode(&wallet); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(wallet.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, wallet.ID), http.StatusBadRequest)
			return
		}

		if err = fin.ChangeWalletName(wallet.ID, wallet.Name); err != nil {
			log.Printf("Can't save item with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	case "DELETE":
		if err := fin.DeleteWallet(uint16(ID)); err != nil {
			log.Printf("Can't delete item with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	}

}

func WalletsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/wallets" {
		http.NotFound(w, r)
		return
	}
	fin := cdb.Finance()
	if r.Method == "POST" {
		item := new(models.Wallet)
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item, err := fin.CreateWallet(item)
		if err != nil {
			log.Printf("Can't create wallet. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(item)
		if err != nil {
			log.Println("Can't json.Marshal(wallet) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
	}

	wallets, err := fin.ListWallets(nil)
	if err != nil {
		log.Println("Can't get wallets error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	b, err := json.Marshal(wallets)
	if err != nil {
		log.Println("Can't json.Marshal(wallets) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	total := strconv.Itoa(len(wallets))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", total)
	fmt.Fprint(w, string(b))
}
