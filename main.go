package main

import (
	"log"
	"net/http"
	"os"

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

func newREST() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/contacts", contacts.contactsHandler).Methods("GET")
	// r.HandleFunc("/banks", newBankHandler).Methods("POST")
	// r.HandleFunc("/banks/{id}", bankChangeHandler).Methods("PUT", "DELETE")
	return r
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
	if err := cdb.Init(dsn); err != nil {
		log.Fatalf("Cant init data base error: %s", err.Error())
	}

	// testdata.Fill()

	// if err := amo.Push_Contacts("../backup/amocrm_export_contacts_2022-04-20.csv"); err != nil {
	// 	log.Println(err)
	// }

	// create_users()
	router := newREST()
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "X-Requested-With", "application/json"})

	// ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200", os.Getenv("ORIGIN_ALLOWED")})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(credentials, methods, origins, headersOk)(router)))
}
