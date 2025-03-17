package dm_actions

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	auth "dm_server/dm_api/dm_authorization"

	h "dm_server/dm_helper"

	mysql "dm_server/dm_db/dm_mysql"

	bp "dm_server/dm_api/dm_bp"
)

var c_invoiceValidators = []int64{1, 2, 3, 4, 5, 6, 8, 9, 10}
var c_requestVerifiers = []int64{5, 6, 8, 9, 10}

// func addVerifier(uID int) (AddedVerifierResponse, bool) {

// 	result := AddedVerifierResponse{}

//		return result, true
//	}

func HandleAddVerifier(w http.ResponseWriter, r *http.Request) {

	var response []UserChatResponse

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Err("Ошибка при преобразовании в JSON: " + err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func HandleInvoicetToVerify(w http.ResponseWriter, r *http.Request) {
	h.Log("HandleInvoicetToVerify...")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var iDB mysql.Invoice
	err := mysql.GormDB.
		Preload("InvoiceItems", "invoice_id = ?", ID).
		First(&iDB, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.Err("Счёт не найден")
			http.Error(w, "Not found", http.StatusBadRequest)
		} else {
			h.Err("Ошибка при получении объекта mysql.Invoice из MySQL: " + err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if iDB.Status > 1 {
		h.Err("Попытка перевести счёт в статус `Создан` из статуса: ", h.YellowColor, strconv.Itoa(iDB.Status))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// add verifications
	ok = initVerifications(iDB.ID, 0, c_invoiceValidators)
	if !ok {
		h.Err("Ошибка при инициализации согласований")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	iIDStr := strconv.FormatInt(iDB.ID, 10)
	h.Log("Созданы согласования для счёта ", h.YellowColor, iIDStr)

	// Create chat
	chatName := h.CS("Чат по счёту ", strconv.FormatInt(iDB.ID, 10))
	chatID := bp.CreateChat(chatName, u, c_invoiceValidators, h.CS("Согласование счёта от ", h.TimeCurrStr()))

	h.Log("Создан чат: ", h.YellowColor+chatName+h.DefaultColor+", ChatID: ", h.YellowColor, strconv.FormatInt(chatID, 10))

	iDB.Status = h.C_InvoiceCreated
	iDB.ChatID = &chatID
	err = mysql.GormDB.Save(&iDB).Error
	if err != nil {
		h.Err("Ошибка при сохранении iDB:\n" + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(responseJSON)
}

func HandleRequestToVerify(w http.ResponseWriter, r *http.Request) {
	h.Log("HandleRequestToWork...")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var rDB mysql.Request
	err := mysql.GormDB.
		Preload("RequestItems", "request_id = ?", ID).
		First(&rDB, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.Err("Заявка не найдена")
			http.Error(w, "Requestnot found", http.StatusBadRequest)
		} else {
			h.Err("Ошибка при получении объекта mysql.Request из MySQL: " + err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if rDB.Status > 1 {
		h.Err("Попытка перевести счёт в статус `Создан` из статуса: ", h.YellowColor, strconv.Itoa(rDB.Status))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, riDB := range rDB.RequestItems {
		// add verifications
		ok := initVerifications(0, riDB.ID, c_requestVerifiers)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		riIDStr := strconv.FormatInt(riDB.ID, 10)
		h.Log("Созданы согласования для заказа ", h.YellowColor, riIDStr)

		riDB.Status = h.C_RequestCreated
		err = mysql.GormDB.Save(&riDB).Error
		if err != nil {
			h.Err("Ошибка при сохранении riDB:\n", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Create chat
	chatName := h.CS("Чат по заявке ", strconv.FormatInt(rDB.ID, 10))
	chatID := bp.CreateChat(chatName, u, c_requestVerifiers, h.CS("Согласование заявки от ", h.TimeCurrStr()))
	h.Log("Создан чат: ", h.YellowColor+chatName+h.DefaultColor+", ChatID: ", h.YellowColor, strconv.FormatInt(chatID, 10))

	rDB.ChatID = &chatID
	rDB.Status = h.C_RequestCreated
	err = mysql.GormDB.Save(&rDB).Error
	if err != nil {
		h.Err("Ошибка при сохранении rDB:\n" + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(responseJSON)
}

func initVerifications(invoiceID int64, requestItemID int64, verifiers []int64) bool {
	h.Log("initVerifications...")

	if requestItemID > 0 {
		var requestItem mysql.RequestItem
		err := mysql.GormDB.First(&requestItem, requestItemID).Error
		if err != nil {
			h.Err("Ошибка при получении mysql.RequestItem")
			return false
		}

		for _, verifierDivisionID := range verifiers {
			//получение отдела
			var divisionDB mysql.Division
			err = mysql.GormDB.First(&divisionDB, verifierDivisionID).Error
			if err != nil {
				h.Err("Ошибка при получении Отдела:" + err.Error())
				return false
			}

			verification := createVerification(
				requestItemID,
				h.C_ObjType_Request,
				requestItem.CreatedByID,
				divisionDB.ManagerID,
				h.C_Verification_Created,
				divisionDB.ID,
			)
			h.Log("Создан запрос на согласование:" + h.YellowColor + strconv.FormatInt(verification.ID, 10))
		}
	} else if invoiceID > 0 {
		var invoiceDB mysql.Invoice
		err := mysql.GormDB.First(&invoiceDB, invoiceID).Error
		if err != nil {
			h.Err("Ошибка при получении mysql.Invoice")
			return false
		}
		for _, verifierDivisionID := range verifiers {
			//получение отдела
			var divisionDB mysql.Division
			err := mysql.GormDB.First(&divisionDB, verifierDivisionID).Error
			if err != nil {
				h.Err("Ошибка при получении Отдела:" + err.Error())
				return false
			}

			verification := createVerification(
				invoiceID,
				h.C_ObjType_Invoice,
				invoiceDB.CreatedByID,
				divisionDB.ManagerID,
				h.C_Verification_Created,
				divisionDB.ID,
			)
			h.Log("Создан запрос на согласование:" + h.YellowColor + strconv.FormatInt(verification.ID, 10))
		}
	} else {
		return false
	}

	return true
}

func HandleInitVerifications(w http.ResponseWriter, r *http.Request) {
	h.Log("HandleInitVerifications...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	// Десериализация тела запроса в структуру VerifyEntityPayload
	var payload InitVerificationsPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.Err("Ошибка десериализации тела запроса:" + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var a []int64
	if payload.InvoiceID > 0 {
		a = append(a, c_invoiceValidators...)
	} else if payload.RequestItemID > 0 {
		a = append(a, c_requestVerifiers...)
	} else {
		h.Err("Ошибка: Должен быть или счёт или заказ")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	ok := initVerifications(payload.InvoiceID, payload.RequestItemID, a)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(responseJSON)
}

func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	h.Log("GetChatMessages...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	ChatID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return
	}

	db := mysql.GormDB

	countRequest := h.GetNamedStringFromURL(w, r, h.C_COUNT)
	fromDate := h.GetNamedStringFromURL(w, r, h.C_FROMDATE)
	toDate := h.GetNamedStringFromURL(w, r, h.C_TODATE)

	if countRequest != "" {
		var count int64
		result := db.Model(&mysql.Message{}).Where("chat_id = ?", ChatID).Count(&count)
		if result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Err("Ошибка при подсчете количества сообщений: " + result.Error.Error())
			return
		}

		countResponse := map[string]int64{
			"Count": count,
		}

		responseJSON, err := json.Marshal(countResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Err("Ошибка при преобразовании в JSON: " + err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
		return
	}

	if (fromDate != "") && (toDate != "") {
		var lastMessages []mysql.Message
		result := db.
			Where("chat_id = ? AND CT >= ? AND CT <= ?", ChatID, fromDate, toDate).
			Order("ct DESC").
			Preload("Author").
			Find(&lastMessages)
		if result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Err("Ошибка при получении последнего сообщения: " + result.Error.Error())
			return
		}

		var response []MessageResponse

		for _, m := range lastMessages {
			message := MessageResponse{
				ID:       m.ID,
				CT:       m.CT,
				Content:  m.Content,
				AuthorID: m.AuthorID,
				ChatID:   m.ChatID,
				Author:   m.Author.Name,
			}

			if nil != m.RepliedToID {
				message.RepliedToID = *m.RepliedToID
			}

			if nil != m.ForwardedID {
				message.ForwardedID = *m.ForwardedID
			}

			if nil != m.FileID {
				message.FileID = *m.FileID
			}

			response = append(response, message)
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Err("Ошибка при преобразовании в JSON: " + err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
		return
	}

	var messages []mysql.Message
	result := db.
		Preload("Author").
		Where("chat_id = ?", ChatID).
		Find(&messages)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Err("Ошибка при извлечении сообщений из базы данных: " + result.Error.Error())
		return
	}

	var response []MessageResponse

	for _, m := range messages {
		message := MessageResponse{
			ID:       m.ID,
			CT:       m.CT,
			Content:  m.Content,
			AuthorID: m.AuthorID,
			ChatID:   m.ChatID,
			Author:   m.Author.Name,
		}

		if nil != m.RepliedToID {
			message.RepliedToID = *m.RepliedToID
		}

		if nil != m.ForwardedID {
			message.ForwardedID = *m.ForwardedID
		}

		if nil != m.FileID {
			message.FileID = *m.FileID
		}

		response = append(response, message)
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		h.Err("Ошибка при преобразовании в JSON: " + err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func createVerification(
	objId int64,
	objType int,
	authorID int64,
	verificatorID int64,
	status int,
	divisionID int64,
) mysql.Verification {
	// Создание объекта Verification на основе данных из RequestItem
	verification := mysql.Verification{
		ObjID:        objId,
		ObjType:      objType,
		CreatedByID:  authorID,
		VerifiedByID: verificatorID,
		DivisionID:   divisionID,
		Status:       status,
	}
	// Сохранение объекта Verification в MySQL
	err := mysql.GormDB.Create(&verification).Error
	if err != nil {
		h.Err("Ошибка при создании документа Verification в MySQL: " + err.Error())
		return mysql.Verification{}
	}

	return verification
}

func HandleVerify(w http.ResponseWriter, r *http.Request) {
	h.Log("HandleVerify...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	var invoiceID, requestItemID, newStatus int
	var ok bool

	newStatus++
	invoiceID, ok = h.GetNamedIntFromURL(w, r, h.C_InvoiceID)
	if !ok {
		h.Log("missing invoiceID")
	}
	requestItemID, ok = h.GetNamedIntFromURL(w, r, h.C_RequestItemID)
	if !ok {
		h.Log("missing requestItemID")
	}
	// newStatus, ok = h.GetNamedIntFromURL(w, r, h.C_RequestItemID)
	if !ok {
		h.Err("некорректный статус")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existingVerification mysql.Verification

	var err error
	if invoiceID > 0 {
		h.Log("InvoiceID=", h.YellowColor, strconv.Itoa(invoiceID))
		// Получение счёта по InvoiceID
		var invoiceDB mysql.Invoice
		err = mysql.GormDB.First(&invoiceDB, requestItemID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				h.Err("mysql.Invoice не найден в MySQL")
				http.Error(w, "mysql.Invoice not found", http.StatusBadRequest)
			} else {
				h.Err("Ошибка при получении объекта mysql.Invoice из MySQL: " + err.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		invoiceDB.Status = h.C_InvoiceCancelled
		err = mysql.GormDB.First(&existingVerification).Error
	} else if requestItemID > 0 {
		h.Log("RequestItemID=", h.YellowColor, strconv.Itoa(requestItemID))
		// Получение заказа по RequestItemID
		var requestItem mysql.RequestItem
		err = mysql.GormDB.First(&requestItem, requestItemID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				h.Err("mysql.RequestItem не найден в MySQL")
				http.Error(w, "RequestItem not found", http.StatusBadRequest)
			} else {
				h.Err("Ошибка при получении объекта mysql.RequestItem из MySQL: " + err.Error())
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
		err = mysql.GormDB.First(&existingVerification).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			existingVerification = mysql.Verification{} // Сброс значения existingVerification, если документ не найден
		} else {
			h.Err("Ошибка при проверке существующего документа Verification в MySQL: " + err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Отправка успешного ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(responseJSON)
}

func GetEnums(w http.ResponseWriter, r *http.Request) {
	schetaStatuses := []string{"Отменён", "Черновик", "Создан", "В работе", "Решение Фин. Отеда", "Отправлен в бухгалтерию", "В оплате", "Частично оплачен", "Оплачен"}
	zakazStatuses := []string{"Отменён", "Черновик", "Создан", "В работе", "Закупка", "Оплата", "Завершён"}
	zayavkaStatuses := []string{"Отменена", "Черновик", "Создана", "В работе", "Закупка", "Оплата", "Завершена"}
	verificationStatuses := []string{"Поступило", "В работе", "Отклонено", "Согласовано"}
	chatIncExcl := []string{"undefined", "успех", "уже", "не существует"}

	schetaResponse := EnumResponse{Statuses: schetaStatuses}
	zakazResponse := EnumResponse{Statuses: zakazStatuses}
	zayavkaResponse := EnumResponse{Statuses: zayavkaStatuses}
	verificationResponse := EnumResponse{Statuses: verificationStatuses}
	chatStatusesResponse := EnumResponse{Statuses: chatIncExcl}

	response := map[string]EnumResponse{
		"Заявка":        zayavkaResponse,
		"Заказ":         zakazResponse,
		"Счёт":          schetaResponse,
		"Согласование":  verificationResponse,
		"ЧленствоВЧате": chatStatusesResponse,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Err("Ошибка при преобразовании в JSON: " + err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func GetUserDivisionsResponse(w http.ResponseWriter, r *http.Request) {
	h.Log("GetUserDivisionsResponse...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return
	}

	db := mysql.GormDB

	var divisions []mysql.Division
	result := db.
		Preload("Manager").
		Preload("Participants").
		Joins("JOIN division_user ON divisions.ID = division_user.division_id").
		Where("division_user.user_id = ?", ID).Find(&divisions)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Err("Ошибка при извлечении подразделений: " + result.Error.Error())
		return
	}

	var response []UserDivisionsResponse

	for _, division := range divisions {
		response = append(response, UserDivisionsResponse{
			ID:   division.ID,
			Name: division.Name,
		})
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Err("Ошибка при преобразовании в JSON: " + err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// передаётся ID RequestItem и время до которого интересует остаток
func calculateResidueInvoice(ID int, d time.Time) ResidueResponse {
	h.Log("calculateResidueInvoice...")

	//колчичество в заказе
	var totalRequestsQty float64
	result := mysql.GormDB.Model(&mysql.RequestItem{}).
		Select("SUM(qty) AS qty").
		Where("request_item_id = ? AND status > ?", ID, 1).
		Scan(&totalRequestsQty)
	if result.Error != nil {
		h.Err("Ошибка при запросе суммы заказов:" + result.Error.Error())
	}

	// сумма в рабочих счетах
	var totalInvoicesQty float64
	result = mysql.GormDB.Preload("Invoice").Model(&mysql.InvoiceItem{}).
		Select("SUM(qty) AS qty").
		Where("request_item_id = ? AND Invoice.status > ? AND ct <= ?", ID, 1, d).
		Scan(&totalInvoicesQty)
	if result.Error != nil {
		h.Err("Ошибка при запросе суммы позиций в счетах:" + result.Error.Error())
	}

	// Вычисляю остаток
	residue := totalRequestsQty - totalInvoicesQty

	r := ResidueResponse{
		Residue: residue,
		Total:   totalRequestsQty,
	}

	return r
}

// передаётся ID ARCWorkItem и время до которого интересует остаток
func calculateResidueRequestItem(ID int64, d time.Time) ResidueResponse {
	h.Log("calculateResidueRequestItem...")
	var totalArcWorkItemQty float64
	//ЛЗВ
	result := mysql.GormDB.Model(&mysql.ARCWorkItem{}).
		Select("SUM(qty) AS qty").
		Where("ID = ?", ID).
		Scan(&totalArcWorkItemQty)
	if result.Error != nil {
		h.Err("Ошибка при запросе суммы позиции в ЛЗВ: " + result.Error.Error())
	}
	h.Log("Сумма по ЛЗВ: " + h.YellowColor + strconv.FormatFloat(totalArcWorkItemQty, 'f', -1, 64))

	//все предыдущие заказы
	var totalRequestItemsQty float64
	result = mysql.GormDB.Model(&mysql.RequestItem{}).
		Select("SUM(qty) AS qty").
		Where("status > ? AND arc_work_item_id = ? AND ct <= ?", 1, ID, d).
		Scan(&totalRequestItemsQty)
	if result.Error != nil {
		h.Err("Ошибка при запросе суммы позиции в заказах: " + result.Error.Error())
	}

	// Вычисляю остаток
	residue := totalArcWorkItemQty - totalRequestItemsQty

	r := ResidueResponse{
		Residue: residue,
		Total:   totalArcWorkItemQty,
	}

	return r
}

func HandleGetResidueInvoice(w http.ResponseWriter, r *http.Request) {
	h.Log("GetResidueInvoice...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	IDStr := r.URL.Query().Get("ID")
	h.Log("ID=" + h.YellowColor + IDStr)

	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Err("Не удалось преобразовать ID в число:" + err.Error() + "\n" + IDStr)
		return
	}

	response := calculateResidueInvoice(ID, time.Now())

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Err("Ошибка сериализации response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func HandleGetResidueRequest(w http.ResponseWriter, r *http.Request) {
	h.Log("GetResidueRequest...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	IDStr := r.URL.Query().Get("ID")
	h.Log("ID=" + h.YellowColor + IDStr)

	ID, err := strconv.ParseInt(IDStr, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Err("Не удалось преобразовать ID в число: " + err.Error() + "\n" + IDStr)
		return
	}

	response := calculateResidueRequestItem(ID, time.Now())

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Err("Ошибка сериализации response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
