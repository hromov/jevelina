package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/routes"
)

// const dsn = "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
const bucketName = "jevelina"

func main() {
	dsn, err := os.ReadFile("_keys/db_google2")
	if err != nil {
		log.Fatal(err)
	}

	db, err := cdb.OpenAndInit(string(dsn))
	if err != nil {
		log.Fatalf("Cant open and init data base error: %s", err.Error())
	}
	db.SetBucket(bucketName)

	router := routes.Base()
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
