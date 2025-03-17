package dm_heartbit

import (
	// "io/ioutil"

	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func HeartbitHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// // Читаем содержимое файла
		// fileContent, err := os.ReadFile("./heatbit")
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }

		// // Возвращаем статус OK и содержимое файла в теле ответа
		// w.WriteHeader(http.StatusOK)
		// w.Write(fileContent)

		// Параметры подключения к MySQL
		dbUser := "dm_easy"
		dbPass := "765567765567"
		dbHost := "dmcorporation.ru"
		dbName := "dm_easy" // замените на имя вашей базы данных

		// Формирование строки подключения
		dbURI := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPass, dbHost, dbName)

		// Подключение к MySQL
		db, err := sql.Open("mysql", dbURI)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Извлечение значения параметра "username"
		username := r.FormValue("username")

		// Выполнение запроса к базе данных
		query := "SELECT 1 FROM user WHERE email = ?"
		var result int
		err = db.QueryRow(query, username).Scan(&result)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		// // Получаем текущий timestamp
		// timestamp := time.Now().Format("2006-01-02 15:04:05")

		// // Записываем timestamp в файл
		// err := os.WriteFile("./heatbit", []byte(timestamp), 0644)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }

		// // Возвращаем статус OK и сообщение об успешной записи
		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("Таймстамп успешно записан"))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
