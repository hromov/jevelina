package fin_api

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/auth"
)

func Rest(r *mux.Router) *mux.Router {
	r.HandleFunc("/wallets", auth.UserCheck(WalletsHandler)).Methods("GET")
	r.HandleFunc("/wallets", auth.AdminCheck(WalletsHandler)).Methods("POST")
	r.HandleFunc("/wallets/{id}", auth.AdminCheck(WalletHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/wallets/{id}/close", auth.AdminCheck(CloseWalletHandler)).Methods("GET")
	r.HandleFunc("/wallets/{id}/open", auth.AdminCheck(OpenWalletHandler)).Methods("GET")

	r.HandleFunc("/transfers", auth.UserCheck(TransfersHandler)).Methods("GET", "POST")
	r.HandleFunc("/transfers/{id}", auth.UserCheck(TransferHandler)).Methods("PUT")
	r.HandleFunc("/transfers/{id}", auth.AdminCheck(TransferHandler)).Methods("DELETE")
	r.HandleFunc("/transfers/{id}/complete", auth.AdminCheck(CompleteTransferHandler)).Methods("GET")
	return r
}
