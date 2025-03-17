package dm_actions

import (
	"encoding/json"
	"net/http"
	"strconv"

	h "dm_server/dm_helper"

	mysql "dm_server/dm_db/dm_mysql"

	auth "dm_server/dm_api/dm_authorization"
)

type DocType int

const (
	DInvoice DocType = iota
	DRequestItem
)

// docType
// 0 - RequestItem
// 1 - Invoice
func GetVerificationStatusByUserIDItemID(v []*mysql.Verification, UserID int64, ItemID int64, docType DocType) (int, bool) {
	h.Log("GetVerificationStatusByUserID...")
	status := -1
	// TODO!!!! проверить валидность запрашиваемого ID типа сущности
	// if docType < h.C_ObjType_Request
	// {
	// 	h.Err("Неизвестный тип документа")
	// 	return status, false
	// }

	var verification mysql.Verification

	rows, errSQL := mysql.GormDB.
		Where("verified_by_id = ? AND obj_id = ? AND obj_type = ?", UserID, ItemID, docType).
		Find(&verification).
		Rows()

	if errSQL != nil {
		h.Err("Ошибка при выполнении запроса mysql.Verification:\n" + errSQL.Error())
		return status, false
	}
	defer rows.Close()

	var rowCount int = 0

	if rows.Next() {
		rowCount++

		if err := mysql.GormDB.ScanRows(rows, &verification); err != nil {
			h.Err("Ошибка при сканировании результатов запроса mysql.Verification:\n" + err.Error())
			return status, false
		}
	} else {
		h.Log("Пользователь " + h.YellowColor + strconv.FormatInt(UserID, 10) + h.DefaultColor + " не является обязательным согласователем")
		return status, true
	}

	if rowCount < 1 {
		h.Err("Нет записей, удовлетворяющих условию")
		return status, false
	}

	status = verification.Status

	return status, true
}

func GetVerifiersFromVerifications(verifications []*mysql.Verification) ([]VerifierResponse, bool) {
	h.Log("GetVerifiersFromVerifications...")

	var verifiers []VerifierResponse
	for _, v := range verifications {

		mysql.GormDB.Preload("VerifiedBy.Divisions").First(&v)
		var divisionID int64
		divisionName := "Без подразделения"
		if len(v.VerifiedBy.Divisions) > 0 {
			divisionID = v.VerifiedBy.Divisions[0].ID
			divisionName = v.VerifiedBy.Divisions[0].Name
		}

		// Preload cursor
		err := mysql.GormDB.
			First(&v).Error
		if err != nil {
			// Обработка ошибки загрузки деталей
			h.Err("Ошибка загрузки деталей v:\n" + err.Error())
			continue
		}

		verifier := VerifierResponse{
			ID:                 v.VerifiedByID,
			Name:               v.VerifiedBy.Name,
			DivisionID:         divisionID,
			DivisionName:       divisionName,
			VerificationStatus: v.Status,
		}

		verifiers = append(verifiers, verifier)
		h.Log("Согласователь добавлен: " + h.YellowColor + verifier.Name + "-" + verifier.DivisionName)
	}

	return verifiers, true
}

func handle_ListInvoice(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_ListInvoice...")

	var invoicesDB []mysql.Invoice
	resultSql := mysql.GormDB.
		Preload("Supplier").
		Preload("Object").
		Where("(Status >= ? ) OR (Status >= ?) AND (created_by_id = ?)", h.C_InvoiceCreated, h.C_InvoiceDraft, u.ID).
		Find(&invoicesDB)
	if resultSql.Error != nil {
		h.Err("Ошибка при извлечении списка счетов" + resultSql.Error.Error())
		return nil
	}

	var result []PliResult
	for _, iDB := range invoicesDB {

		i := PliResult{
			ID:               iDB.ID,
			CounteragentID:   iDB.Supplier.ID,
			CounteragentName: iDB.Supplier.Name,
			InvoiceNumber:    iDB.NumberStr,
			CT:               iDB.CT,
			ObjectID:         iDB.ObjectID,
			ObjectName:       iDB.Object.Name,
			Status:           iDB.Status,
		}

		if nil != iDB.ChatID {
			i.ChatID = *iDB.ChatID
		}
		if nil != iDB.FPNDS {
			i.FPWithNDS = *iDB.FPNDS
		}
		if iDB.Date != nil {
			i.InvoiceDate = *iDB.Date
		}

		result = append(result, i)
	}

	return result
}

