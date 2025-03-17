package dm_list

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	mysql "dm_server/dm_db/dm_mysql"

	auth "dm_server/dm_api/dm_authorization"

	h "dm_server/dm_helper"
)

func getLastMessage(chatID int64) (mysql.Message, bool) {
	var result mysql.Message
	mysql.GormDB.
		Where("chat_id = ?", chatID).
		Order("CT DESC").
		Limit(1).
		Preload("Author").
		Find(&result)
	return result, true
}

// TODO
// если уполномоченный (реестр уполномоченных) - полные списки
// в противном случае доступные
func GetEntitiesByUserID(entityName string, w http.ResponseWriter, r *http.Request, u *mysql.User) {
	h.Log("GetEntities...")

	currPage, pageExist := h.TryGetNamedIntFromURL(w, r, h.C_Page)
	if !pageExist {
		currPage = 1
	}

	pageSize, sizeExist := h.TryGetNamedIntFromURL(w, r, h.C_Size)
	if !sizeExist {
		pageSize = C_pageSize
	}

	orderBy := h.GetNamedStringFromURL(w, r, "orderby")

	desc := h.GetNamedStringFromURL(w, r, "desc")
	shouldOrderByDescending := len(desc) > 0

	entities := make([]interface{}, 0)

	offset := (currPage - 1) * pageSize

	orderStr := orderBy
	if shouldOrderByDescending {
		orderStr = h.CS(orderStr, " desc")
	}

	if entityName == "chat" {

		var chatEntities []mysql.Chat
		if err := mysql.GormDB.
			Limit(pageSize).
			Offset(offset).
			Order(orderStr).
			Preload("Participants").
			Find(&chatEntities).Error; err != nil {
			h.Err("Не удалось прочитать из базы " + h.YellowColor + entityName + "\n" + h.RedColor + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, chat := range chatEntities {

			lastMessage, _ := getLastMessage(chat.ID)
			authorName := lastMessage.Author.Name
			authorID := lastMessage.Author.ID
			lastMessage.Author = nil
			cri := ChatItemResponse{
				ID:          chat.ID,
				Name:        chat.Name,
				LastMessage: lastMessage,
				AuthorName:  authorName,
				AuthorID:    authorID,
			}

			// chat.Participants = nil
			entities = append(entities, cri)
		}
		// entities = append(entities, chatEntities...)
	} else if entityName == "arcworkitem" {
		if len(u.Divisions) == 0 {
			return
		}
		shouldDisplayData := false
		for _, dDB := range u.Divisions {
			if dDB.ID <= 11 {
				shouldDisplayData = true
			}
		}
		if shouldDisplayData {
			entities = GetBPEntities(w, entityName, pageSize, offset, orderStr, entities)
		}
	} else {
		entities = GetBPEntities(w, entityName, pageSize, offset, orderStr, entities)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entities)
	h.Log("Список " + h.YellowColor + entityName + h.DefaultColor + " отправлен.")
}

func GetBPEntities(w http.ResponseWriter, entityName string, pageSize int, offset int, orderStr string, entities []interface{}) []interface{} {
	e, ok1 := mysql.GetBPStruct(entityName)
	if !ok1 {
		return entities
	}
	array := reflect.New(reflect.SliceOf(e)).Interface()
	h.Log(reflect.TypeOf(array).String())

	if err := mysql.GormDB.
		Limit(pageSize).
		Offset(offset).
		Order(orderStr).
		Find(array).Error; err != nil {
		h.Err("Не удалось прочитать из базы " + h.YellowColor + entityName + "\n" + h.RedColor + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return entities
	}
	entities = append(entities, array)
	return entities
}
func ListHandler(w http.ResponseWriter, r *http.Request) {
	h.Log("ListHandler...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	entityName := strings.TrimPrefix(r.URL.Path, EndpointStringList)
	entityName = strings.ToLower(entityName)

	switch r.Method {
	case "GET":
		GetEntitiesByUserID(entityName, w, r, u)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}
