package database

import (
	"database/sql"
	"fmt"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Name     string
	SSLMode  string
	Password string
}

func NewPostgresConnection(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.SSLMode, cfg.Password))
	if err != nil {
		return nil, fmt.Errorf("database open err | %s", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping err | %s", err)
	}

	return db, nil
}
