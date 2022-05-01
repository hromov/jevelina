package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/base"
)

// const dsn = "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

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
	r.HandleFunc("/steps", api.StepsHandler).Methods("GET", "POST")
	r.HandleFunc("/steps/{id}", api.StepHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/products", api.ProductsHandler).Methods("GET", "POST")
	r.HandleFunc("/products/{id}", api.ProductHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/manufacturers", api.ManufacturersHandler).Methods("GET", "POST")
	r.HandleFunc("/manufacturers/{id}", api.ManufacturerHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/tags", api.TagsHandler).Methods("GET", "POST")
	r.HandleFunc("/tags/{id}", api.TagHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/tasks", api.TasksHandler).Methods("GET", "POST")
	r.HandleFunc("/tasks/{id}", api.TaskHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/tasktypes", api.TaskTypesHandler).Methods("GET", "POST")
	r.HandleFunc("/tasktypes/{id}", api.TaskTypeHandler).Methods("GET", "PUT", "DELETE")

	return r
}

func main() {
	dsn, err := os.ReadFile("_keys/db")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(dsn))
	if err := base.Init(string(dsn)); err != nil {
		log.Fatalf("Cant init data base error: %s", err.Error())
	}

	// n := 3000

	// if err := amoimport.Push_Misc("_import/amocrm_export_leads_2022-04-20.csv", n); err != nil {
	// 	log.Println(err)
	// }

	// if err := amoimport.Push_Contacts("_import/amocrm_export_contacts_2022-04-20.csv", n); err != nil {
	// 	log.Println(err)
	// }

	// if err := amoimport.Push_Leads("_import/amocrm_export_leads_2022-04-20.csv", n); err != nil {
	// 	log.Println(err)
	// }

	router := newREST()
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "X-Requested-With", "application/json", "Authorization"})

	// ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200", "https://d3qttgy7smx7mi.cloudfront.net", "https://front-dot-vorota-ua.ew.r.appspot.com", os.Getenv("ORIGIN_ALLOWED")})

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(credentials, methods, origins, headersOk)(router)))

	// log.Fatal(http.ListenAndServeTLS(":5000", "_keys/public.crt", "_keys/private.pem", handlers.CORS(credentials, methods, origins, headersOk)(router)))
}
