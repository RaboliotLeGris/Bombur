package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/franela/goblin"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"bombur/db"
	"bombur/handlers"
)

func Test_Routes(t *testing.T) {
	DB_URI := os.Getenv("BOMBUR_DB_URI")
	pool, err := pgxpool.Connect(context.Background(), DB_URI)
	require.NoError(t, err)
	defer pool.Close()

	g := Goblin(t)
	g.Describe("Routes >", func() {
		g.Before(func() {
			// If logging is required
			log.SetLevel(log.DebugLevel)
			require.NoError(t, db.InitDB(DB_URI))
		})

		g.Describe("POST /link >", func() {
			g.It("Works with no expiration", func() {
				givenLink := "https://raboland.fr"

				w := httptest.NewRecorder()
				r := mux.NewRouter()
				r.PathPrefix("/link").Handler(handlers.CreateLink{Pool: pool}).Methods("POST")

				rawBuf, err := json.Marshal(handlers.LinkCreatePayload{Link: givenLink})
				require.NoError(t, err)
				buffer := bytes.NewBuffer(rawBuf)

				req := httptest.NewRequest("POST", "/link", buffer)
				req.Header.Add("Origin", "http://SomeOrigin")
				r.ServeHTTP(w, req)

				require.Equal(t, 200, w.Code)
				var response handlers.LinkResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				require.True(t, strings.Contains(response.Link, "http://SomeOrigin/s/"))
				require.Equal(t, 28, len(response.Link))
			})

			g.It("Works with expiration", func() {
				givenLink := "https://raboland.fr"

				w := httptest.NewRecorder()
				r := mux.NewRouter()
				r.PathPrefix("/link").Handler(handlers.CreateLink{Pool: pool}).Methods("POST")

				rawBuf, err := json.Marshal(handlers.LinkCreatePayload{Link: givenLink, Expire: "10m"})
				require.NoError(t, err)
				buffer := bytes.NewBuffer(rawBuf)

				req := httptest.NewRequest("POST", "/link", buffer)
				req.Header.Add("Origin", "http://SomeOrigin")
				r.ServeHTTP(w, req)

				require.Equal(t, 200, w.Code)
				var response handlers.LinkResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				require.True(t, strings.Contains(response.Link, "http://SomeOrigin/s/"))
				require.Equal(t, 28, len(response.Link))
			})

			g.It("Failed if link is missing", func() {
				w := httptest.NewRecorder()
				r := mux.NewRouter()
				r.PathPrefix("/link").Handler(handlers.CreateLink{Pool: pool}).Methods("POST")

				rawBuf, err := json.Marshal(handlers.LinkCreatePayload{})
				require.NoError(t, err)
				buffer := bytes.NewBuffer(rawBuf)

				req := httptest.NewRequest("POST", "/link", buffer)
				req.Header.Add("Origin", "http://SomeOrigin")
				r.ServeHTTP(w, req)

				require.Equal(t, 400, w.Code)
			})

		})

		g.Describe("GET /s/{slug} >", func() {
			g.It("GET link without expiration", func() {
				// Creating routes
				givenLink := "https://raboland.fr"
				r := mux.NewRouter()
				r.PathPrefix("/s/{slug}").Handler(handlers.GetLink{Pool: pool}).Methods("GET")
				r.PathPrefix("/link").Handler(handlers.CreateLink{Pool: pool}).Methods("POST")

				// Creating link
				rawBuf, err := json.Marshal(handlers.LinkCreatePayload{Link: givenLink})
				require.NoError(t, err)
				buffer := bytes.NewBuffer(rawBuf)
				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/link", buffer)
				req.Header.Add("Origin", "http://SomeOrigin")
				r.ServeHTTP(w, req)

				require.Equal(t, 200, w.Code)
				var response handlers.LinkResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Checking the link exists
				w = httptest.NewRecorder()
				req = httptest.NewRequest("GET", fmt.Sprintf("/s/%s", response.Link[20:]), nil)
				req.Header.Add("Origin", "http://SomeOrigin")
				r.ServeHTTP(w, req)

				require.Equal(t, 307, w.Code)
				require.Equal(t, givenLink, w.Header().Get("Location"))
			})

			g.It("GET link with expiration", func() {
				// Creating routes
				givenLink := "https://raboland.fr"
				r := mux.NewRouter()
				r.PathPrefix("/s/{slug}").Handler(handlers.GetLink{Pool: pool}).Methods("GET")
				r.PathPrefix("/link").Handler(handlers.CreateLink{Pool: pool}).Methods("POST")

				// Creating link
				rawBuf, err := json.Marshal(handlers.LinkCreatePayload{Link: givenLink, Expire: "1m"})
				require.NoError(t, err)
				buffer := bytes.NewBuffer(rawBuf)
				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/link", buffer)
				req.Header.Add("Origin", "http://SomeOrigin")
				r.ServeHTTP(w, req)

				require.Equal(t, 200, w.Code)
				var response handlers.LinkResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Checking the link exists
				w = httptest.NewRecorder()
				req = httptest.NewRequest("GET", fmt.Sprintf("/s/%s", response.Link[20:]), nil)
				req.Header.Add("Origin", "http://SomeOrigin")
				r.ServeHTTP(w, req)

				require.Equal(t, 307, w.Code)
				require.Equal(t, givenLink, w.Header().Get("Location"))
			})

			g.It("Failed to GET link due to expiration", func() {
				// Creating routes
				givenLink := "https://raboland.fr"
				r := mux.NewRouter()
				r.PathPrefix("/s/{slug}").Handler(handlers.GetLink{Pool: pool}).Methods("GET")
				r.PathPrefix("/link").Handler(handlers.CreateLink{Pool: pool}).Methods("POST")

				// Creating link
				rawBuf, err := json.Marshal(handlers.LinkCreatePayload{Link: givenLink, Expire: "1s"})
				require.NoError(t, err)
				buffer := bytes.NewBuffer(rawBuf)
				w := httptest.NewRecorder()
				req := httptest.NewRequest("POST", "/link", buffer)
				req.Header.Add("Origin", "http://SomeOrigin")
				r.ServeHTTP(w, req)

				require.Equal(t, 200, w.Code)
				var response handlers.LinkResponse
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				time.Sleep(2 * time.Second)

				// Checking the link exists
				w = httptest.NewRecorder()
				req = httptest.NewRequest("GET", fmt.Sprintf("/s/%s", response.Link[20:]), nil)
				req.Header.Add("Origin", "http://SomeOrigin")
				r.ServeHTTP(w, req)

				require.Equal(t, 200, w.Code)
				require.Equal(t, "{\"reason\":\"Link doesn't exist or have expired\"}\n", w.Body.String())
			})
		})

		g.Describe("GET static files >", func() {
			g.It("GET /", func() {
				w := httptest.NewRecorder()
				r := mux.NewRouter()
				r.PathPrefix("/").Handler(handlers.StaticHandler{StaticPath: "static", IndexPath: "index.html"})

				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Add("Origin", "http://SomeOrigin")
				r.ServeHTTP(w, req)

				require.Equal(t, 200, w.Code)
			})
		})
	})
}
