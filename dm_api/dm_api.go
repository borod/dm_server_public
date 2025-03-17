package dm_api

import (
	"dm_server/dm_api/dm_authorization"

	"dm_server/dm_api/dm_entity"

	bp "dm_server/dm_api/dm_bp"

	l "dm_server/dm_api/dm_list"

	"dm_server/dm_api/dm_actions"

	h "dm_server/dm_helper"

	"net/http"
	"strings"
)

var routes = map[string]http.HandlerFunc{
	// "/api/upload":    API_upload.HandleUpload,
	// "/api/heartbeat": dm_heartbeat.HeartbitHandler,
	// "/api/estimate":  api_estimate.Handler,
	// "/api/entities":  api_entities.HandleEntities,
	// "/api/test":      api_entities.HandleTestData,
}

func InitRoutes() {
	h.Log("Создание общей карты маршрутов...")

	routes = appendRoutes(routes, l.Routes)
	h.Log("Добавлены маршруты " + h.YellowColor + " dm_list")

	routes = appendRoutes(routes, dm_authorization.Routes)
	h.Log("Добавлены маршруты " + h.YellowColor + " dm_authorization")

	routes = appendRoutes(routes, dm_entity.Routes)
	h.Log("Добавлены маршруты " + h.YellowColor + " dm_entity")

	routes = appendRoutes(routes, bp.Routes)
	h.Log("Добавлены маршруты " + h.YellowColor + " dm_bp")

	routes = appendRoutes(routes, dm_actions.Routes)
	h.Log("Добавлены маршруты " + h.YellowColor + " dm_actions")
}

func appendRoutes(baseRoutes map[string]http.HandlerFunc, newRoutes map[string]http.HandlerFunc) map[string]http.HandlerFunc {
	for key, handler := range newRoutes {
		baseRoutes[key] = handler
	}
	return baseRoutes
}

func enableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Log("http.HandlerFunc...")
		// Определение пути запроса
		path := r.URL.Path

		h.Log(r.Method + " " + h.YellowColor + path)

		h.Log("Получение всех заголовков запроса...")
		headers := r.Header

		originStr := "*"
		methodsStr := "POST, GET, PUT, PATCH, DELETE, OPTIONS"
		headersStr := "Origin, Content-Type, X-Requested-With"
		for key, values := range headers {
			for _, value := range values {
				h.Log(key + ": " + value)
				switch key {
				case "Origin":
					originStr = value
				case "Access-Control-Request-Method":
					methodsStr = value
				case "Access-Control-Request-Headers":
					headersStr = value
				}
			}
		}

		// Если пришёл OPTIONS
		if r.Method == http.MethodOptions {
			h.Log("Запрос OPTIONS")
			w.Header().Set("Access-Control-Allow-Origin", originStr)
			w.Header().Set("Access-Control-Allow-Methods", methodsStr)
			w.Header().Set("Access-Control-Allow-Headers", headersStr)

			w.WriteHeader(http.StatusNoContent) // 204 на запрос OPTIONS
			return
		}

		shouldHandleRedirect := false
		redirectHeaderValue := r.Header.Get("X-dmredirect")
		if redirectHeaderValue == h.C_Redirect_redmine {
			shouldHandleRedirect = true
		}

		if shouldHandleRedirect {
			w.Header().Set("Access-Control-Allow-Origin", originStr)
			handleRedirect(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", originStr)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Requested-With")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		bShouldReturn := false
		if strings.HasPrefix(path, l.EndpointStringList) {
			routeHandler := routes[l.EndpointStringList]
			routeHandler(w, r)
			bShouldReturn = true
		} else if strings.HasPrefix(path, dm_entity.EndpointStringEntity) {
			routeHandler := routes[dm_entity.EndpointStringEntity]
			routeHandler(w, r)
			bShouldReturn = true
		} else if strings.HasPrefix(path, bp.EndpointStringChat) {
			routeHandler := routes[bp.EndpointStringChat]
			routeHandler(w, r)
			bShouldReturn = true
		}
		if bShouldReturn {
			return
		}

		// Получение соответствующего обработчика из карты маршрутов
		routeHandler, exists := routes[path]
		if exists {
			// Вызов соответствующей функции эндпоинта
			routeHandler(w, r)
			w.Header().Set("Access-Control-Allow-Origin", originStr)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Requested-With")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			return
		}

		// Продолжение обработку
		handler.ServeHTTP(w, r)
	})
}

func StartServer() {
	// Создание обработчика маршрутов
	mux := http.NewServeMux()

	// Добавление обработчик middleware для разрешения CORS
	handler := enableCORS(mux)

	err := http.ListenAndServe(":32730", handler)

	if err != nil {
		panic(err)
	}
}
