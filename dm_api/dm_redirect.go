package dm_api

import (
	"io"
	"net/http"
)

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Создание нового запроса к redmine.dmcorporation.ru с теми же параметрами, что и изначальный запрос.
	// Включение тела запроса, если оно было.
	req, err := http.NewRequest(r.Method, "https://redmine.dmcorporation.ru"+r.RequestURI, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Копирвание заголовки из оригинального запроса в новый запрос.
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Отправляем новый запрос на redmine.dmcorporation.ru.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Копируем заголовки и статус ответа от redmine.dmcorporation.ru в ответ клиенту.
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Передаем тело ответа от redmine.dmcorporation.ru обратно клиенту.
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
