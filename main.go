package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/base"
)

const dsn = "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

func newREST() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/contacts", api.ContactsHandler).Methods("GET", "POST")
	r.HandleFunc("/contacts/{id}", api.ContactHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/leads", api.LeadsHandler).Methods("GET", "POST")
	r.HandleFunc("/leads/{id}", api.LeadHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/users", api.UsersHandler).Methods("GET", "POST")
	r.HandleFunc("/users/{id}", api.UserHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/sources", api.SourcesHandler).Methods("GET", "POST")
	r.HandleFunc("/sources/{id}", api.SourceHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/roles", api.RolesHandler).Methods("GET", "POST")
	r.HandleFunc("/roles/{id}", api.RoleHandler).Methods("GET", "PUT", "DELETE")

	return r
}

func main() {
	if err := base.Init(dsn); err != nil {
		log.Fatalf("Cant init data base error: %s", err.Error())
	}
	// testdata.Fill()

	// if err := amoimport.Push_Contacts("/home/serhii/git/backup/amocrm_export_contacts_2022-04-20.csv"); err != nil {
	// 	log.Println(err)
	// }

	// if err := amoimport.Push_Leads("/home/serhii/git/backup/amocrm_export_leads_2022-04-20.csv"); err != nil {
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
