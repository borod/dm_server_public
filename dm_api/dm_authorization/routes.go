package dm_authorization

import (
	h "dm_server/dm_helper"
	"net/http"
)

var Routes = map[string]http.HandlerFunc{}

func InitRoutes() {
	h.Log("Инициализация карты маршрутов...")

	Routes["/api/register"] = CreateUserHandler
	h.Log("Добавлен маршрут /api/register -> CreateUserHandler")

	Routes["/api/resetpassword"] = ResetPassword
	h.Log("Добавлен маршрут /api/resetpassword -> ResetPassword")

	Routes["/api/login"] = LoginHandler
	h.Log("Добавлен маршрут /api/login -> LoginHandler")

	Routes["/api/user"] = ReadUserHandler
	h.Log("Добавлен маршрут /api/user -> ReadUserHandler")

	Routes["/api/getme"] = GetMe
	h.Log("Добавлен маршрут /api/getme -> GetMe")
}
