package dm_helper

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// APIs
const C_API_Redmine = 1 // "redmine"

const C_Xdmtoken = "X-dmtoken"

// const c_ProjectID = "ProjectID"
const C_ID = "ID"
const C_ObjectID = "ObjectID"
const C_Token = "token"
const C_COUNT = "COUNT"
const C_FROMDATE = "FROMDATE"
const C_TODATE = "TODATE"
const C_Page = "page"
const C_Size = "size"

// Сброс цвета
const ResetColor = "\033[0m"

// Базовые цвета
const BlackColor = "\033[38;2;0;0;0m"
const RedColor = "\033[38;2;255;0;0m"
const GreenColor = "\033[38;2;0;255;0m"
const YellowColor = "\033[38;2;255;255;0m"
const BlueColor = "\033[38;2;0;0;255m"
const PurpleColor = "\033[38;2;128;0;128m"
const CyanColor = "\033[38;2;0;255;255m"
const GrayColor = "\033[38;2;128;128;128m"
const WhiteColor = "\033[38;2;255;255;255m"

// Цвета с яркостью
const BrightRedColor = "\033[38;2;255;0;0m"
const BrightGreenColor = "\033[38;2;0;255;0m"
const BrightYellowColor = "\033[38;2;255;255;0m"
const BrightBlueColor = "\033[38;2;0;0;255m"
const BrightPurpleColor = "\033[38;2;128;0;128m"
const BrightCyanColor = "\033[38;2;0;255;255m"

var logIterator int = 1
var DefaultColor string

const C__t_no_ms = "2006-01-02 15:04:05"
const C__t_full = "2006-01-02 15:04:05.000"

const C_Str_Field_CreatedByID = "CreatedByID"
const C_Str_Field_Filter = "filter"

const Empty = "{\"}"

// concatenate strings
// func concatenateStrings(values ...*string) string {
// 	var builder strings.Builder

// 	for _, val := range values {
// 		builder.WriteString(*val)
// 	}

//		return builder.String()
//	}
func TimeCurrStr() string {
	return time.Now().Format(C__t_no_ms)
}

func TimeCurrStrMS() string {
	return time.Now().Format(C__t_full)
}

func CS(values ...string) string {
	var builder strings.Builder

	for _, val := range values {
		builder.WriteString(val)
	}

	return builder.String()
}

func JsonToString(obj any) string {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		Err("Ошибка при преобразовании:\n" + err.Error())
		return ""
	}
	return string(jsonBytes)
}

func ReverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func GetAccessRightsIDFromValues(create, read, update, delete, verify, owner bool) int {
	values := [6]bool{create, read, update, delete, verify, owner}
	binaryID := ""

	for _, val := range values {
		bit := "0"
		if val {
			bit = "1"
		}
		binaryID += bit
	}

	id, _ := strconv.ParseInt(ReverseString(binaryID), 2, 64)
	return int(id + 1)
}

func GetNamedStringFromURL(w http.ResponseWriter, r *http.Request, name string) string {
	Log("GetNamedStringFromURL...")
	Str := r.URL.Query().Get(name)
	Log(name + "=" + YellowColor + Str)
	return Str
}

func TryGetNamedIntFromURL(w http.ResponseWriter, r *http.Request, name string) (int, bool) {
	Log("TryGetNamedIntFromURL...")
	Str := GetNamedStringFromURL(w, r, name)
	Log(name + "=" + YellowColor + Str)
	ID, err := strconv.Atoi(Str)
	if err != nil {
		Log("Не удалось преобразовать " + YellowColor +
			name + DefaultColor + " в число: " + YellowColor + Str)
		return ID, false
	}
	return ID, true
}

func GetNamedIntFromURL(w http.ResponseWriter, r *http.Request, name string) (int, bool) {
	Log("getNamedIntFromURL...")
	Str := GetNamedStringFromURL(w, r, name)
	Log(name + "=" + YellowColor + Str)
	ID, err := strconv.Atoi(Str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		Err("Не удалось преобразовать " + YellowColor +
			name + DefaultColor + " в число: " + YellowColor + Str)
		return ID, false
	}
	return ID, true
}

func SetDefaultColor(c string) {
	DefaultColor = c
}

func GenerateUniqueInt(userID, itemID int) int {
	// Конкатенация userID и itemID
	str := fmt.Sprintf("%d%d", userID, itemID)

	// Вычисление хеша
	h := fnv.New32a()
	h.Write([]byte(str))
	hash := h.Sum32()

	// Преобразование хеша в int
	uniqueInt := int(hash)

	return uniqueInt
}

func LogRoutineStart(msg string) {
	logPrefix()
	fmt.Printf("%s-> %s%s %spackage initialization...\n", CyanColor, GreenColor, msg, CyanColor)
}

func LogRoutineEnd(msg string) {
	logPrefix()
	fmt.Printf("%s<- %s%s %spackage initialization finished.\n", CyanColor, GreenColor, msg, CyanColor)
}

func Err(msgs ...string) {
	logPrefixError()

	var builder strings.Builder
	builder.WriteString(RedColor)

	for _, val := range msgs {
		builder.WriteString(val)
	}

	fmt.Println(builder.String())
}

func Log(msgs ...string) {
	logPrefix()

	var builder strings.Builder
	builder.WriteString(DefaultColor)

	for _, val := range msgs {
		builder.WriteString(val)
	}

	fmt.Println(builder.String())
}

func logPrefixError() {
	fmt.Printf("%s[%s]:%d ", RedColor, time.Now().Format(C__t_no_ms), logIterator)
	logIterator++
}

func logPrefix() {
	fmt.Printf("%s[%s]:%d ", YellowColor, time.Now().Round(time.Microsecond).Format(C__t_no_ms), logIterator)
	logIterator++
}

func DMFilePath(path string) string {
	// получить путь к выполняемому файлу
	executablePath, err001 := os.Executable()
	if err001 != nil {
		Err("Ошибка: \n" + err001.Error())
		return ""
	}

	// получить каталог выполняемого файла
	executableDir := filepath.Dir(executablePath)

	// объединить путь к каталогу и имя файла
	filePath := filepath.Join(executableDir, path)

	return filePath
}

func ReplaceNaN(v interface{}) interface{} {
	switch v := v.(type) {
	case float64:
		if math.IsNaN(v) {
			return nil
		}
		return v
	case map[string]interface{}:
		for k, vv := range v {
			v[k] = ReplaceNaN(vv)
		}
	case []interface{}:
		for i, vv := range v {
			v[i] = ReplaceNaN(vv)
		}
	}
	return v
}

func InitColors() {
	DefaultColor = GrayColor

	logMessage := "Rhis is Default log color, "
	logMessage += WhiteColor + "This is White, "
	logMessage += RedColor + "This is Red, "
	logMessage += GreenColor + "This is Green, "
	logMessage += YellowColor + "This is Yellow, "
	logMessage += BlueColor + "This is Blue, "
	logMessage += PurpleColor + "This is Purple, "
	logMessage += CyanColor + "This is Cyan, "
	logMessage += GrayColor + "This is Gray."

	Log(logMessage)
}
