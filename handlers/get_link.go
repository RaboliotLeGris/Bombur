package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	"bombur/db"
)

type GetLink struct {
	Pool *pgxpool.Pool
}

func (l GetLink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("Request POST /link")

	slug, found := mux.Vars(r)["slug"]
	if !found || slug == "" {
		log.Debug("Missing link slug")
		http.Error(w, "empty slug provided", http.StatusBadRequest)
		return
	}

	link, err := db.NewLinkDAO(l.Pool).GetLink(r.Context(), slug)
	if err != nil {
		log.Debug("Failed to get link ", err)
		w.WriteHeader(404)
		return
	}

	http.Redirect(w, r, link, http.StatusMovedPermanently)
}
