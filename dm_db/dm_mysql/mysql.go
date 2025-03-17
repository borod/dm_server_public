package dm_mysql

import (
	conf "dm_server/dm_configuration"
	crypto "dm_server/dm_crypto"
	"reflect"
	"time"

	h "dm_server/dm_helper"

	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlDB *sql.DB
var GormDB *gorm.DB

func CreatePasswordResetToken(email string) (Token, bool) {
	h.Log("CreatePasswordResetToken...")

	u, ok := GetUserByEmail(email)
	if !ok {
		return Token{}, false
	}

	authorizationToken := crypto.Sum([]byte("Authorization" + email + time.Now().String() + crypto.Sum([]byte(u.Password))))
	refreshToken := crypto.Sum([]byte("Refresh" + email + time.Now().String() + crypto.Sum([]byte(u.Password))))

	token := Token{
		UserID:        u.ID,
		Authorization: authorizationToken,
		Refresh:       refreshToken,
	}

	GormDB.Save(&token)
	h.Log("Создан токен сброса пароля")

	return token, true
}

func GetUserByEmail(email string) (User, bool) {
	var u User
	result := GormDB.Where("email = ?", email).First(&u)
	if result.Error != nil {
		h.Log("Пользователь не найден: ", h.YellowColor, email)
		return User{}, false
	}
	return u, true
}

func GetBPStruct(entityName string) (reflect.Type, bool) {
	entityMap := map[string]reflect.Type{
		//files
		"dmfile": reflect.TypeOf(&DMFile{}),

		//chat
		"message":        reflect.TypeOf(&Message{}),
		"chat":           reflect.TypeOf(&Chat{}),
		"chatgroup":      reflect.TypeOf(&ChatGroup{}),
		"accessright":    reflect.TypeOf(&AccessRights{}),
		"userchataccess": reflect.TypeOf(&UserChatAccess{}),
		"user":           reflect.TypeOf(&User{}),
		"division":       reflect.TypeOf(&Division{}),

		//business process
		"paymentorder": reflect.TypeOf(&PaymentOrder{}),
		"payment":      reflect.TypeOf(&Payment{}),
		"invoice":      reflect.TypeOf(&Invoice{}),
		"invoiceitem":  reflect.TypeOf(&InvoiceItem{}),
		"request":      reflect.TypeOf(&Request{}),
		"requestitem":  reflect.TypeOf(&RequestItem{}),
		"analogue":     reflect.TypeOf(&Analogue{}),
		"verification": reflect.TypeOf(&Verification{}),
		"counteragent": reflect.TypeOf(&Counteragent{}),
		"arc":          reflect.TypeOf(&ARC{}),
		"arcwork":      reflect.TypeOf(&ARCWork{}),
		"arcworkitem":  reflect.TypeOf(&ARCWorkItem{}),
		"object":       reflect.TypeOf(&Object{}),
		"project":      reflect.TypeOf(&Project{}),
	}

	t, exists := entityMap[entityName]
	if !exists {
		return nil, false
	}

	return t.Elem(), true
}

func InitMySQLDB() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Conf_MySQL.User,
		conf.Conf_MySQL.Password,
		conf.Conf_MySQL.Host,
		conf.Conf_MySQL.Port,
		conf.Conf_MySQL.DBName)

	h.Log("Подключение к базе данных MySQL... \n" + dsn)

	var err error
	SqlDB, err = sql.Open("mysql", dsn)
	if err != nil {
		h.Log(h.RedColor + " Ошибка:\n" + err.Error())
	}

	GormDB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: SqlDB,
	}), &gorm.Config{})

	if err != nil {
		h.Log(h.RedColor + " Ошибка:\n" + err.Error())
	}
}

func CreateDB(db *gorm.DB) error {
	h.Log(h.GrayColor + "Создание таблиц в базе данных...")
	// Создание таблиц в базе данных
	err := db.AutoMigrate(
		&DMFile{},

		//Core structs
		&Category{},
		&Obj{},
		&Measure{},
		&Asset{},

		//Authorization + user
		&User{},
		&Token{},
		&APIKey{},

		//Divisions
		&Division{},

		//Chats
		&Message{},
		&Chat{},
		&ChatGroup{},
		&AccessRights{},
		&UserChatAccess{},

		//ARCs
		&ARC{},
		&ARCWork{},
		&ARCWorkItem{},

		//Requests
		&Request{},
		&RequestItem{},

		//Analogues
		&Analogue{},

		//Invoices
		&Invoice{},
		&InvoiceItem{},

		//Paymenrs
		&Payment{},
		&PaymentOrder{},

		//Verifications
		&Verification{},

		//Contragents
		&Counteragent{},

		//Projects Objects
		&Object{},
		&Project{},
	)

	if err != nil {
		fmt.Println(h.RedColor + " Error: " + err.Error() + h.ResetColor)
		return err
	}
	h.Log(h.GrayColor + "... Успешно завершено")

	return nil
}
