package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/hromov/jevelina/config"
	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/finances"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/domain/misc/files"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/http/rest"
	"github.com/hromov/jevelina/http/rest/auth"
	"github.com/hromov/jevelina/storage/gcloud"
	"github.com/hromov/jevelina/storage/mysql"
	"github.com/hromov/jevelina/useCases/orders"
	"github.com/hromov/jevelina/useCases/tasks"
	"github.com/hromov/jevelina/utils/events"
	tokenvalidator "github.com/hromov/jevelina/utils/tokenValidator"
)

func main() {
	cfg := config.Get()
	log.Println(cfg)
	dns, err := os.ReadFile(cfg.Dsn)
	if err != nil {
		log.Fatal(err)
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
	fin := finances.NewService(storage)
	es := events.NewService(storage)
	gc, err := gcloud.NewService(context.Background(), cfg.BucketName)
	if err != nil {
		log.Println("Can't create google cloud client error: ", err.Error())
	}
	fs := files.NewService(storage, gc)
	ordersService := orders.NewService(cs, ls, us, ts)
	router := rest.InitRouter(us, cs, ls, ordersService, ms, ts, as, fs, fin, es)
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
