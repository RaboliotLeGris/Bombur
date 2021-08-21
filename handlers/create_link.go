package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	"bombur/db"
)

type CreateLink struct {
	Pool   *pgxpool.Pool
	UseTLS bool
}

type LinkCreatePayload struct {
	Link   string `json:"link"`
	Expire string `json:"expire"`
}

type LinkResponse struct {
	Link string `json:"link"`
}

func (l CreateLink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("Request POST /link")

	var linkToCreate LinkCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&linkToCreate); err != nil {
		log.Debug("Failed to unmarshal payload with error: ", err)
		w.WriteHeader(400)
		return
	}

	if linkToCreate.Link == "" {
		log.Debug("Empty link value")
		w.WriteHeader(400)
		return
	}

	var expireIn *time.Duration
	if linkToCreate.Expire != "" {
		got, err := time.ParseDuration(linkToCreate.Expire)
		if err != nil {
			log.Debug("Failed to parse duration ", linkToCreate.Expire)
			w.WriteHeader(400)
			return
		}
		expireIn = &got
	}

	slug, err := db.NewLinkDAO(l.Pool).CreateLink(r.Context(), linkToCreate.Link, expireIn)
	if err != nil {
		log.Debug("Failed to create link ", err)
		w.WriteHeader(500)
		return
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		log.Debug("Missing Origin Header ")
		w.WriteHeader(500)
		return
	}

	link := fmt.Sprintf("%s/s/%s", origin, slug)
	log.Debug("Shorten URL: ", link)

	_ = json.NewEncoder(w).Encode(LinkResponse{Link: link})
	w.Header().Add("Content-type", "application/json")
}
