package app

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/FreshOfficeFriends/SSO/internal/config"
	"github.com/FreshOfficeFriends/SSO/internal/service/auth"
	"github.com/FreshOfficeFriends/SSO/internal/storage/psql"
	rds "github.com/FreshOfficeFriends/SSO/internal/storage/redis"
	"github.com/FreshOfficeFriends/SSO/internal/transport/rest"
	"github.com/FreshOfficeFriends/SSO/pkg/database"
	"github.com/FreshOfficeFriends/SSO/pkg/hash"
	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

const (
	pathMigrations = "file://migrations"
)

func Run(cfg *config.Config) {
	db, err := database.NewPostgresConnection(&cfg.DB)
	if err != nil {
		logger.Debug(fmt.Sprintf("%+v", cfg.DB))
		logger.Fatal(err.Error())
	}

	defer db.Close()

	logger.Info("pg connected")

	migrateDB(db)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisDB := rds.New(rdb)

	logger.Info("redis connected")

	hasher := hash.NewSHA1Hasher(os.Getenv("hash_salt"))

	usersRepo := psql.NewUsers(db)
	usersService := auth.NewAuth(usersRepo, hasher, redisDB)

	handler := rest.NewHandler(usersService)

	srv := &http.Server{
		Addr:              "localhost:8080",
		Handler:           handler.InitRouter(),
		ReadHeaderTimeout: 0,
	}

	logger.Info("localhost:8080")

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("")
	}
}

func migrateDB(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal("Couldn't get database instance for running migrations.", zap.Error(err))
	}

	m, err := migrate.NewWithDatabaseInstance(pathMigrations, os.Getenv("DB_NAME"), driver)
	if err != nil {
		logger.Fatal("Couldn't create migrate instance.", zap.Error(err))
	}

	if err := m.Up(); err != nil && errors.Is(err, errors.New("no change")) {
		logger.Fatal("Couldn't run database migration.", zap.Error(err))
	} else {
		logger.Info("Database migration was run successfully")
	}
}
