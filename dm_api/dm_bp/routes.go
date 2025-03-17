package dm_bp

import (
	h "dm_server/dm_helper"
	"net/http"
)

var Routes = map[string]http.HandlerFunc{}

const EndpointStringChat = "/api/bp/"

func InitRoutes() {
	h.Log("Инициализация карты маршрутов...")

	Routes[EndpointStringChat] = BPHandler

	h.Log("Добавлен маршрут " + EndpointStringChat + " -> ChatHandler")
}
