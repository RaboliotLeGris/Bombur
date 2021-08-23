package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

func Create_router(pool *pgxpool.Pool) *mux.Router {
	log.Info("Creating routers")
	// Routes order creation matter.
	r := mux.NewRouter()

	r.PathPrefix("/s/{slug}").Handler(GetLink{Pool: pool}).Methods("GET")
	r.PathPrefix("/link").Handler(CreateLink{Pool: pool}).Methods("POST")
	r.PathPrefix("/").Handler(StaticHandler{StaticPath: "static", IndexPath: "index.html"})

	return r
}

func Launch(router *mux.Router, port uint64) error {
	log.Info("Launching HTTP server")

	srv := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("0.0.0.0:%v", port),
	}

	return srv.ListenAndServe()
}
