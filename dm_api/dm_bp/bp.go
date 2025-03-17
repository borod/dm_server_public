package dm_bp

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	mysql "dm_server/dm_db/dm_mysql"

	auth "dm_server/dm_api/dm_authorization"

	h "dm_server/dm_helper"
)

func CreateChat(chatName string, u *mysql.User, mandatoryValidators []int64, msg string) int64 {

	// Создание чата
	chat := mysql.Chat{
		Name:        chatName + " " + time.Now().Format("2006.01.02 15:04:05"),
		Description: "Обсуждение " + chatName,
		OwnerID:     u.ID,
	}
	db := mysql.GormDB

	// Добавление обязательных участников в Participants
	participants := []mysql.User{}
	for _, userID := range mandatoryValidators {
		user := mysql.User{}
		if err := db.First(&user, userID).Error; err != nil {
			// Обработка ошибки, если не удалось найти пользователя
			h.Err("Failed to find user: " + err.Error())
			//   http.Error(w, fmt.Sprintf("Failed to find user: %s", err), http.StatusInternalServerError)
			return -1
		}

		participants = append(participants, user)
	}
	//добавление создателя чата в участники
	participants = append(participants, *u)
	chat.Participants = participants

	if err := db.Create(&chat).Error; err != nil {
		h.Err("Failed to create chat: " + err.Error())
		// http.Error(w, fmt.Sprintf("Failed to create chat: %s", err), http.StatusInternalServerError)
		return -1
	}

	// // Создание записи UserChatAccess для обязательных участников
	// accessRights := mysql.AccessRights{
	// 	Read:   true,
	// 	Update: true,
	// 	Verify: true,
	// }
	for _, participant := range participants {
		userChatAccess := mysql.UserChatAccess{
			UserID:         participant.ID,
			ChatID:         chat.ID,
			AccessRightsID: h.C_RightsFullAccess, // Замените на соответствующий ID записи AccessRights в вашей базе данных
		}
		if err := db.Create(&userChatAccess).Error; err != nil {
			h.Err("Failed to create user chat access 1: " + err.Error())
			//   http.Error(w, fmt.Sprintf("Failed to create user chat access: %s", err), http.StatusInternalServerError)
			return -1
		}

		participant.Chats = append(participant.Chats, chat)

	}

	// Создание первого сообщения в чате
	var content string
	if msg == "" {
		content = "Ожидание согласования " + time.Now().Format("2006.01.02 15:04:05")
	} else {
		content = msg
	}
	message := mysql.Message{
		Content:  content,
		AuthorID: u.ID,
		ChatID:   chat.ID,
	}
	if err := db.Create(&message).Error; err != nil {
		h.Err("Failed to create message: " + err.Error())
		return -1
	}

	return chat.ID
}

func BPHandler(w http.ResponseWriter, r *http.Request) {
	h.Log("BPHandler...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	entityName := strings.TrimPrefix(r.URL.Path, EndpointStringChat)
	entityName = strings.ToLower(entityName)

	switch r.Method {
	case http.MethodGet:
		GetEntity(entityName, w, r)
	case http.MethodPost:
		CreateEntity(entityName, w, r, u)
	case http.MethodPatch:
	case http.MethodPut:
		UpdateEntity(entityName, w, r)
	case http.MethodDelete:
		DeleteEntity(entityName, w, r)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func GetEntity(entityName string, w http.ResponseWriter, r *http.Request) {
	h.Log("GetEntity...")
	e, ok1 := mysql.GetBPStruct(entityName)
	if !ok1 {
		return
	}
	entity := reflect.New(e).Interface()
	h.Log(reflect.TypeOf(entity).String())
	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return
	}
	IDStr := strconv.Itoa(ID)
	h.Log("ID=", h.YellowColor, IDStr)

	if err := mysql.GormDB.First(entity, ID).Error; err != nil {
		h.Err("Не удалось прочитать из базы " + h.YellowColor + entityName + " " + IDStr + "\n" + h.RedColor + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}

func CreateEntity(entityName string, w http.ResponseWriter, r *http.Request, u *mysql.User) {
	h.Log("CreateEntity...")
	e, ok1 := mysql.GetBPStruct(entityName)
	if !ok1 {
		return
	}
	entity := reflect.New(e).Interface()

	h.Log(reflect.TypeOf(entity).String())

	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		h.Err("CreateChatEntity - Ошибка десериализации объекта из тела запроса\n" + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := reflect.ValueOf(entity).Elem()

	// for i := 0; i < v.NumField(); i++ {
	// 	field := v.Field(i)
	// 	fieldType := v.Type().Field(i)
	// 	h.Log(fmt.Sprintf("%s (%s) = %v\n", fieldType.Name, fieldType.Type, field))
	// }

	CreatedByID := v.FieldByName(h.C_Str_Field_CreatedByID)
	if CreatedByID.IsValid() {
		CreatedByID.Set(reflect.ValueOf(u.ID))
	}

	if err := mysql.GormDB.Create(entity).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}

func UpdateEntity(entityName string, w http.ResponseWriter, r *http.Request) {
	h.Log("UpdateEntity...")
	e, ok1 := mysql.GetBPStruct(entityName)
	if !ok1 {
		return
	}
	entity := reflect.New(e).Interface()
	h.Log(reflect.TypeOf(entity).String())

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return
	}
	IDStr := strconv.Itoa(ID)
	h.Log("ID=", h.YellowColor, IDStr)

	if err := mysql.GormDB.First(entity, ID).Error; err != nil {
		h.Err("Не удалось прочитать из базы " + h.YellowColor + entityName + " " + IDStr + "\n" + h.RedColor + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := mysql.GormDB.Save(entity).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity)
}

func DeleteEntity(entityName string, w http.ResponseWriter, r *http.Request) {
	h.Log("DeleteEntity...")
	e, ok1 := mysql.GetBPStruct(entityName)
	if !ok1 {
		return
	}
	entity := reflect.New(e).Interface()
	h.Log(reflect.TypeOf(entity).String())

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return
	}
	IDStr := strconv.Itoa(ID)
	h.Log("ID=", h.YellowColor, IDStr)

	if err := mysql.GormDB.First(&entity, ID).Error; err != nil {
		h.Err("Не удалось прочитать из базы " + h.YellowColor + entityName + " " + IDStr + "\n" + h.RedColor + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := mysql.GormDB.Delete(entity).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
