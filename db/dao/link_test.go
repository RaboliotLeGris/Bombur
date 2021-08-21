package dao_test

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/franela/goblin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"

	"bombur/db"
	"bombur/db/dao"
)

func Test_Link_DAO(t *testing.T) {
	DB_URI := os.Getenv("BOMBUR_DB_URI")
	pool, err := pgxpool.Connect(context.Background(), DB_URI)
	require.NoError(t, err)
	defer pool.Close()

	g := Goblin(t)
	g.Describe("Link >", func() {
		g.Before(func() {
			require.NoError(t, db.InitDB(DB_URI))
		})

		g.It("Create link without expiration", func() {
			linkDAO := dao.NewLinkDAO(pool)
			slug, err := linkDAO.CreateLink(context.Background(), "https://raboland.fr", nil)
			require.NoError(t, err)
			var count int
			err = pool.QueryRow(context.Background(), "SELECT COUNT(*) from link WHERE slug=$1;", slug).Scan(&count)
			require.NoError(t, err)
			require.Equal(t, 1, count)
		})

		g.It("Create link with expiration", func() {
			linkDAO := dao.NewLinkDAO(pool)
			oneMinute := time.Minute

			slug, err := linkDAO.CreateLink(context.Background(), "https://raboland.fr", &oneMinute)
			require.NoError(t, err)

			var count int
			var expireAt time.Time
			err = pool.QueryRow(context.Background(), "SELECT COUNT(*), expire from link WHERE slug=$1 GROUP BY expire;", slug).Scan(&count, &expireAt)
			require.NoError(t, err)
			require.Equal(t, 1, count)
		})

		g.It("Fail to create link with empty string", func() {
			linkDAO := dao.NewLinkDAO(pool)
			_, err := linkDAO.CreateLink(context.Background(), "", nil)
			require.Error(t, err)
		})

		g.It("Get link without expiration", func() {
			givenLink := "https://raboland.fr"
			linkDAO := dao.NewLinkDAO(pool)

			slug, err := linkDAO.CreateLink(context.Background(), givenLink, nil)
			require.NoError(t, err)

			gotLink, err := linkDAO.GetLink(context.Background(), slug)
			require.NoError(t, err)
			require.Equal(t, givenLink, gotLink)
		})

		g.It("Get link with expiration and valid", func() {
			oneMinute := time.Minute
			givenLink := "https://raboland.fr"
			linkDAO := dao.NewLinkDAO(pool)

			slug, err := linkDAO.CreateLink(context.Background(), givenLink, &oneMinute)
			require.NoError(t, err)

			gotLink, err := linkDAO.GetLink(context.Background(), slug)
			require.NoError(t, err)
			require.Equal(t, givenLink, gotLink)
		})

		g.It("Get link with expiration and expired", func() {
			g.Timeout(10 * time.Second)
			linkDAO := dao.NewLinkDAO(pool)
			oneSecond := time.Second

			slug, err := linkDAO.CreateLink(context.Background(), "https://raboland.fr", &oneSecond)
			require.NoError(t, err)
			time.Sleep(2 * time.Second)

			_, err = linkDAO.GetLink(context.Background(), slug)
			require.Error(t, err)
		})
	})
}
