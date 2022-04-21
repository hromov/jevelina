package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	amo "./pkg/amo_import"

	testdata "./pkg/test_data"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hromov/cdb"
)

// func create_users() {

// 	role := cdb.Role{
// 		Name: "User",
// 	}
// 	if r_role, err := cdb.Create(&role); err != nil {
// 		log.Println(r_role)
// 	}
// 	user := cdb.User{
// 		Name:  "Vasya A",
// 		Email: "vasya2@gmail.com",
// 	}
// 	if r_user, err := cdb.Create(&user); err != nil {
// 		log.Println(r_user)
// 	}
// }

func contactsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contacts" {
		http.NotFound(w, r)
		return
	}
	contacts, err := cdb.Contacts()
	if err != nil {
		log.Println("Can't get contacts error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(contacts)
	if err != nil {
		log.Println("Can't json.Marchal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(b))
}

func newREST() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/contacts", contactsHandler).Methods("GET")
	// r.HandleFunc("/banks", newBankHandler).Methods("POST")
	// r.HandleFunc("/banks/{id}", bankChangeHandler).Methods("PUT", "DELETE")
	return r
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	if err := cdb.Init(dsn); err != nil {
		log.Fatalf("Cant init data base error: %s", err.Error())
	}

	testdata.Fill()

	if err := amo.Push_Contacts("../backup/amocrm_export_contacts_2022-04-20.csv"); err != nil {
		log.Println(err)
	}
	// create_users()
	router := newREST()
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "X-Requested-With", "application/json"})

	// ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200", os.Getenv("ORIGIN_ALLOWED")})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(credentials, methods, origins, headersOk)(router)))
}
