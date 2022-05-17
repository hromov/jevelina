package routes

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api/fin_api"
	"github.com/hromov/jevelina/auth"
)

func FinRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/wallets", auth.UserCheck(fin_api.WalletsHandler)).Methods("GET")
	r.HandleFunc("/wallets", auth.AdminCheck(fin_api.WalletsHandler)).Methods("POST")
	r.HandleFunc("/wallets/{id}", auth.AdminCheck(fin_api.WalletHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/wallets/{id}/close", auth.AdminCheck(fin_api.CloseWalletHandler)).Methods("GET")
	r.HandleFunc("/wallets/{id}/open", auth.AdminCheck(fin_api.OpenWalletHandler)).Methods("GET")

	r.HandleFunc("/transfers", auth.UserCheck(fin_api.TransfersHandler)).Methods("GET", "POST")
	r.HandleFunc("/transfers/{id}", auth.UserCheck(fin_api.TransferHandler)).Methods("PUT")
	r.HandleFunc("/transfers/{id}", auth.AdminCheck(fin_api.TransferHandler)).Methods("DELETE")
	r.HandleFunc("/transfers/{id}/complete", auth.AdminCheck(fin_api.CompleteTransferHandler)).Methods("GET")
	r.HandleFunc("/categories", auth.UserCheck(fin_api.CategoriesHandler)).Methods("GET")

	r.HandleFunc("/analytics/categories", auth.UserCheck(fin_api.CategoriesSumHandler)).Methods("GET")
	return r
}
