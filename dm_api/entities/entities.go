package api_entities

import (
	h "dm_server/dm_helper"
	"encoding/json"
	"net/http"
)

const __name = "DM XML"

func GetEntities() []byte {
	// Создание объектов для первого проекта
	// objects1 := []Object{
	// 	{
	// 		ID:             1,
	// 		ProjectID:      1,
	// 		ProjectName:    "Подпорная стена",
	// 		Address:        "Москва",
	// 		CompletionDate: "01.01.2024",
	// 		Name:           "Люберцы",
	// 	},
	// 	{
	// 		ID:             2,
	// 		ProjectID:      1,
	// 		ProjectName:    "Подпорная стена",
	// 		Address:        "Москва",
	// 		CompletionDate: "01.01.2024",
	// 		Name:           "первый объект всегда будет по имени проекта?",
	// 	},
	// 	{
	// 		ID:             3,
	// 		ProjectID:      1,
	// 		ProjectName:    "Подпорная стена",
	// 		Address:        "Москва",
	// 		CompletionDate: "01.01.2024",
	// 		Name:           "Что-то",
	// 	},
	// }

	// // Создание объектов для второго проекта
	// objects2 := []Object{
	// 	{
	// 		ID:             4,
	// 		ProjectID:      2,
	// 		ProjectName:    "Какой-то другой проект",
	// 		CompletionDate: "В следующем году",
	// 		Name:           "Object 2.1",
	// 	},
	// 	{
	// 		ID:             5,
	// 		ProjectID:      2,
	// 		ProjectName:    "Какой-то другой проект",
	// 		CompletionDate: "В следующем году",
	// 		Name:           "Object 2.2",
	// 	},
	// 	{
	// 		ID:             6,
	// 		ProjectID:      2,
	// 		ProjectName:    "Какой-то другой проект",
	// 		CompletionDate: "В следующем году",
	// 		Name:           "Object 2.3",
	// 	},
	// }

	// works := []Work{
	// 	{
	// 		ID:   1,
	// 		Name: "Установка окон",
	// 		Resources: []Resource{
	// 			{
	// 				ID:         1,
	// 				Type:       "материал",
	// 				Name:       "Стеклопакеты",
	// 				Measure:    "шт",
	// 				Qty:        10,
	// 				Expediture: 5000,
	// 			},
	// 			{
	// 				ID:         2,
	// 				Type:       "материал",
	// 				Name:       "Крепежные элементы",
	// 				Measure:    "шт",
	// 				Qty:        50,
	// 				Expediture: 1000,
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:   2,
	// 		Name: "Отделка стен",
	// 		Resources: []Resource{
	// 			{
	// 				ID:         3,
	// 				Type:       "материал",
	// 				Name:       "Гипсокартонные панели",
	// 				Measure:    "лист",
	// 				Qty:        100,
	// 				Expediture: 20000,
	// 			},
	// 			{
	// 				ID:         4,
	// 				Type:       "материал",
	// 				Name:       "Шпатлевка",
	// 				Measure:    "кг",
	// 				Qty:        5,
	// 				Expediture: 1000,
	// 			},
	// 		},
	// 	},
	// 	{
	// 		ID:   3,
	// 		Name: "Укладка плитки",
	// 		Resources: []Resource{
	// 			{
	// 				ID:         5,
	// 				Type:       "материал",
	// 				Name:       "Керамическая плитка",
	// 				Measure:    "м²",
	// 				Qty:        20,
	// 				Expediture: 8000,
	// 			},
	// 			{
	// 				ID:         6,
	// 				Type:       "материал",
	// 				Name:       "Клей для плитки",
	// 				Measure:    "кг",
	// 				Qty:        3,
	// 				Expediture: 1500,
	// 			},
	// 		},
	// 	},
	// }

	// arc := AvailableResourcesCatalog{
	// 	ID:    1,
	// 	Name:  "Лимитно-забоная ведомость",
	// 	Works: works,
	// }

	// Создание двух проектов с объектами
	// projects := []Project{
	// 	{
	// 		ID:             1,
	// 		Name:           "Подпорная стена",
	// 		Objects:        objects1,
	// 		Address:        "Москва",
	// 		CompletionDate: "01.01.2024",
	// 		ARC:            arc,
	// 		Price:          "44587687.0932",
	// 	},
	// 	{
	// 		ID:             2,
	// 		Name:           "Какой-то другой проект",
	// 		Objects:        objects2,
	// 		Address:        "какой-то очередной адрес",
	// 		CompletionDate: "В следующем году",
	// 		ARC:            arc,
	// 		Price:          "123.000.000,567",
	// 	},
	// }

	// requests := GetDummyRequests()
	// entities := Entities{
	// 	Projects: projects,
	// 	Objects:  append(objects1, objects2...),
	// 	ARC:      arc,
	// 	Requests: requests,
	// }

	type Entities struct {
		Projects []interface{} `json:"Projects"`
		Objects  []interface{} `json:"Objects"`
		ARC      interface{}   `json:"ARC"`
		Requests []interface{} `json:"Requests"`
	}

	entities := Entities{
		Projects: nil,
		Objects:  nil,
		ARC:      nil,
		Requests: nil,
	}

	// Преобразование в JSON
	result, err := json.Marshal(entities)

	if err != nil {
		h.Log("Ошибка преобразования в JSON: \n" + err.Error())
		return nil
	}

	return result
}

func HandleTestData(w http.ResponseWriter, r *http.Request) {
	// w.Write(GetEntities())
	w.Write([]byte(""))
	w.WriteHeader(http.StatusOK)
}

func HandleEntities(w http.ResponseWriter, r *http.Request) {
	// w.Write(GetEntities())
	w.Write([]byte(""))
	w.WriteHeader(http.StatusOK)
}
