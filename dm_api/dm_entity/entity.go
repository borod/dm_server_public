package dm_entity

import (
	_ "dm_server/dm_helper"
	h "dm_server/dm_helper"

	_ "dm_server/dm_db/dm_mongo"
	mongo "dm_server/dm_db/dm_mongo"

	_ "dm_server/dm_db/dm_mysql"
	mysql "dm_server/dm_db/dm_mysql"

	_ "dm_server/dm_api/dm_authorization"
	auth "dm_server/dm_api/dm_authorization"

	"encoding/json"
	"fmt"
	"io"
	"strings"

	"net/http"
	"net/url"
)

// create
func Post(entityName string, u mysql.User, w http.ResponseWriter, r *http.Request) {
	h.Log("HTTP PUT/PATCH Handler...")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %s", err), http.StatusInternalServerError)
		return
	}

	// Преобразование байтов в строку
	payload := string(body)

	err = mongo.CreateEntities(entityName, payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create documents: %s", err), http.StatusInternalServerError)
		return
	}
}

// update
func PutOrPatch(entityName string, w http.ResponseWriter, r *http.Request) {
	h.Log("HTTP PUT/PATCH Handler...")

	queryString, err := url.QueryUnescape(r.URL.Query().Get("mdbQuery"))
	// queryString := r.URL.RawQuery

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %s", err), http.StatusInternalServerError)
		return
	}

	// Преобразование байтов в строку
	payload := string(body)

	err = mongo.UpdateWithJSONFilterAndPayload(entityName, queryString, payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update documents: %s", err), http.StatusInternalServerError)
		return
	}
}

// read
func Get(entityName string, w http.ResponseWriter, r *http.Request) {
	h.Log("HTTP GET Handler...")

	// queryString := r.URL.RawQuery
	queryString, err := url.QueryUnescape(r.URL.Query().Get("mdbQuery"))

	h.Log("Строка запроса: " + h.YellowColor + " " + queryString)

	objects, err := mongo.GetCollectionObjects(entityName, queryString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get collection objects: %s", err), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(objects)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal collection objects: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// delete
func Delete(entityName string, w http.ResponseWriter, r *http.Request) {
	h.Log("HTTP DELETE Handler...")

	queryString, err := url.QueryUnescape(r.URL.Query().Get("mdbQuery"))

	h.Log("Строка запроса: " + h.YellowColor + " " + queryString)

	errMongo := mongo.DeleteMongoEntities(entityName, queryString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete documents: %s", errMongo), http.StatusInternalServerError)
		return
	}
}

func EntityHandler(w http.ResponseWriter, r *http.Request) {
	h.Log("EntityHandler...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	entityName := strings.TrimPrefix(r.URL.Path, EndpointStringEntity)

	switch r.Method {
	case "GET":
		Get(entityName, w, r)
	case "POST":
		Post(entityName, *u, w, r)
	case "PUT", "PATCH":
		PutOrPatch(entityName, w, r)
	case "DELETE":
		Delete(entityName, w, r)

	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}
