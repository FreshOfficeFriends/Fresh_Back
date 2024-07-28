package main

import (
	"github.com/FreshOfficeFriends/SSO/internal/app"
	"github.com/FreshOfficeFriends/SSO/internal/config"

	_ "github.com/lib/pq"
)

//todo реализовать хранение refresh-tokens в redis

//todo запретить запросы к ендпоинтам микросервисов не через nginx
//todo убрать все os.GETenv
//todo разбить интерфейсы
//todo нейминг пофикси.............
//todo отрефачить email connection
//todo шатдаун (closer)

//----------------------------DONE
//добавить endpoint для восстановления паролей
//сконфигурировать адекватный api-gateway (nginx)
//вынести валидацию в отдельный ендпоинт

func main() {
	app.Run(config.New())
}
