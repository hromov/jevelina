package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/auth"
	"github.com/hromov/jevelina/base"
)

// const dsn = "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

func usersRest(r *mux.Router) *mux.Router {
	r.HandleFunc("/contacts", auth.UserCheck(api.ContactsHandler)).Methods("GET", "POST")
	r.HandleFunc("/contacts/{id}", auth.UserCheck(api.ContactHandler)).Methods("GET", "PUT", "DELETE")
	// r.HandleFunc("/leads", auth.UserCheck(api.LeadsHandler)).Methods("GET", "POST")
	r.HandleFunc("/leads", api.LeadsHandler).Methods("GET", "POST")
	r.HandleFunc("/leads/{id}", auth.UserCheck(api.LeadHandler)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/users", auth.UserCheck(api.UsersHandler)).Methods("GET")
	r.HandleFunc("/users/{id}", auth.UserCheck(api.UserHandler)).Methods("GET")
	r.HandleFunc("/sources", auth.UserCheck(api.SourcesHandler)).Methods("GET")
	r.HandleFunc("/sources/{id}", auth.UserCheck(api.SourceHandler)).Methods("GET")
	r.HandleFunc("/roles", auth.UserCheck(api.RolesHandler)).Methods("GET")
	r.HandleFunc("/roles/{id}", auth.UserCheck(api.RoleHandler)).Methods("GET")
	r.HandleFunc("/steps", auth.UserCheck(api.StepsHandler)).Methods("GET")
	r.HandleFunc("/steps/{id}", auth.UserCheck(api.StepHandler)).Methods("GET")
	r.HandleFunc("/products", auth.UserCheck(api.ProductsHandler)).Methods("GET")
	r.HandleFunc("/products/{id}", auth.UserCheck(api.ProductHandler)).Methods("GET")
	r.HandleFunc("/manufacturers", auth.UserCheck(api.ManufacturersHandler)).Methods("GET")
	r.HandleFunc("/manufacturers/{id}", auth.UserCheck(api.ManufacturerHandler)).Methods("GET")
	r.HandleFunc("/tags", auth.UserCheck(api.TagsHandler)).Methods("GET")
	r.HandleFunc("/tags/{id}", auth.UserCheck(api.TagHandler)).Methods("GET")
	r.HandleFunc("/tasks", auth.UserCheck(api.TasksHandler)).Methods("GET", "POST")
	r.HandleFunc("/tasks/{id}", auth.UserCheck(api.TaskHandler)).Methods("GET", "PUT")
	r.HandleFunc("/tasktypes", auth.UserCheck(api.TaskTypesHandler)).Methods("GET")
	r.HandleFunc("/tasktypes/{id}", auth.UserCheck(api.TaskTypeHandler)).Methods("GET")
	return r
}

func adminRest(r *mux.Router) *mux.Router {
	r.HandleFunc("/users", auth.AdminCheck(api.UsersHandler)).Methods("POST")
	r.HandleFunc("/users/{id}", auth.AdminCheck(api.UserHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/sources", auth.AdminCheck(api.SourcesHandler)).Methods("POST")
	r.HandleFunc("/sources/{id}", auth.AdminCheck(api.SourceHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/roles", auth.AdminCheck(api.RolesHandler)).Methods("POST")
	r.HandleFunc("/roles/{id}", auth.AdminCheck(api.RoleHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/steps", auth.AdminCheck(api.StepsHandler)).Methods("POST")
	r.HandleFunc("/steps/{id}", auth.AdminCheck(api.StepHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/products", auth.AdminCheck(api.ProductsHandler)).Methods("POST")
	r.HandleFunc("/products/{id}", auth.AdminCheck(api.ProductHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/manufacturers", auth.AdminCheck(api.ManufacturersHandler)).Methods("POST")
	r.HandleFunc("/manufacturers/{id}", auth.AdminCheck(api.ManufacturerHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/tags", auth.AdminCheck(api.TagsHandler)).Methods("POST")
	r.HandleFunc("/tags/{id}", auth.AdminCheck(api.TagHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/tasks/{id}", auth.AdminCheck(api.TaskHandler)).Methods("DELETE")
	r.HandleFunc("/tasktypes", auth.AdminCheck(api.TaskTypesHandler)).Methods("POST")
	r.HandleFunc("/tasktypes/{id}", auth.AdminCheck(api.TaskTypeHandler)).Methods("PUT", "DELETE")
	return r
}

func finRest(r *mux.Router) *mux.Router {
	r.HandleFunc("/wallets", auth.UserCheck(api.WalletsHandler)).Methods("GET")
	r.HandleFunc("/wallets", auth.AdminCheck(api.WalletsHandler)).Methods("POST")
	r.HandleFunc("/wallets/{id}", auth.AdminCheck(api.WalletHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/wallets/{id}/close", auth.AdminCheck(api.CloseWalletHandler)).Methods("PUT")
	r.HandleFunc("/wallets/{id}/open", auth.AdminCheck(api.OpenWalletHandler)).Methods("PUT")
	return r
}

func newREST() *mux.Router {
	r := mux.NewRouter()
	r = usersRest(r)
	r = adminRest(r)
	r.HandleFunc("/usercheck", auth.UserCheckHandler).Methods("GET")
	r.HandleFunc("/orders", api.OrderHandler).Methods("POST")
	return r
}

func main() {
	dsn, err := os.ReadFile("_keys/db_local")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(dsn))
	if err := base.Init(string(dsn)); err != nil {
		log.Fatalf("Cant init data base error: %s", err.Error())
	}

	// n := 1500

	// if err := amoimport.Push_Misc("_import/amocrm_export_leads_2022-04-20.csv", n); err != nil {
	// 	log.Println(err)
	// }

	// adminRoleID := uint8(1)
	// admin := &models.User{
	// 	Name:   "Admin User",
	// 	Email:  "melifarowow@gmail.com",
	// 	Hash:   "melifarowow@gmail.com",
	// 	RoleID: &adminRoleID,
	// }
	// if err := base.GetDB().DB.Omit(clause.Associations).Create(admin).Error; err != nil {
	// 	log.Fatalf("Can't create admin error: %s", err.Error())
	// }

	// userRoleID := uint8(2)
	// user := &models.User{
	// 	Name:   "Random User",
	// 	Email:  "random@random.org",
	// 	Hash:   "random@random.org",
	// 	RoleID: &userRoleID,
	// }
	// if err := base.GetDB().DB.Omit(clause.Associations).Create(user).Error; err != nil {
	// 	log.Printf("Can't create random error: %s", err.Error())
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
