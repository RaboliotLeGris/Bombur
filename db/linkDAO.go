package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/thanhpk/randstr"
)

type LinkDAO struct {
	conn *pgxpool.Pool
}

func NewLinkDAO(conn *pgxpool.Pool) *LinkDAO {
	return &LinkDAO{conn: conn}
}

func (l LinkDAO) CreateLink(ctx context.Context, link string, expire *time.Duration) (string, error) {
	if link == "" {
		return "", fmt.Errorf("empty link string")
	}

	slug := randstr.String(8)
	if expire != nil {
		return slug, l.createLinkWithExpiration(ctx, link, slug, expire)
	}

	return slug, l.createLink(ctx, link, slug)
}

func (l LinkDAO) createLink(ctx context.Context, link, slug string) error {
	commandTag, err := l.conn.Exec(ctx, "INSERT INTO link (link, slug) VALUES ($1, $2);", link, slug)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("failed to create new link")
	}
	return nil
}

func (l LinkDAO) createLinkWithExpiration(ctx context.Context, link, slug string, expire *time.Duration) error {
	// TODO(Rabo): expire must be current time of request + duration
	now := time.Now()
	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		return err
	}
	expireAt := now.In(loc).Add(*expire)

	commandTag, err := l.conn.Exec(ctx, "INSERT INTO link (link, slug, expire) VALUES ($1, $2, $3);", link, slug, expireAt)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("failed to create new link")
	}

	return nil
}

func (l LinkDAO) GetLink(ctx context.Context, slug string) (string, error) {
	var link string
	err := l.conn.QueryRow(ctx, "SELECT link FROM link WHERE slug=$1;", slug).Scan(&link)
	return link, err
}