func handle_ListStaff(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_ListStaff...")

	// список пользователей
	var usersDB []mysql.User
	resultSql := mysql.GormDB.
		Preload("Divisions").
		Preload("APIKeys").
		Find(&usersDB)
	if resultSql.Error != nil {
		h.Err("Ошибка при извлечении списка пользователей с подразделениями и ключами" + resultSql.Error.Error())
		return nil
	}

	var result []PageListStaffResponse
	for _, uDB := range usersDB {

		var divisionsResponse []UserDivisionsResponse

		for _, division := range uDB.Divisions {
			divisionsResponse = append(divisionsResponse, UserDivisionsResponse{
				ID:   division.ID,
				Name: division.Name,
			})
		}

		redmineUserID := "not set"
		if len(uDB.APIKeys) > 0 {
			redmineUserID = uDB.APIKeys[0].PartyID
		}

		u := PageListStaffResponse{
			ID:            uDB.ID,
			Name:          uDB.Name,
			Divisions:     divisionsResponse,
			RedmineUserID: redmineUserID,
			Email:         uDB.Email,
			Nickname:      uDB.Nickname,
			Phone:         uDB.Phone,
			CT:            uDB.CT,
		}

		// h.Log("Пользователь " + uDB.Name + " - " + strconv.Itoa(uDB.ID))

		result = append(result, u)
	}

	return result
}

func handle_ViewingOnePayments(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_ViewingOnePayments...")
	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return nil
	}

	var paymentDB mysql.Payment
	err := mysql.GormDB.
		Preload("PaymentOrder").
		Preload("CreatedBy").
		Preload("Payer").
		Preload("Acceptor").
		// Where("payment_order_id = ?", ID).
		First(&paymentDB, ID).Error
	if err != nil {
		h.Err("Не удалось найти mysql.RequestItem " + h.YellowColor + strconv.Itoa(ID) + "\n" + h.RedColor + err.Error())
		return nil
	}

	payer := PaymentCounteragentDetails{
		Name:                 paymentDB.Payer.Name,
		INN:                  paymentDB.Payer.Iinn,
		KPP:                  paymentDB.Payer.Kpp,
		BankName:             paymentDB.Payer.BankName,
		BankBIK:              paymentDB.Payer.BankBIK,
		BankAccountNumber:    paymentDB.Payer.BankAccountNumber,
		CorrespondentAccount: paymentDB.Payer.CorrespondentAccount,
	}

	supplier := PaymentCounteragentDetails{
		Name:                 paymentDB.Acceptor.Name,
		INN:                  paymentDB.Acceptor.Iinn,
		KPP:                  paymentDB.Acceptor.Kpp,
		BankName:             paymentDB.Acceptor.BankName,
		BankBIK:              paymentDB.Acceptor.BankBIK,
		BankAccountNumber:    paymentDB.Acceptor.BankAccountNumber,
		CorrespondentAccount: paymentDB.Acceptor.CorrespondentAccount,
	}

	result := ViewingOnePaymentResponse{
		ID:       paymentDB.ID,
		Summ:     paymentDB.Summ,
		Payer:    payer,
		Supplier: supplier,
	}

	if paymentDB.Date != nil {
		result.PaymentDate = *paymentDB.Date
	}
	return result
}

func handle_ListProvider(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_ListProvider...")

	// список заказов со статусом 4 (закупка) по объекту
	var counteragentsDB []mysql.Counteragent
	resultSql := mysql.GormDB.
		Find(&counteragentsDB)
	if resultSql.Error != nil {
		h.Err("Ошибка при извлечении контрагентов: " + resultSql.Error.Error())
		return nil
	}

	var result []PageListProviderResponse
	for _, cDB := range counteragentsDB {
		c := PageListProviderResponse{
			ID:          cDB.ID,
			Address:     cDB.LegalAddress,
			INN:         cDB.Iinn,
			ServiceType: cDB.ServiceType,
		}

		result = append(result, c)
	}

	return result
}

