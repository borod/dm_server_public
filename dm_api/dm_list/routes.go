package dm_list

import (
	h "dm_server/dm_helper"
	"net/http"
)

var Routes = map[string]http.HandlerFunc{}

const EndpointStringList = "/api/list/"

func InitRoutes() {
	h.Log("Инициализация карты маршрутов...")

	Routes[EndpointStringList] = ListHandler

	h.Log("Добавлен маршрут " + EndpointStringList + " -> ListHandler")
}
