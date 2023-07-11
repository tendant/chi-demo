package dbconn

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/exp/slog"
)

type DBConnSettings struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func OpenDBConn(driver string, dsn string, settings DBConnSettings) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		slog.Error("Failed open database", "driver", driver)
	} else {
		if settings.ConnMaxLifetime > 0 {
			db.SetConnMaxLifetime(settings.ConnMaxLifetime)
		}
		if settings.ConnMaxIdleTime > 0 {
			db.SetConnMaxIdleTime(settings.ConnMaxIdleTime)
		}
		if settings.MaxIdleConns > 0 {
			db.SetMaxIdleConns(settings.MaxIdleConns)
		}
		if settings.MaxOpenConns > 0 {
			db.SetMaxOpenConns(settings.MaxOpenConns)
		}
	}
	return db, err
}

func CheckDBConn(dbConn *sql.DB) (bool, error) {
	if dbConn != nil {
		err := dbConn.Ping()
		if err != nil {
			return false, errors.New("Failed Ping Database Connection!")
		} else {
			return true, nil
		}
	} else {
		return false, errors.New("No Database Connection!")
	}
}
