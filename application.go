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
	"github.com/hromov/jevelina/routes"
)

// const dsn = "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

func newREST() *mux.Router {
	r := mux.NewRouter()
	r = routes.UserRoutes(r)
	r = routes.AdminRoutes(r)
	r = routes.FinRoutes(r)
	r = routes.FilesRoutes(r)
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

	// const leads = "_import/amocrm_export_leads_2022-04-20.csv"
	// const contacts = "_import/amocrm_export_contacts_2022-04-20.csv"
	// if err := amoimport.Import(base.GetDB().DB, leads, contacts, 1500); err != nil {
	// 	log.Fatalf("Can't import error: %s", err.Error())
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

	router := newREST()
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "X-Requested-With", "application/json", "Authorization"})

	// ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200", "https://d3qttgy7smx7mi.cloudfront.net", "https://front-dot-vorota-ua.ew.r.appspot.com", os.Getenv("ORIGIN_ALLOWED")})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handlers.CORS(credentials, methods, origins, headersOk)(router)); err != nil {
		log.Fatal(err)
	}
	// log.Fatal(http.ListenAndServeTLS(":5000", "_keys/public.crt", "_keys/private.pem", handlers.CORS(credentials, methods, origins, headersOk)(router)))
}
