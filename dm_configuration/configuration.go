package dm_configuration

import (
	h "dm_server/dm_helper"

	"encoding/json"
	"os"
	"time"
)

var Conf_MySQL DBConf
var Conf_Mongo DBConf
var Crypto CryptoConf
var Mail MailConf
var _configuration Configuration

func ReloadConfig() {
	h.LogRoutineStart(__name + " reload")

	filePath := h.DMFilePath(confFilePath)

	h.Log("Чтение файла конфигурации: \n" + filePath)
	content, err := os.ReadFile(filePath)
	if err != nil {
		h.Err("Ошибка чтения файла:\n" + err.Error())
		return
	}

	if err := json.Unmarshal(content, &_configuration); err != nil {
		if timeErr, ok := err.(*time.ParseError); ok {
			h.Err("Ошибка разбора времени:\n" + timeErr.Error())
		} else {
			h.Err("Ошибка разбора JSON:\n" + err.Error())
			return
		}
	}
	Conf_MySQL = _configuration.DB.MySQL
	Conf_Mongo = _configuration.DB.Mongo
	Crypto = _configuration.Crypto
	Mail = _configuration.Mail

	h.Log("name: " + _configuration.Name)

	h.Log("Значение поля DB:\n" + h.JsonToString(_configuration.DB))

	h.Log("Значение поля Crypto:\n" + h.JsonToString(_configuration.Crypto))

	h.LogRoutineEnd(__name + " reload")
}