func handle_ListZakazZakupka(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_ListZakazZakupka...")

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ObjectID)
	if !ok {
		return nil
	}

	// список заказов со статусом 4 (закупка) по объекту
	var requestItemsDB []mysql.RequestItem
	resultSql := mysql.GormDB.
		Where("object_id = ? AND status = ?", ID, h.C_RequestPurchasing).Find(&requestItemsDB)
	if resultSql.Error != nil {
		h.Err("Ошибка при извлечении заказов: " + resultSql.Error.Error())
		return nil
	}

	var result []PageListZakazZakupkaResponse
	for _, riDB := range requestItemsDB {
		var residue ResidueResponse
		if nil != riDB.ARCWorkItemID {
			residue = calculateResidueRequestItem(*riDB.ARCWorkItemID, riDB.CT)
		} else {
			residue = ResidueResponse{}
		}

		if residue.Residue <= 0 {
			continue
		}

		r := PageListZakazZakupkaResponse{
			ID:                       riDB.ID,
			Name:                     riDB.Name,
			Residue:                  residue,
			Measure:                  riDB.Measure,
			SummMoneyTotalExpended:   0,
			SummMoneyResidueToExpend: 0,
			SummMoneyExpectedViaARC:  0,
			SummMoneyTotalPrediction: 0,
			RequestID:                riDB.RequestID,
			Qty:                      riDB.Qty,
			Note:                     riDB.Note,
		}

		if nil != riDB.ARCWorkItemID {
			r.ARCWorkItemID = *riDB.ARCWorkItemID
		}

		if nil != riDB.RequiredAtDateTime {
			r.RequiredAtDateTime = *riDB.RequiredAtDateTime
		}
		result = append(result, r)
	}

	return result
}

// просмотр заявки
func handle_ViewingZayavka(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_ViewingInvoice...")

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return nil
	}

	var rDB mysql.Request
	err := mysql.GormDB.
		Preload("ARCWork").
		Preload("Project").
		Preload("Object").
		Preload("RequestItems").
		First(&rDB, ID).Error
	if err != nil {
		h.Err("Не удалось найти mysql.Request " + h.YellowColor + strconv.Itoa(ID) + "\n" + h.RedColor + err.Error())
		return nil
	}
	var risResult []RequestedItemViewZayavkaResponse
	var workName string
	for _, riDB := range rDB.RequestItems {
		// Preload cursor
		if riDB.IsAnalogueSet {
			err = mysql.GormDB.
				Preload("ARCWork").
				Preload("ARCWorkItem").
				First(&riDB).Error
		} else {
			err = mysql.GormDB.
				Preload("ARCWork").
				Preload("ARCWorkItem").
				First(&riDB).Error
		}
		if err != nil {
			// Обработка ошибки загрузки деталей
			h.Err("Ошибка загрузки деталей v:\n" + err.Error())
			continue
		}
		workName = ""
		if nil != riDB.ARCWork {
			workName = riDB.ARCWork.Name
		}
		ri := RequestedItemViewZayavkaResponse{
			ID:      riDB.ID,
			Name:    riDB.Name,
			Qty:     riDB.Qty,
			Measure: riDB.Measure,
			Status:  riDB.Status,
		}
		if nil != riDB.ARCWorkItem {
			ri.ARCWorkItemName = riDB.ARCWorkItem.Name
			ri.ARCWorkItemID = riDB.ARCWorkItem.ID
		}
		if nil != riDB.AnalogueID {
			ri.AnalogueID = *riDB.AnalogueID
		}
		if nil != riDB.RequiredAtDateTime {
			ri.RequiredAtDateTime = *riDB.RequiredAtDateTime
		}
		ri.AnalogueInitialName = ""
		if riDB.IsAnalogueSet {
			if nil != riDB.ARCWorkItem {
				ri.AnalogueInitialName = riDB.ARCWorkItem.Name
			}
		}

		var residue ResidueResponse
		if nil != riDB.ARCWorkItemID {
			residue = calculateResidueRequestItem(*riDB.ARCWorkItemID, riDB.CT)
		} else {
			residue = ResidueResponse{}
		}

		excess := residue.Residue < 0

		var percs InfoPercentages
		if residue.Total != 0 {
			percs.QtyOrdered = residue.Residue / residue.Total
			percs.QtyDelivered = .1
			percs.QtyPayed = .1
		}

		ri.Excess = excess
		ri.Percentages = percs
		ri.Note = riDB.Note
		ri.Residue = residue

		risResult = append(risResult, ri)
	}
	result := PageViewingZayavkaResponse{
		CT:           rDB.CT,
		RequestItems: risResult,
		ProjectName:  rDB.Project.Name,
		ObjectID:     rDB.ObjectID,
		ObjectName:   rDB.Object.Name,
		WorkName:     workName,
		Status:       rDB.Status,
	}

	if nil != rDB.ARCWorkID {
		result.ARCWorkID = *rDB.ARCWorkID
	} else {
		result.ARCWorkID = 0
	}

	if nil != rDB.ARCID {
		result.ARCID = *rDB.ARCID
	} else {
		result.ARCID = 0
	}

	return result
}

