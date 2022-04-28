package main

import (
	"log"
	"net/http"
	"os"

	correlationid "github.com/dmytrohridin/correlation-id"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
)

func main() {
	arg := os.Args[1]
	log.Println(arg)
	middleware := correlationid.New()
	var routerHandler http.Handler
	switch arg {
	case "mux":
		mux := http.NewServeMux()
		mux.Handle("/ping", middleware.Handle(http.HandlerFunc(defaultHandler)))
		routerHandler = mux
	case "gorilla":
		gorillaRouter := mux.NewRouter()
		gorillaRouter.Use(middleware.Handle)
		gorillaRouter.HandleFunc("/ping", defaultHandler)
		routerHandler = gorillaRouter
	case "chi":
		chiRouter := chi.NewRouter()
		chiRouter.Use(middleware.Handle)
		chiRouter.Get("/ping", defaultHandler)
		routerHandler = chiRouter
	default:
		panic("Unexpected argument")
	}

	log.Println("Listening...")
	http.ListenAndServe(":3000", routerHandler)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	correlationId := correlationid.FromContext(r.Context())
	w.Write([]byte(correlationId))
}
