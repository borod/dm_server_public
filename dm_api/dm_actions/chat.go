package dm_actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	h "dm_server/dm_helper"

	mysql "dm_server/dm_db/dm_mysql"

	auth "dm_server/dm_api/dm_authorization"

	bp "dm_server/dm_api/dm_bp"

	"gorm.io/gorm"
)

func HandleExcludeChatParticipants(w http.ResponseWriter, r *http.Request) {
	h.Log("HandleExcludeChatParticipants...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return
	}

	var chatDB mysql.Chat
	err := mysql.GormDB.
		Preload("Participants").
		First(&chatDB, ID).Error
	if err != nil {
		h.Err("Не удалось найти mysql.Chat " + h.YellowColor + strconv.Itoa(ID) + "\n" + h.RedColor + err.Error())
		http.Error(w, "No such chat", http.StatusBadRequest)
		return
	}

	// Десериализация тела запроса в структуру RemoveChatParticipantsPayload
	var payload CreateChatPayload
	errPayload := json.NewDecoder(r.Body).Decode(&payload)
	if errPayload != nil {
		h.Err("Ошибка десериализации тела запроса:" + errPayload.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var responses []IncExclChatParticipantsResponse

	for _, uID := range payload.Participants {
		var userDB mysql.User
		err := mysql.GormDB.First(&userDB, uID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Пользователь с ID uID не существует
				responses = append(responses, IncExclChatParticipantsResponse{
					ID:     uID,
					Status: h.C_ChatUserDoesNotExist,
				})
				h.Log(fmt.Sprintf("Пользователь с ID %d не существует", uID))
			} else {
				responses = append(responses, IncExclChatParticipantsResponse{
					ID:     uID,
					Status: h.C_ChatUserUndefined,
				})
				h.Log(fmt.Sprintf("Ошибка при поиске пользователя с ID %d: %s", uID, err.Error()))
			}
			continue
		}

		// Проверяем, включен ли пользователь в chat
		var newParticipants []mysql.User
		userExcluded := false
		for i, uDB := range chatDB.Participants {
			if uDB.ID == userDB.ID {
				// Удаляем пользователя из slice
				newParticipants = append(chatDB.Participants[:i], chatDB.Participants[i+1:]...)
				responses = append(responses, IncExclChatParticipantsResponse{
					ID:     userDB.ID,
					Status: h.C_ChatUserSuccess,
				})
				userExcluded = true
				// newParticipants = append(newParticipants, uDB)
				break
			}
		}

		if !userExcluded {
			responses = append(responses, IncExclChatParticipantsResponse{
				ID:     userDB.ID,
				Status: h.C_ChatUserAlready,
			})
		}

		// Обновляем участников чата
		chatDB.Participants = newParticipants
		err = mysql.GormDB.Model(&chatDB).Association("Participants").Replace(newParticipants)
		if err != nil {
			h.Err(fmt.Sprintf("Ошибка при обновлении участников чата: %s", err.Error()))
			return
		}

		err = mysql.GormDB.Save(&chatDB).Error
		if err != nil {
			h.Err(fmt.Sprintf("Ошибка при сохранении chat: %s", err.Error()))
			return
		}
	}

	// Отправляем JSON-ответ с массивом RemoveChatParticipantsResponse
	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		h.Err(fmt.Sprintf("Ошибка при сериализации JSON-ответа: %s", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		h.Err(fmt.Sprintf("Ошибка при записи JSON-ответа: %s", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func HandleIncludeChatParticipants(w http.ResponseWriter, r *http.Request) {
	h.Log("HandleIncludeChatParticipants...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return
	}

	var chatDB mysql.Chat
	err := mysql.GormDB.
		Preload("Participants").
		First(&chatDB, ID).Error
	if err != nil {
		h.Err("Не удалось найти mysql.Chat " + h.YellowColor + strconv.Itoa(ID) + "\n" + h.RedColor + err.Error())
		http.Error(w, "No such chat", http.StatusBadRequest)
		return
	}

	// Десериализация тела запроса в структуру CreateChatPayload
	var payload CreateChatPayload
	errPayload := json.NewDecoder(r.Body).Decode(&payload)
	if errPayload != nil {
		h.Err("Ошибка десериализации тела запроса:" + errPayload.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var responses []IncExclChatParticipantsResponse

	for _, uID := range payload.Participants {
		var userDB mysql.User
		err := mysql.GormDB.First(&userDB, uID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Пользователь с ID uID не существует
				responses = append(responses, IncExclChatParticipantsResponse{
					ID:     uID,
					Status: h.C_ChatUserDoesNotExist,
				})
				h.Log(fmt.Sprintf("Пользователь с ID %d не существует", uID))
			} else {
				responses = append(responses, IncExclChatParticipantsResponse{
					ID:     uID,
					Status: h.C_ChatUserUndefined,
				})
				h.Log(fmt.Sprintf("Ошибка при поиске пользователя с ID %d: %s", uID, err.Error()))
			}
			continue
		}

		// Проверяем, включен ли уже пользователь в chat
		alreadyIncluded := false
		for _, uDB := range chatDB.Participants {
			if uDB.ID == userDB.ID {
				responses = append(responses, IncExclChatParticipantsResponse{
					ID:     userDB.ID,
					Status: h.C_ChatUserAlready,
				})
				alreadyIncluded = true
				break
			}
		}

		if !alreadyIncluded {
			chatDB.Participants = append(chatDB.Participants, userDB)
			responses = append(responses, IncExclChatParticipantsResponse{
				ID:     userDB.ID,
				Status: h.C_ChatUserSuccess,
			})
		}
	}

	err = mysql.GormDB.Save(&chatDB).Error
	if err != nil {
		h.Err(fmt.Sprintf("Ошибка при сохранении chat: %s", err.Error()))
		return
	}

	// Отправляем JSON-ответ с массивом IncludeChatParticipantsResponse
	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		h.Err(fmt.Sprintf("Ошибка при сериализации JSON-ответа: %s", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		h.Err(fmt.Sprintf("Ошибка при записи JSON-ответа: %s", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func HandleCreateChat(w http.ResponseWriter, r *http.Request) {
	h.Log("HandleCreateChat...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	// Десериализация тела запроса в структуру CreateChatPayload
	var payload CreateChatPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.Err("Ошибка десериализации тела запроса:" + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Использование десериализованных данных
	h.Log("Имя чата:" + h.YellowColor + payload.Name)

	chatID := bp.CreateChat(payload.Name, u, payload.Participants, payload.Justification)

	// Отправка успешного ответа

	response := struct {
		ID int64 `json:"ID"`
	}{
		ID: chatID,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Err("Ошибка при преобразовании в JSON: " + err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func GetChatParticipants(w http.ResponseWriter, r *http.Request) {
	h.Log("GetChatParticipants...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	chatID, OK := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !OK {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var allowedChats []mysql.UserChatAccess
	err := mysql.GormDB.Where("chat_id = ?", chatID).Preload("User").Find(&allowedChats).Error
	if err != nil {
		//!!!TODO
		h.Err("Ошибка при извлечении участников чата из базы данных: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response []ChatParticipantsResponse

	for _, ac := range allowedChats {
		r := ChatParticipantsResponse{
			ID:       ac.User.ID,
			Name:     ac.User.Name,
			Nickname: ac.User.Nickname,
			Email:    ac.User.Email,
		}
		response = append(response, r)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func HandleUserChats(w http.ResponseWriter, r *http.Request) {
	h.Log("HandleUserChats...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	var allowedChats []mysql.UserChatAccess

	// Assuming you have a user object `u` with an `ID` field
	err := mysql.GormDB.Where("user_id = ?", u.ID).Preload("Chat").Find(&allowedChats).Error

	if err != nil {
		// Handle the error
		//!!!TODO
		h.Err("Error:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response []UserChatResponse

	for _, ac := range allowedChats {
		// h.Log("Chat ID:", strconv.Itoa(int(chat.ID)))
		// h.Log("Chat Name:", chat.Chat.Name)
		mysql.GormDB.Preload("Messages").First(&ac.Chat)
		c := UserChatResponse{
			ID:            ac.Chat.ID,
			Name:          ac.Chat.Name,
			Description:   ac.Chat.Description,
			CountMessages: len(ac.Chat.Messages),
		}
		response = append(response, c)
	}

	mysql.GormDB.Preload("Chats").First(&u)

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Err("Ошибка при преобразовании в JSON: " + err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func handle_ChatDetails(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_ChatDetails...")

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return nil
	}

	var chatDB mysql.Chat
	err := mysql.GormDB.
		Preload("Participants").
		Preload("Owner").
		Preload("Messages").
		Preload("Messages.Author").
		First(&chatDB, ID).Error
	if err != nil {
		h.Err("Не удалось найти mysql.RequestItem " + h.YellowColor + strconv.Itoa(ID) + "\n" + h.RedColor + err.Error())
		return nil
	}

	var messages []MessageResponse
	for _, cmDB := range chatDB.Messages {
		m := MessageResponse{
			CT:       cmDB.CT,
			AuthorID: cmDB.AuthorID,
			Author:   cmDB.Author.Name,
			Content:  cmDB.Content,
			ID:       cmDB.ID,
			ChatID:   cmDB.ChatID,
		}
		messages = append(messages, m)
	}

	result := PageChatDetilsResponse{
		Name:              chatDB.Name,
		CountMessages:     len(chatDB.Messages),
		CountParticipants: len(chatDB.Participants),
		Messages:          messages,
	}

	return result
}
