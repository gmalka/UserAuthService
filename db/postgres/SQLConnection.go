package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DbConfig struct {
	Host string
	Port string
	User string
	Dbname string
	Password string
	Sslmode string
}

func NewPostgresConnection(config DbConfig) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.Dbname, config.Sslmode))
}