package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/hromov/jevelina/config"
	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/http/rest"
	"github.com/hromov/jevelina/storage/mysql"
)

// const dsn = "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
const bucketName = "jevelina"

func main() {
	cfg := config.Get()
	log.Println(cfg)
	dsn, err := os.ReadFile(cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}

	db, err := mysql.OpenAndInit(string(dsn))
	if err != nil {
		log.Fatalf("Cant open and init data base error: %s", err.Error())
	}
	db.SetBucket(bucketName)
	//TODO: repo
	us := users.NewService(mysql.Misc())
	cs := contacts.NewService(mysql.Contacts())
	ls := leads.NewService(mysql.Leads())
	router := rest.InitRouter(us, cs, ls)
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "X-Requested-With", "application/json", "Authorization"})

	origins := handlers.AllowedOrigins([]string{"http://localhost:4200", "https://front-dot-vorota-ua.ew.r.appspot.com", os.Getenv("ORIGIN_ALLOWED")})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handlers.CORS(credentials, methods, origins, headersOk)(router)); err != nil {
		log.Fatal(err)
	}
}
