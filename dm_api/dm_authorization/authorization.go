package dm_authorization

import (
	h "dm_server/dm_helper"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	mysql "dm_server/dm_db/dm_mysql"
	mail "dm_server/dm_mailer"

	crypto "dm_server/dm_crypto"
	redmine "dm_server/dm_redmine"

	"io"
	"net/http"
	"unsafe"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	h.Log("ResetPassword...")
	if r.Method != "POST" {
		h.Err("Метод " + r.Method + " вместо ожидаемого POST")
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	ipAddress := r.RemoteAddr

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Err("Ошибка чтения body")
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	payload := string(body)

	// Десериализация
	var rpb ResetPasswordBody
	errJson := json.Unmarshal([]byte(payload), &rpb)
	if errJson != nil {
		h.Err("Ошибка десериализации:\n" + payload)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// token := h.GetNamedStringFromURL(w, r, h.C_Token)
	var existingUser mysql.User
	if rpb.Email != "" {
		result := mysql.GormDB.Where("email = ?", rpb.Email).First(&existingUser)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Пользователя с таким Email не существует
				h.Err("Пользователя с таким Email не существует:\n" + payload)
				http.Error(w, "Пользователя с таким Email не существует", http.StatusBadRequest)
				return
			} else {
				// Другая ошибка при выполнении запроса
				h.Err("Ошибка при получении пользователя:\n" + result.Error.Error())
				http.Error(w, "Ошибка при получении пользователя", http.StatusInternalServerError)
				return
			}
		}

		h.Log(rpb.Email)
		mail.SendResetPasswordEmail(rpb.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
		return
	}

	if rpb.Token == "" || rpb.Password == "" {
		h.Err("Не задан token или пароль:\n" + payload)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Хеширую пароль с использованием SHA-256
	hashedPassword := crypto.Sum(*(*[]byte)(unsafe.Pointer(&rpb.Password)))

	// Проверка наличия пользователя с таким же token
	var existingToken mysql.Token
	euRes := mysql.GormDB.Where("refresh = ?", rpb.Token).First(&existingToken)
	if euRes.Error != nil {
		http.Error(w, "Not found", http.StatusBadRequest)
		h.Err(ipAddress, " !!! Хеша не существует в базе данных")
		return
	}

	err = mysql.GormDB.First(&existingUser, existingToken.UserID).Error
	if err != nil {
		http.Error(w, "Not found", http.StatusBadRequest)
		h.Err("!!!(НЕВОЗМОЖНО) Пользователя с указанным хешем не существует")
		return
	}
	existingUser.Password = hashedPassword
	mysql.GormDB.Save(existingUser)
	h.Log("Новый хеш пароля пользотвателя сохранён в базу")

	// удаление всех существующих токенов авторизации и сброса пароля
	mysql.GormDB.Where("user_id = ?", existingUser.ID).Delete(&mysql.Token{})
	h.Log("Удалены все предыдущие токены авторизации и сброса пароля")

	// отпраавка уведомление о сбросе пароля

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		h.Err("Метод " + r.Method + " вместо ожидаемого POST")
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Err("Ошибка чтения body")
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	// Преобразование байтов в строку
	payload := string(body)

	if payload == "" {
		h.Err("Пустой body")
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	// Десериализация
	var userNamePassword UserNamePasswordBody
	errJson := json.Unmarshal([]byte(payload), &userNamePassword)
	if errJson != nil {
		h.Err("Ошибка десериализации:\n" + payload)
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	if userNamePassword.Email == "" || userNamePassword.Password == "" {
		h.Err("Не задан email или пароль:\n" + payload)
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	// Хеширую пароль с использованием SHA-256
	hashedPassword := crypto.Sum(*(*[]byte)(unsafe.Pointer(&userNamePassword.Password)))

	// Проверка наличия пользователя с таким же Email
	var existingUser mysql.User
	result := mysql.GormDB.Where("email = ?", userNamePassword.Email).First(&existingUser)
	if result.Error == nil {
		// Пользователь с таким Email уже существует
		h.Err("Пользователь с таким Email уже существует:\n" + payload)
		http.Error(w, "Пользователь с таким Email уже существует", http.StatusConflict)
		return
	}

	// Создаю нового пользователя
	u := mysql.User{
		Email:    userNamePassword.Email,
		Password: hashedPassword,
		StatusID: 1,
		Name:     userNamePassword.Name,
	}

	// Выполняю запрос на добавление пользователя в базу данных
	if err := mysql.GormDB.Create(&u).Error; err != nil {
		h.Err("Ошибка создания пользователя в базе:\n" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redmine.CreateUser(u)

	// Отправляю успешный ответ
	w.WriteHeader(http.StatusOK)
}

func GenerateOrUpdateTokens(email, password string, userID int64) (mysql.Token, error) {
	h.Log("GenerateOrUpdateTokens...")

	// Расчет значения токена
	authorizationToken := crypto.Sum([]byte("Authorization" + email + time.Now().String() + crypto.Sum([]byte(password))))
	refreshToken := crypto.Sum([]byte("Refresh" + email + time.Now().String() + crypto.Sum([]byte(password))))

	// Поиск существующего токена
	var token mysql.Token
	mysqlResult := mysql.GormDB.Where("user_id = ?", userID).First(&token)
	if mysqlResult.Error != nil {
		if errors.Is(mysqlResult.Error, gorm.ErrRecordNotFound) {
			// Токен не найден, создаем новый
			token = mysql.Token{
				UserID:        userID,
				Authorization: authorizationToken,
				Refresh:       refreshToken,
			}
			result := mysql.GormDB.Create(&token)
			if result.Error != nil {
				return token, result.Error
			}
		} else {
			return token, mysqlResult.Error
		}
	} else {
		// Токен найден, обновляем его значение
		// TODO: всегдя создаём новый
		// Токен не найден, создаем новый
		token = mysql.Token{
			UserID:        userID,
			Authorization: authorizationToken,
			Refresh:       refreshToken,
		}
		mysqlResult := mysql.GormDB.Create(&token)
		if mysqlResult.Error != nil {
			return token, mysqlResult.Error
		}
		// token.Value = tokenValue
		// result := mysql.GormDB.Save(&token)
		// if result.Error != nil {
		// 	return "", result.Error
		// }
	}

	return token, nil
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	h.Log("GetMe...")

	u := GetAuthUser(w, r)
	if nil == u {
		return
	}

	response := LoginResponse{
		UserID: u.ID,
		Name:   u.Name,
		Email:  u.Email,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func ReadUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")

	// Получение данных пользователя по email из базы данных
	var user mysql.User
	result := mysql.GormDB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// Формирование ответа
	response := struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		// тут остальные поля пользователя
	}{
		Email: user.Email,
		Name:  user.Name,
		// тут маппинг базы данных
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func GetAuthUser(w http.ResponseWriter, r *http.Request) *mysql.User {
	token := r.Header.Get(h.C_Xdmtoken)
	h.Log("Значение заголовка " + h.YellowColor + h.C_Xdmtoken + h.DefaultColor + ": " + token)
	currentUser, ok := Authorize(token)
	if !ok {
		h.Log("Ошибка авторизации по токену:\n" + token)
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return nil
	}
	h.Log(currentUser.Email + " успешно авторизован")
	return currentUser
}

func Authorize(token string) (*mysql.User, bool) {
	// Поиск токена в базе данных
	var tokenObj mysql.Token
	result := mysql.GormDB.Where("Authorization = ?", token).First(&tokenObj)
	if result.Error != nil {
		return nil, false
	}

	// Получение пользователя на основе связи с токеном
	var user mysql.User
	result = mysql.GormDB.
		Preload("Divisions").
		First(&user, tokenObj.UserID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrInvalidField) {
			// Handle the "invalid field found" error
			h.Err("!!! Ошибка с полями в структуре БД: " + result.Error.Error())
		}
		h.Err("Ошибка при получении пользователя: " + result.Error.Error())
		return nil, false
	}

	return &user, true
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	h.Log("LoginHandler...")

	if r.Method != "POST" {
		h.Err("Метод " + r.Method + " вместо ожидаемого POST")
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	h.Log("Правильный метод " + h.YellowColor + "POST")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Err("Ошибка чтения body")
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	h.Log("Есть body")

	// Преобразование байтов в структуру UserNamePassword
	var loginData UserNamePasswordBody
	err = json.Unmarshal(body, &loginData)
	if err != nil {
		h.Err("Ошибка десериализации:\n" + string(body))
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	h.Log("Успешная десериализация структура " + h.YellowColor + " UserNamePassword")

	if loginData.Email == "" || loginData.Password == "" {
		h.Err("Не задан email или пароль:\n" + string(body))
		http.Error(w, "DM: Bad request", http.StatusBadRequest)
		return
	}

	h.Log("Ненулевые поля:" + h.YellowColor + " loginData.Email, loginData.Password")

	// Поиск пользователя по Email
	var user mysql.User

	result := mysql.GormDB.
		Preload("APIKeys", "api = "+strconv.Itoa(h.C_API_Redmine)).
		Where("email = ?", loginData.Email).
		First(&user)
	if result.Error != nil {
		// Пользователь с указанным Email не найден
		h.Err("Пользователь с указанным Email не найден\n" + result.Error.Error())
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}
	h.Log("Пользователь с указанным Email найден")

	// Проверка хеша пароля
	hashedPassword := crypto.Sum([]byte(loginData.Password))
	if user.Password != hashedPassword {
		// Неудачная попытка входа - неверный пароль
		h.Err("Неудачная попытка входа - неверная хэш-сумма пароля")
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	// Генерация или обновление токена авторизации
	var token mysql.Token
	var err1 error

	token, err1 = GenerateOrUpdateTokens(user.Email, user.Password, user.ID)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}

	redmineAPIKey := "not set"
	redmineUserID := "not set"
	if len(user.APIKeys) > 0 {
		redmineAPIKey = user.APIKeys[0].Key
		redmineUserID = user.APIKeys[0].PartyID
	}

	// Формирование ответа с токеном и данными пользователя
	response := LoginResponse{
		UserID:             user.ID,
		AuthorizationToken: token.Authorization,
		RefreshToken:       token.Authorization,
		Name:               user.Name,
		Nickname:           user.Nickname,
		Email:              user.Email,
		Phone:              user.Phone,
		RedmineToken:       redmineAPIKey,
		RedmineUserID:      redmineUserID,
	}

	// Кодирование ответа в JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		h.Err("Ошибка кодирования JSON(LoginResponse)")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	h.Log("Отправка ответа завершена статус " + h.YellowColor + "http.StatusOK")
	h.Log("Пользователь " + h.YellowColor + loginData.Email + h.DefaultColor + " авторизован.")

	// Чтение данных о пользователе по email
}
