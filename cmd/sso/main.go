package main

import (
	"github.com/FreshOfficeFriends/SSO/internal/app"
	"github.com/FreshOfficeFriends/SSO/internal/config"

	_ "github.com/lib/pq"
)

//todo проверить количества/все месте использования os.GETENV

func main() {
	app.Run(config.New())
}
