package dm_xml

import (
	h "dm_server/dm_helper"
	"io"

	"encoding/json"
	"encoding/xml"

	"os"
)

const smetaPath = "/files/smeta.xml"
const __name = "DM Smeta"

func identReader(encoding string, input io.Reader) (io.Reader, error) {
	return input, nil
}

// Функция для поиска нод с атрибутом ABC, равным "1.b.$" по указанному пути
func FindNodes_OldStyle(xmlPath string) []byte {
	// Открываем XML-файл для чтения
	file, err := os.Open(xmlPath)
	if err != nil {
		h.Err(err.Error())
		return []byte(h.Empty)
	}
	defer file.Close()

	/////////////////
	// Создаем объект декодера XML на основе файла
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = identReader

	var xmlEstimate XMLEstimate

	err = decoder.Decode(&xmlEstimate)
	if err != nil {
		h.Err("Ошибка декодирования XML: \n" + err.Error())
		return []byte(h.Empty)
	}

	// Преобразование в JSON
	result, err := json.Marshal(xmlEstimate)
	if err != nil {
		h.Err("Ошибка преобразования в JSON: \n" + err.Error())
		return []byte(h.Empty)
	}
	return result
}

func ParseSmeta() []byte {
	h.LogRoutineStart(__name + " ParseSmeta")
	// Путь к XML-файлу

	smeta := h.DMFilePath(smetaPath)

	h.LogRoutineEnd(__name + " ParseSmeta")
	return FindNodes_OldStyle(smeta)
}