// просмотр счёта
func handle_ViewingInvoice(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_ViewingInvoice...")

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return nil
	}

	var iDB mysql.Invoice
	err := mysql.GormDB.
		Preload("CreatedBy").
		Preload("Supplier").
		Preload("Customer").
		Preload("Payer").
		Preload("Project").
		Preload("Object").
		Preload("LastModifiedBy").
		Preload("Verifications").
		Preload("Verifications.VerifiedBy").
		Preload("InvoiceItems").
		Preload("Payments").
		Preload("Payments.PaymentOrder").
		First(&iDB, ID).Error
	if err != nil {
		h.Err("Не удалось найти mysql.Invoice " + h.YellowColor + strconv.Itoa(ID) + "\n" + h.RedColor + err.Error())
		return nil
	}
	h.Log("Загружен mysql.Invoice")

	verifiers, ok := GetVerifiersFromVerifications(iDB.Verifications)
	if !ok {
		h.Log("Ошибка при получении списка согласователей для счёта ID: " + strconv.Itoa(ID))
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	// статус согласования у текущего пользователя
	userStatus, vOk := GetVerificationStatusByUserIDItemID(iDB.Verifications, u.ID, iDB.ID, h.C_ObjType_Invoice)
	if !vOk {
		h.Log("Ошибка при вызове GetVerificationStatusByUserIDItemID")
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	// строки в счёте
	var iis []InvoiceItemResponse
	for _, iiDB := range iDB.InvoiceItems {
		// Preload cursor
		err := mysql.GormDB.
			Preload("ARC").
			Preload("RequestItem").
			First(&iiDB).Error
		if err != nil {
			// Обработка ошибки загрузки деталей
			h.Err("Ошибка загрузки деталей v:\n" + err.Error())
			continue
		}

		ii := InvoiceItemResponse{
			Name:               iiDB.RequestItem.Name,
			Measure:            iiDB.Measure,
			TotalPrice:         iiDB.UPWithNDS,
			QtyInvoice:         iiDB.Qty,
			QtyRequest:         iiDB.RequestItem.Qty,
			InvoiceItemID:      iiDB.ID,
			RequestItemID:      iiDB.RequestItemID,
			RequiredAtDateTime: *iiDB.RequiredAtDateTime,
			DiscountType:       iiDB.DiscountType,
			UPWithoutNDS:       iiDB.UPWithoutNDS,
			UPWithNDS:          iiDB.UPWithNDS,
			UPWithDiscount:     iiDB.UPWithDiscount,
		}

		if nil != iiDB.DiscountValue {
			ii.DiscountValue = *iiDB.DiscountValue
		}
		if nil != iiDB.ARC {
			ii.ARCID = *iiDB.ARCID
		}

		if nil != iiDB.RequiredAtDateTime {
			ii.RequiredAtDateTime = *iiDB.RequiredAtDateTime
		}

		iis = append(iis, ii)
	}

	// ответ клиенту
	result := PageViewingInvoiceResponse{
		CT:            iDB.CT,
		CreatedByName: iDB.CreatedBy.Name,
		ProjectName:   iDB.Project.Name,
		ObjectID:      iDB.ObjectID,
		ObjectName:    iDB.Object.Name,
		Verifiers:     verifiers,
		Status:        iDB.Status,
		UserStatus:    userStatus,
		Items:         iis,
		Date:          *iDB.Date,
		//!!!TODO
		// Verifications:      invoiceDB.Verifications,
		Payments:      iDB.Payments,
		SupplierID:    iDB.SupplierID,
		SupplierName:  iDB.Supplier.Name,
		CustomerID:    iDB.CustomerID,
		CustomerName:  iDB.Customer.Name,
		InvoiceNumber: iDB.NumberStr,
		CreatedByID:   iDB.CreatedByID,
		PayerID:       iDB.PayerID,
		PayerName:     iDB.Payer.Name,
		MeasureWeight: iDB.MeasureWeight,
	}

	if nil != iDB.ChatID {
		result.ChatID = *iDB.ChatID
	}
	if nil != iDB.Date {
		result.Date = *iDB.Date
	}
	if nil != iDB.FPWithDiscount {
		result.FPWithDiscount = *iDB.FPWithDiscount
	}
	if nil != iDB.FPNDS {
		result.FPNDS = *iDB.FPNDS
	}
	if nil != iDB.FPWithNDS {
		result.FPWithNDS = *iDB.FPWithNDS
	}
	if nil != iDB.FPWithNDS {
		result.FPWithoutNDS = *iDB.FPWithoutNDS
	}
	if nil != iDB.TotalWeight {
		result.TotalWeight = *iDB.TotalWeight
	}

	return result
	// return response
}

func GetAvailableARCWorks(u *mysql.User) ([]mysql.ARCWork, bool) {
	h.Log("GetAvailableARCWorks...")

	var arcWorksDB []mysql.ARCWork

	resultGORM := mysql.GormDB.
		Preload("ARC").
		Preload("ARC.Object").
		Find(&arcWorksDB)
	if resultGORM.Error != nil {
		h.Err("Ошибка при извлечении ЛЗВ: " + resultGORM.Error.Error())
		return nil, false
	}

	return arcWorksDB, true
}

func GetAvailableObjectsWithARC(u *mysql.User) ([]*mysql.Object, bool) {
	h.Log("GetAvailableObjectsWithARC...")

	var arcsDB []mysql.ARC

	resultGORM := mysql.GormDB.
		Preload("Object").
		Find(&arcsDB)
	if resultGORM.Error != nil {
		h.Err("Ошибка при извлечении ЛЗВ: " + resultGORM.Error.Error())
		return nil, false
	}

	var result []*mysql.Object

	for _, a := range arcsDB {
		result = append(result, a.Object)
	}

	return result, true
}

func handle_CreateRequest(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_CreateRequest...")

	filterStr := h.GetNamedStringFromURL(w, r, h.C_Str_Field_Filter)
	h.Log(h.C_Str_Field_Filter, "=", filterStr)

	id, idOK := h.TryGetNamedIntFromURL(w, r, h.C_ID)
	if !idOK {
		h.Log("Работа явно не указана")
	}

	var works []CRCard
	if id > 0 {
		var arcWorkDB mysql.ARCWork

		resultGORM := mysql.GormDB.
			Preload("ARC").
			Preload("ARC.Object").
			Find(&arcWorkDB, id)
		if resultGORM.Error != nil {
			h.Err("Ошибка при извлечении ЛЗВ: " + resultGORM.Error.Error())
			return nil
		}

		work := CRCard{
			ObjectID:    arcWorkDB.ARC.ObjectID,
			ObjectName:  arcWorkDB.ARC.Object.Name,
			ARCWorkName: arcWorkDB.Name,
			ARCCT:       arcWorkDB.ARC.CT,
			Qty:         arcWorkDB.Qty,
			Measure:     arcWorkDB.Measure,
			Address:     arcWorkDB.ARC.Object.Address,
			ID:          arcWorkDB.ID,
			CODE:        arcWorkDB.Code,
		}

		works = append(works, work)
	} else {
		worksDB, ok := GetAvailableARCWorks(u)
		if !ok {
			return nil
		}
		for _, w := range worksDB {
			work := CRCard{
				ObjectID:    w.ARC.ObjectID,
				ObjectName:  w.ARC.Object.Name,
				ARCWorkName: w.Name,
				ARCCT:       w.ARC.CT,
				Qty:         w.Qty,
				Measure:     w.Measure,
				Address:     w.ARC.Object.Address,
				ID:          w.ID,
				CODE:        w.Code,
			}
			works = append(works, work)
		}
	}

	response := PageCreateRequestResponse{
		List: works,
	}

	h.Log("Данные для страницы ", h.YellowColor, "CreateRequest", h.DefaultColor, " отправлены.")

	return response
}

// стартовая страница
func handle_PageWorkMobile(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_PageWorkMobile...")

	var projectsDB []*mysql.Project
	err := mysql.GormDB.
		Find(&projectsDB).Error
	if err != nil {
		h.Log("Ошибка при получении списка проектов:\n" + err.Error())
		return nil
	}
	var projects []PwmProject
	for _, prj := range projectsDB {
		p := PwmProject{
			ID:      prj.ID,
			Name:    prj.Name,
			Address: prj.Address,
			Price:   prj.Price,
		}

		if nil != prj.CompletionDate {
			p.CompletionDate = *prj.CompletionDate
		}

		projects = append(projects, p)
	}
	h.Log("Загружены проекты")

	var objectsDB []*mysql.Object
	err = mysql.GormDB.
		Preload("Project").
		Find(&objectsDB).Error
	if err != nil {
		h.Log("Ошибка при получении списка проектов:\n" + err.Error())
		return nil
	}
	var objects []PwmObject
	for _, obj := range objectsDB {
		o := PwmObject{
			ID:          obj.ID,
			Name:        obj.Name,
			ProjectName: obj.Project.Name,
			Address:     obj.Address,
		}

		if nil != obj.CompletionDate {
			o.CompletionDate = *obj.CompletionDate
		}

		objects = append(objects, o)
	}
	h.Log("Загружены объекты")

	var requestsDB []mysql.Request
	err = mysql.GormDB.
		Preload("Project").
		Preload("Object").
		Preload("ARC").
		Preload("ARCWork").
		Find(&requestsDB).Error
	if err != nil {
		h.Err("Ошибка при получении списка проектов:\n" + err.Error())
		return nil
	}
	var requests []PwmRequest
	var allRis []PwmRequestItem
	for _, req := range requestsDB {
		var requestItemsDB []*mysql.RequestItem
		err = mysql.GormDB.
			Preload("ARCWorkItem").
			Where("request_id = ?", req.ID).
			Find(&requestItemsDB).Error
		if err != nil {
			h.Err("Ошибка при получении списка проектов:\n" + err.Error())
			return nil
		}
		projectName := req.Project.Name
		objectName := req.Object.Name
		var ris []PwmRequestItem
		for _, riDB := range requestItemsDB {
			var residue ResidueResponse
			if nil != riDB.ARCWorkItemID {
				residue = calculateResidueRequestItem(*riDB.ARCWorkItemID, riDB.CT)
			} else {
				residue = ResidueResponse{}
			}
			excess := residue.Residue < 0

			var percs InfoPercentages
			if residue.Total != 0 {
				percs.QtyOrdered = residue.Residue / residue.Total
				percs.QtyDelivered = .1
				percs.QtyPayed = .1
			}

			ri := PwmRequestItem{
				ID:          riDB.ID,
				Excess:      excess,
				Name:        riDB.Name,
				Qty:         riDB.Qty,
				Measure:     riDB.Measure,
				Note:        riDB.Note,
				ProjectName: projectName,
				ObjectID:    riDB.ObjectID,
				ObjectName:  objectName,
				Status:      riDB.Status,
				Percentages: percs,
			}
			if nil != riDB.ARCWorkItem {
				ri.ARCWorkItemID = riDB.ARCWorkItem.ID
				ri.ARCWorkItemName = riDB.ARCWorkItem.Name
			}
			if nil != riDB.RequiredAtDateTime {
				ri.RequiredAtDateTime = *riDB.RequiredAtDateTime
			}

			ris = append(ris, ri)
			allRis = append(allRis, ri)
		}

		r := PwmRequest{
			ID:           req.ID,
			Status:       req.Status,
			CT:           req.CT,
			ProjectName:  req.Project.Name,
			ObjectName:   req.Object.Name,
			RequestItems: ris,
		}
		if nil != req.ARCWork {
			r.WorkName = req.ARCWork.Name
		}
		if req.ChatID != nil {
			r.ChatID = *req.ChatID
		}

		requests = append(requests, r)
	}
	h.Log("Загружены заявки")

	var invoices []*mysql.Invoice
	err = mysql.GormDB.
		Preload("Verifications").
		Preload("Supplier").
		Preload("Customer").
		Preload("Payer").
		Preload("InvoiceItems").
		Preload("Verifications").
		Find(&invoices).Error
	if err != nil {
		h.Log("Ошибка при получении списка счетов:\n" + err.Error())
		return nil
	}
	h.Log("Загружены счета")

	response := PageWorkMobileResponse{
		Projects:     projects,
		Objects:      objects,
		Requests:     requests,
		RequestItems: allRis,
		Invoices:     invoices,
	}

	h.Log("Данные для страницы " + h.YellowColor + "PageWorkMobile" + h.DefaultColor + " отправлены.")

	return response
}

func handle_ListZayavka(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_ListZayavka...")

	var requestsDB []*mysql.Request
	err := mysql.GormDB.
		Preload("Project").
		Preload("Object").
		Preload("ARC").
		Preload("ARCWork").
		Find(&requestsDB).Error
	if err != nil {
		h.Log("Ошибка при получении списка проектов:\n" + err.Error())
		return nil
	}
	h.Log("Загружены заявки")

	var requests []PwmRequest
	for _, req := range requestsDB {
		var requestItemsDB []*mysql.RequestItem
		err = mysql.GormDB.
			Preload("ARCWorkItem").
			Where("request_id = ?", req.ID).
			Find(&requestItemsDB).Error
		if err != nil {
			h.Err("Ошибка при получении списка проектов:\n" + err.Error())
			return nil
		}
		projectName := req.Project.Name
		objID := req.ObjectID
		objectName := req.Object.Name
		var ris []PwmRequestItem
		for _, ri := range requestItemsDB {
			var residue ResidueResponse
			if nil != ri.ARCWorkItemID {
				residue = calculateResidueRequestItem(*ri.ARCWorkItemID, ri.CT)
			} else {
				residue = ResidueResponse{}
			}
			excess := residue.Residue < 0

			var percs InfoPercentages
			if residue.Total != 0 {
				percs.QtyOrdered = residue.Residue / residue.Total
				percs.QtyDelivered = .1
				percs.QtyPayed = .1
			}

			riResult := PwmRequestItem{
				ID:          ri.ID,
				Excess:      excess,
				Qty:         ri.Qty,
				Measure:     ri.Measure,
				Note:        ri.Note,
				ProjectName: projectName,
				ObjectID:    objID,
				ObjectName:  objectName,
				Status:      ri.Status,
				Percentages: percs,
				Name:        ri.Name,
			}

			if nil != ri.ARCWorkItem {
				riResult.ARCWorkItemID = ri.ARCWorkItem.ID
				riResult.ARCWorkItemName = ri.ARCWorkItem.Name
			}

			if nil != ri.RequiredAtDateTime {
				riResult.RequiredAtDateTime = *ri.RequiredAtDateTime
			}
			ris = append(ris, riResult)
		}

		r := PwmRequest{
			ID:           req.ID,
			Status:       req.Status,
			CT:           req.CT,
			ProjectName:  req.Project.Name,
			ObjectID:     req.ObjectID,
			ObjectName:   req.Object.Name,
			RequestItems: ris,
		}
		if nil != req.ARCWork {
			r.WorkName = req.ARCWork.Name
		}
		if req.ChatID != nil {
			r.ChatID = *req.ChatID
		}

		requests = append(requests, r)
	}
	h.Log("Загружены заявки")

	response := PageWorListZayavkaResponse{
		Requests: requests,
	}
	h.Log("Данные для страницы " + h.YellowColor + "ListZayavka" + h.DefaultColor + " отправлены.")
	return response
}

// просмотр заказа
func handle_ViewingZakaz(w http.ResponseWriter, r *http.Request, u *mysql.User) interface{} {
	h.Log("handle_ViewingZakaz...")

	ID, ok := h.GetNamedIntFromURL(w, r, h.C_ID)
	if !ok {
		return nil
	}

	var riDB mysql.RequestItem
	err := mysql.GormDB.Preload("Project").
		Preload("Object").
		Preload("ARCWorkItem").
		Preload("ARCWork").
		Preload("Request").
		Preload("Verifications").
		Preload("Verifications.VerifiedBy").
		Preload("CreatedBy").
		First(&riDB, ID).Error
	if err != nil {
		h.Err("Не удалось найти RequestItem " + h.YellowColor + strconv.Itoa(ID) + "\n" + h.RedColor + err.Error())
		return nil
	}
	h.Log("Загружен RequestItem")

	verifiers, ok := GetVerifiersFromVerifications(riDB.Verifications)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	rir := RequestedItemResponse{
		Name:       riDB.Name,
		Qty:        riDB.Qty,
		IsAnalogue: riDB.IsAnalogueSet,
		Measure:    riDB.Measure,
	}

	rir.AnalogueInitialName = ""
	if riDB.IsAnalogueSet {
		if nil != riDB.ARCWorkItem {
			rir.AnalogueInitialName = riDB.ARCWorkItem.Name
		}
	}
	if nil != riDB.RequiredAtDateTime {
		rir.RequiredAtDateTime = *riDB.RequiredAtDateTime
	}

	h.Log("Информация о заказе сформирована: ", h.YellowColor, strconv.FormatInt(riDB.ID, 10))

	//todo - возможна инвариантность с датой, сейчас используется дата заказа
	var residue ResidueResponse
	if nil != riDB.ARCWorkItemID {
		residue = calculateResidueRequestItem(*riDB.ARCWorkItemID, riDB.CT)
	} else {
		residue = ResidueResponse{}
	}

	excess := residue.Residue < 0

	var percs InfoPercentages
	if residue.Total != 0 {
		percs.QtyOrdered = residue.Residue / residue.Total
		percs.QtyDelivered = .1
		percs.QtyPayed = .1
	}

	// статус согласования у текущего пользователя
	userStatus, vOk := GetVerificationStatusByUserIDItemID(riDB.Verifications, u.ID, riDB.ID, DInvoice)
	if !vOk {
		h.Log("Ошибка при вызове GetVerificationStatusByUserID")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// ответ клиенту
	response := PageViewingZakazResponse{
		CreatedDateTime: riDB.CT,
		CreatedByName:   u.Name,
		ProjectName:     riDB.Project.Name,
		ObjectName:      riDB.Object.Name,
		Verifiers:       verifiers,
		RequestedItem:   rir,
		ARCTotal:        residue.Total,
		ARCResidue:      residue.Residue,
		Excess:          excess,
		UserStatus:      userStatus,
		Status:          riDB.Status,
		Percentages:     percs,
		Note:            riDB.Note,
		CreatedByID:     riDB.CreatedByID,
		ObjectID:        riDB.ObjectID,
	}

	if nil != riDB.ARCWorkID {
		response.ARCWorkID = *riDB.ARCWorkID
	} else {
		response.ARCWorkID = 0
	}

	if nil != riDB.ARCID {
		response.ARCID = *riDB.ARCID
	} else {
		response.ARCID = 0
	}

	if nil != riDB.ARCWork {
		response.WorkName = riDB.ARCWork.Name
	} else {
		response.WorkName = ""
	}

	if riDB.Request.ChatID != nil {
		response.ChatID = *riDB.Request.ChatID
	}

	h.Log("Данные для страницы " + h.YellowColor + "ViewingZakaz" + h.DefaultColor + " отправлены.")

	return response
}

func GetPageData(w http.ResponseWriter, r *http.Request) {
	h.Log("GetPageData...")

	u := auth.GetAuthUser(w, r)
	if nil == u {
		return
	}

	page := r.URL.Query().Get("page")
	if len(page) > 0 {
		var result interface{}
		switch page {
		case "CreateRequest":
			{
				result = handle_CreateRequest(w, r, u)
				break
			}
		case "PageWorkMobile":
			{
				result = handle_PageWorkMobile(w, r, u)
				break
			}
		case "ListZayavka":
			{
				result = handle_ListZayavka(w, r, u)
				break
			}
		case "ViewingZakaz":
			{
				result = handle_ViewingZakaz(w, r, u)
				break
			}
		case "ViewingInvoice":
			{
				result = handle_ViewingInvoice(w, r, u)
				break
			}
		case "ViewingZayavka":
			{
				result = handle_ViewingZayavka(w, r, u)
				break
			}
		case "ListZakazZakupka":
			{
				result = handle_ListZakazZakupka(w, r, u)
				break
			}
		case "ListProvider":
			{
				result = handle_ListProvider(w, r, u)
				break
			}
		case "ChatDetails":
			{
				result = handle_ChatDetails(w, r, u)
				break
			}
		case "ViewingOnePayments":
			{
				result = handle_ViewingOnePayments(w, r, u)
				break
			}
		case "ListStaff":
			{
				result = handle_ListStaff(w, r, u)
				break
			}
		case "ListInvoice":
			{
				result = handle_ListInvoice(w, r, u)
				break
			}
		}
		resultWithoutNan := h.ReplaceNaN(result)
		jsonBytes, err := json.Marshal(resultWithoutNan)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// w.Write(responseJSON)
		w.Write(jsonBytes)
		h.Log("Данные для страницы " + h.YellowColor + page + h.DefaultColor + " отправлены.")
		return
	}

	http.Error(w, "Неизвестная страница", http.StatusBadRequest)
}
