package dm_mailer

import (
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"

	mysql "dm_server/dm_db/dm_mysql"
	h "dm_server/dm_helper"

	conf "dm_server/dm_configuration"
)

func GetResetBody(u mysql.User) (string, bool) {
	h.Log("GetResetBody...")

	// Подгружаем HTML-шаблон из файла
	templateFile := "/opt/dm/dm_server/files/mail_templates/resetpasswordlink.html"

	// Чтение содержимого файла
	content, err := os.ReadFile(templateFile)
	if err != nil {
		h.Log("Ошибка при чтении файла:", err.Error())
		return "", false
	}

	// Преобразование содержимого файла в строку
	result := string(content)
	if err != nil {
		log.Fatalf("Ошибка чтения файла шаблона: %v", err)
	}

	// создание нового токена для сброса пароля
	token, ok := mysql.CreatePasswordResetToken(u.Email)
	if !ok {
		return "", false
	}

	mysql.GormDB.Preload("")
	result = strings.ReplaceAll(result, ANCHOR_TOKEN, token.Refresh)

	return result, true
}

func SendResetPasswordEmail(email string) bool {
	u, ok := mysql.GetUserByEmail(email)
	if !ok {
		return false
	}

	messageBody, okBody := GetResetBody(u)
	if !okBody {
		return false
	}

	SendMessage(u.Email, messageBody, "DM - изменение пароля")

	return true
}

func SendMessage(to string, messageBodyText string, subject string) {
	// // Данные пользователя и адрес SMTP-сервера
	// host := conf.Mail.Host
	// port := 587 //conf.Mail.Port
	// user := conf.Mail.User
	// password := conf.Mail.Password

	// // Адрес отправителя и получателя
	// from := conf.Mail.From

	// // Формируем сообщение
	// message := []byte("To: " + to + "\r\n" +
	// 	"From: " + from + "\r\n" +
	// 	"Subject: Your Subject Here\r\n" +
	// 	"MIME-version: 1.0\r\n" +
	// 	"Content-Type: text/html; charset=UTF-8\r\n" +
	// 	"\r\n" +
	// 	"<h1>Hello!</h1><p>This is a test email.</p>")

	// // Настройка аутентификации
	// auth := smtp.PlainAuth("", user, password, host)

	// // Подключение к серверу SMTP
	// smtpAddr := fmt.Sprintf("%s:%d", host, port)
	// client, err := smtp.Dial(smtpAddr)
	// if err != nil {
	// 	fmt.Println("Ошибка при подключении к серверу SMTP:", err)
	// 	return
	// }
	// defer client.Close()

	// // Аутентификация
	// if err := client.Auth(auth); err != nil {
	// 	fmt.Println("Ошибка аутентификации:", err)
	// 	return
	// }

	// // Отправка письма
	// if err := client.Mail(from); err != nil {
	// 	fmt.Println("Ошибка при указании отправителя:", err)
	// 	return
	// }
	// if err := client.Rcpt(to); err != nil {
	// 	fmt.Println("Ошибка при указании получателя:", err)
	// 	return
	// }
	// w, err := client.Data()
	// if err != nil {
	// 	fmt.Println("Ошибка при отправке данных:", err)
	// 	return
	// }
	// _, err = w.Write(message)
	// if err != nil {
	// 	fmt.Println("Ошибка при записи данных письма:", err)
	// 	return
	// }
	// err = w.Close()
	// if err != nil {
	// 	fmt.Println("Ошибка при закрытии записи:", err)
	// 	return
	// }

	// fmt.Println("Письмо успешно отправлено!")

	// Данные пользователя и адрес SMTP-сервера
	host := conf.Mail.Host
	port := conf.Mail.Port
	user := conf.Mail.User
	password := conf.Mail.Password

	// Адрес отправителя и получателя
	from := conf.Mail.From

	message := []byte(
		"MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			"Email body" + "\r\n" +
			messageBodyText + "\r\n")

	// Настройка SMTP-клиента
	auth := smtp.PlainAuth("", user, password, host)

	toArr := []string{
		to,
	}

	// Отправка письма по SMTP
	addr := host + ":" + strconv.Itoa(port)
	err := smtp.SendMail(addr, auth, from, toArr, message)
	if err != nil {
		h.Err("Ошибка при отправке письма: ", err.Error())
	}

	h.Log("Письмо успешно отправлено")
}
