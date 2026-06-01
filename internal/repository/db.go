package repository

import (
	"errors"

	"github.com/hoshina-dev/ticketing-service/internal/config"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func isPgUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
