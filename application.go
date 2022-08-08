package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/hromov/jevelina/config"
	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/domain/misc/files"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/http/rest"
	"github.com/hromov/jevelina/http/rest/auth"
	"github.com/hromov/jevelina/storage/gcloud"
	"github.com/hromov/jevelina/storage/mysql"
	tokenvalidator "github.com/hromov/jevelina/tokenValidator"
	"github.com/hromov/jevelina/useCases/orders"
	"github.com/hromov/jevelina/useCases/tasks"
)

// const dns = "root:password@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	cfg := config.Get()
	log.Println(cfg)
	dns, err := os.ReadFile(cfg.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: remove after all transition finished on storage and services
	if _, err = mysql.OpenAndInit(string(dns)); err != nil {
		log.Fatalf("Cant open and init data base error: %s", err.Error())
	}

	storage, err := mysql.NewStorage(string(dns))
	if err != nil {
		log.Fatal("Can't init storage error: ", err.Error())
	}

	cs := contacts.NewService(storage)
	ls := leads.NewService(storage)
	us := users.NewService(storage)
	ts := tasks.NewService(storage)
	ms := misc.Service(storage)
	tv := tokenvalidator.NewService()
	as := auth.NewService(us, tv)
	gc, err := gcloud.NewService(context.Background(), cfg.BucketName)
	if err != nil {
		log.Println("Can't create google cloud client error: ", err.Error())
	}
	fs := files.NewService(storage, gc)
	ordersService := orders.NewService(cs, ls, us, ts)
	router := rest.InitRouter(us, cs, ls, ordersService, ms, ts, as, fs)
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
