package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

func InitDB(db_uri string) error {
	if db_uri == "" {
		return fmt.Errorf("empty DB URI")
	}

	log.Info("DB - initDB - Checking database status")
	conn, err := pgx.Connect(context.Background(), db_uri)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	needInit, err := isInitRequired(conn)
	if err != nil {
		return err
	}

	if needInit {
		log.Info("DB - initDB - Initializing database")
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return err
		}

		// Create the schema
		if _, err = tx.Exec(context.Background(), Schema); err != nil {
			return err
		}

		if err = tx.Commit(context.Background()); err != nil {
			return err
		}

	}
	return nil
}

func isInitRequired(conn *pgx.Conn) (bool, error) {
	var version int64
	if err := conn.QueryRow(context.Background(), "SELECT * FROM meta_info;").Scan(&version); err != nil {
		if err.Error() == "ERROR: relation \"meta_info\" does not exist (SQLSTATE 42P01)" {
			return true, nil
		}
		return false, err
	}
	log.Debug("DB VERSION: ", version)
	return false, nil
}
