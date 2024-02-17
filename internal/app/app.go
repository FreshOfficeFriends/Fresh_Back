package app

import (
	"fmt"

	"github.com/FreshOfficeFriends/SSO/internal/config"
	"github.com/FreshOfficeFriends/SSO/pkg/database"
	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

func Run(cfg *config.Config) {
	db, err := database.NewPostgresConnection(&cfg.DB)
	if err != nil {
		logger.Debug(fmt.Sprintf("%+v", cfg.DB))
		logger.Fatal(err.Error())
	}

	defer db.Close()

	logger.Info("DB connected")
}
