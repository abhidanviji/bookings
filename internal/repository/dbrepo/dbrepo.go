package dbrepo

import (
	"database/sql"

	"github.com/abhidanviji/bookings/internal/config"
	"github.com/abhidanviji/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo gets the postgres db repo
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
