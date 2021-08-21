package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	"github.com/raboliotlegris/bombur/db/dao"
)

type GetLink struct {
	Pool *pgxpool.Pool
}

type failureResponse struct {
	Reason string `json:"reason"`
}

func (l GetLink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slug, found := mux.Vars(r)["slug"]
	if !found || slug == "" {
		log.Debug("Missing link slug")
		http.Error(w, "empty slug provided", http.StatusBadRequest)
		return
	}
	log.Debug("Request GET /s/", slug)

	link, err := dao.NewLinkDAO(l.Pool).GetLink(r.Context(), slug)
	if err != nil {
		log.Debug("Failed to get link ", err)
		_ = json.NewEncoder(w).Encode(failureResponse{Reason: "Link doesn't exist or have expired"})
		w.Header().Add("Content-type", "application/json")
		return
	}

	// Permanent redirect seems to use browser cache
	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}
