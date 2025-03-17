package dm_entity

import (
	_ "dm_server/dm_helper"
	h "dm_server/dm_helper"
	"net/http"
)

var Routes = map[string]http.HandlerFunc{}

const EndpointStringEntity = "/api/entity/"

func InitRoutes() {
	h.Log("Инициализация карты маршрутов...")

	Routes[EndpointStringEntity] = EntityHandler

	h.Log("Добавлен маршрут " + EndpointStringEntity + " -> EntityHandler")
}
