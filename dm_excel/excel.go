package dm_excel

import (
	h "dm_server/dm_helper"
	"fmt"
	"strconv"
	"strings"

	mysql "dm_server/dm_db/dm_mysql"

	excelize "github.com/xuri/excelize/v2"
)

func importARC(filePath string) (mysql.ARC, bool) {
	h.Log("importARC...")

	// Открытие файла
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		h.Log(err.Error())
		return mysql.ARC{}, false
	}
	newARC := mysql.ARC{}

	// Имя активного листа
	sheetName := f.GetSheetName(1)

	// название объекта
	objectName, err := f.GetCellValue(sheetName, fmt.Sprintf("B1"))
	h.Log(objectName)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		h.Log(err.Error())
		return mysql.ARC{}, false
	}

	// Начинаем чтение из 4-й строки
	startRow := 3 // Индексы строк начинаются с 0
	for i := startRow; i < len(rows); i++ {
		cellValue, err := f.GetCellValue(sheetName, fmt.Sprintf("A%d", i+1))
		if err != nil {
			h.Log(err.Error())
			return mysql.ARC{}, false
		}
		// Проверяем наличие пустой ячейки
		if cellValue == "" {
			break
		}

		isWork := false
		if strings.Contains(cellValue, ".") {
			isWork = true
		}

		if isWork {

		} else {

		}
		h.Log("Значение:", cellValue)
	}

	// Сохранение ЛЗВ в MySQL
	err = mysql.GormDB.Create(&newARC).Error
	if err != nil {
		h.Err("Ошибка при создании документа mysql.ARC: ", err.Error())
		return mysql.ARC{}, false
	}

	h.Log("ЛЗВ ", h.YellowColor, strconv.FormatInt(newARC.ID, 10), h.DefaultColor, " успешно сохранена в БД")
	return newARC, true
}

func Do() {
	importARC("./files/exampleARC.xlsx")
}
