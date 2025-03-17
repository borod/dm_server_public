package dm_redmine

import (
	"bytes"
	mysql "dm_server/dm_db/dm_mysql"
	h "dm_server/dm_helper"
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateUser(u mysql.User) (id int, ok bool) {
	// Формирование тела запроса в формате JSON

	requestBody := RedmineCreateUserBody{
		User: UserBody{
			Login:     "dm_" + strconv.Itoa(int(u.ID)),
			Password:  u.Password,
			Firstname: u.Name,
			Lastname:  "DM",
			Mail:      u.Email,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		h.Log("Ошибка при кодировании в JSON:", err.Error())
		return
	}
	h.Log(string(jsonData))

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		h.Log("Error marshalling request body:", err.Error())
		return 0, false
	}

	// Создание HTTP запроса к API Redmine
	req, err := http.NewRequest("POST", "https://redmine.dmcorporation.ru/users.json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		h.Log("Error creating HTTP request:", err.Error())
		return 0, false
	}

	// Добавление заголовка X-Redmine-API-Key к запросу
	req.Header.Set("X-Redmine-API-Key", "35b1f2f7fea72ac883d08c9e49605b4055599c7f")
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.Log("Error sending HTTP request:", err.Error())
		return 0, false
	}
	defer resp.Body.Close()

	// Чтение и размаршалирование JSON из тела ответа
	var redmineResponse RedmineUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&redmineResponse); err != nil {
		h.Log("Error decoding redmine response body:", err.Error())
		return 0, false
	}

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusCreated {
		h.Log("Error creating user in Redmine. Status:", resp.Status)
		return 0, false
	}

	redmineAPIKey := mysql.APIKey{
		UserID:  u.ID,
		API:     h.C_API_Redmine,
		Key:     redmineResponse.User.ApiKey,
		PartyID: strconv.Itoa(redmineResponse.User.ID),
	}
	mysql.GormDB.Save(&redmineAPIKey)

	return redmineResponse.User.ID, true
}

func ActivateUser(ID int) {

}
