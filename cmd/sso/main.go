package main

import (
	"github.com/FreshOfficeFriends/SSO/internal/app"
	"github.com/FreshOfficeFriends/SSO/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	app.Run(config.New())
}
