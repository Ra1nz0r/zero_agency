package handlers

import (
	"database/sql"

	cfg "github.com/Ra1nz0r/zero_agency/internal/config"
)

type HandleQueries struct {
	*sql.DB
	//*db.Queries
	cfg.Config
}

func NewHandlerQueries(connect *sql.DB, cfg cfg.Config) *HandleQueries {
	return &HandleQueries{
		connect,
		//db.New(connect),
		cfg,
	}
}
