package dm_mongo

import (
	"time"
)

type Counteragent struct {
	ID                   int    `json:"ID" bson:"ID"`
	Name                 string `json:"Name" bson:"Name"`
	Shortname            string `json:"Shortname" bson:"Shortname"`
	ServiceType          string `json:"ServiceType" bson:"ServiceType"`
	LegalAddress         string `json:"LegalAddress" bson:"LegalAddress"`
	PhysicalAddress      string `json:"PhysicalAddress" bson:"PhysicalAddress"`
	Iinn                 string `json:"Iinn" bson:"Iinn"`
	Ogrn                 string `json:"Ogrn" bson:"Ogrn"`
	Kpp                  string `json:"Kpp" bson:"Kpp"`
	BankAccountNumber    string `json:"BankAccountNumber" bson:"BankAccountNumber"`
	CorrespondentAccount string `json:"CorrespondentAccount" bson:"CorrespondentAccount"`
	BankBIK              string `json:"BankBIK" bson:"BankBIK"`
	BankName             string `json:"BankName" bson:"BankName"`
	BankAddress          string `json:"BankAddress" bson:"BankAddress"`
	PhoneNumber          string `json:"PhoneNumber" bson:"PhoneNumber"`
	ContactEmail         string `json:"ContactEmail" bson:"ContactEmail"`
	Website              string `json:"Website" bson:"Website"`
	CEO                  string `json:"CEO" bson:"CEO"`
	Acting_basis         string `json:"Acting_basis" bson:"Acting_basis"`
	Deputy_CEO           string `json:"Deputy_CEO" bson:"Deputy_CEO"`
	Chief_accountant     string `json:"Chief_accountant" bson:"Chief_accountant"`
}

type Invoice struct {
	ID            int            `json:"ID" bson:"ID"`
	NumberStr     string         `json:"NumberStr" bson:"NumberStr"`
	Date          time.Time      `json:"Date" bson:"Date"`
	SupplierID    int            `json:"SupplierID" bson:"SupplierID"`
	CustomerID    int            `json:"CustomerID" bson:"CustomerID"`
	PayerID       int            `json:"PayerID" bson:"PayerID"`
	Status        int            `json:"Status" bson:"Status"`
	ProjctID      int            `json:"ProjectID" bson:"ProjectID"`
	ObjectID      int            `json:"ObjectID" bson:"ObjectID"`
	FPWithoutNDS  float64        `json:"FPWithoutNDS" bson:"FPWithoutNDS"`
	FPWithNDS     float64        `json:"FPWithNDS" bson:"FPWithNDS"`
	FPNDS         float64        `json:"FPNDS" bson:"FPNDS"`
	TotalWeight   float64        `json:"TotalWeight" bson:"TotalWeight"`
	MeasureWeight string         `json:"MeasureWeight" bson:"MeasureWeight"`
	Verifications []Verification `json:"Verifications" bson:"Verifications"`
	InvoiceItems  []InvoiceItem  `json:"InvoiceItems" bson:"InvoiceItems"`
	ChatID        int            `json:"ChatID" bson:"ChatID"`
}

// строка в счёте
// UP Unit Price
// TC Total Cost
type InvoiceItem struct {
	ID                 int       `json:"ID" bson:"ID"`
	InvoiceID          int       `json:"InvoiceID" bson:"InvoiceID"`
	ProjectID          int       `json:"ProjectID" bson:"ProjectID"`
	ObjectID           int       `json:"ObjectID" bson:"ObjectID"`
	ARCID              int       `json:"ARCID" bson:"ARCID"`
	ARCWorkItemID      int       `json:"ARCWorkItemID" bson:"ARCWorkItemID"`
	RequestID          int       `json:"RequestID" bson:"RequestID"`
	RequestItemID      int       `json:"RequestItemID" bson:"RequestItemID"`
	Asset              string    `json:"Asset" bson:"Asset"`
	Measure            string    `json:"Measure" bson:"Measure"`
	Qty                float64   `json:"Qty" bson:"Qty"`
	UPWithoutNDS       float64   `json:"UPWithoutNDS" bson:"UPWithoutNDS"`
	UPWithNDS          float64   `json:"UPWithNDS" bson:"UPWithNDS"`
	UPWithoutDiscount  float64   `json:"UPWithoutDiscount" bson:"UPWithoutDiscount"`
	UPWithDiscount     float64   `json:"UPWithDiscount" bson:"UPWithDiscount"`
	TCWithoutNDS       float64   `json:"TCWithoutNDS" bson:"TCWithoutNDS"`
	TCWithNDS          float64   `json:"TCWithNDS" bson:"TCWithNDS"`
	TCWithoutDiscount  float64   `json:"TCWithoutDiscount" bson:"TCWithoutDiscount"`
	TCWithDiscount     float64   `json:"TCWithDiscount" bson:"TCWithDiscount"`
	RequiredAtDateTime time.Time `json:"RequiredAtDateTime" bson:"RequiredAtDateTime"`
}

type ARCWorkItem struct {
	ID        int     `json:"ID" bson:"ID"`
	Code      string  `json:"Code" bson:"Code"`
	ARCWorkID int     `json:"ARCWorkID" bson:"ARCWorkID"`
	Name      string  `json:"Name" bson:"Name"`
	Measure   string  `json:"Measure" bson:"Measure"`
	Qty       float64 `json:"Qty" bson:"Qty"`
	CPWMU     string  `json:"CPW" bson:"CPW"` //consumption per Work Measurement Unit - расход на единицу измерения работы
	Note      string  `json:"Note" bson:"Note"`
}

type ARCWork struct {
	ID           int           `json:"ID" bson:"ID"`
	Code         string        `json:"Code" bson:"Code"`
	ARCID        int           `json:"ARCID" bson:"ARCID"`
	ProjectID    int           `json:"ProjectID" bson:"ProjectID"`
	ProjectName  string        `json:"ProjectName" bson:"ProjectName"`
	ObjectID     int           `json:"ObjectID" bson:"ObjectID"`
	ObjectName   string        `json:"ObjectName" bson:"ObjectName"`
	Name         string        `json:"Name" bson:"Name"`
	ACRWorkItems []ARCWorkItem `json:"ACRWorkItemsList" bson:"ACRWorkItemsList"`
	Qty          float64       `json:"Qty" bson:"Qty"`
	Measure      string        `json:"Measure" bson:"Measure"`
}

type ARC struct {
	ID                   int       `json:"ID" bson:"ID"`
	CreatedByID          int       `json:"CreatedBy" bson:"CreatedBy"`
	CreatedByName        string    `json:"CreatedByName" bson:"CreatedByName"`
	CreatedDateTime      time.Time `json:"CreatedDateTime" bson:"CreatedDateTime"`
	LastModifiedByID     int       `json:"LastModifiedBy" bson:"LastModifiedBy"`
	LastModifiedByName   string    `json:"LastModifiedByName" bson:"LastModifiedByName"`
	LastModifiedDateTime time.Time `json:"LastModifiedDateTime" bson:"LastModifiedDateTime"`
	ProjectID            int       `json:"ProjectID" bson:"ProjectID"`
	ProjectName          string    `json:"ProjectName" bson:"ProjectName"`
	ObjectID             int       `json:"ObjectID" bson:"ObjectID"`
	ObjectName           string    `json:"ObjectName" bson:"ObjectName"`
	Name                 string    `json:"Name" bson:"Name"`
	ACRWorksList         []ARCWork `json:"ACRItemsList" bson:"ACRItemsList"`
}

type RequestItem struct {
	ID                   int            `json:"ID" bson:"ID"`
	Type                 int            `json:"Type" bson:"Type"`     // 0 - ? 1 - материалы, 2 - механизмы, 3 - трудозатраты, 4 - накладные расходы
	Status               int            `json:"Status" bson:"Status"` //0 — отменён, 1 — черновик,  2 — создан, 3 — в работе, 4 — закупка, 5 — оплата, 6 — завершён
	RequestID            int            `json:"RequestID" bson:"RequestID"`
	ObjectID             int            `json:"ObjectID" bson:"ObjectID"`
	ObjectName           string         `json:"ObjectName" bson:"ObjectName"`
	ProjectID            int            `json:"ProjectID" bson:"ProjectID"`
	ProjectName          string         `json:"ProjectName" bson:"ProjectName"`
	LastModifiedByID     string         `json:"LastModifiedByID" bson:"LastModifiedByID"`
	LastModifiedByName   string         `json:"LastModifiedByName" bson:"LastModifiedByName"`
	LastModifiedDateTime time.Time      `json:"LastModifiedDateTime" bson:"LastModifiedDateTime"`
	ARCID                int            `json:"ARCID" bson:"ARCID"`
	ARCWorkItemID        int            `json:"ARCWorkItemID" bson:"ARCWorkItemID"`
	ARCWowkItemName      string         `json:"ARCWowkItemName" bson:"ARCWowkItemName"`
	ARCWorkID            int            `json:"ARCWorkID" bson:"ARCWorkID"`
	ARCWorkName          string         `json:"ARCWorkName" bson:"ARCWorkName"`
	Measure              string         `json:"Measure" bson:"Measure"`
	Name                 string         `json:"Name" bson:"Name"`
	Qty                  float64        `json:"Qty" bson:"Qty"`
	QtyOrdered           float64        `json:"QtyOrdered" bson:"QtyOrdered"`
	QtyDelivered         float64        `json:"QtyDelivered" bson:"QtyDelivered"`
	QtyPayed             float64        `json:"QtyPayed" bson:"QtyPayed"`
	IsAnalogueSet        bool           `json:"IsSet" bson:"IsSet"`
	Analogue             Analogue       `json:"Analogue" bson:"Analogue"`
	Note                 string         `json:"Note" bson:"Note"`
	RequiredAtDateTime   time.Time      `json:"RequiredAtDateTime" bson:"RequiredAtDateTime"`
	Verifications        []Verification `json:"Verifications" bson:"Verifications"`
}

type Analogue struct {
	ID              int       `json:"ID" bson:"ID"`
	CreatedByID     int       `json:"CreatedBy" bson:"CreatedBy"`
	CreatedByName   string    `json:"CreatedByName" bson:"CreatedByName"`
	CreatedDateTime time.Time `json:"CreatedDateTime" bson:"CreatedDateTime"`
	RequestItemID   int       `json:"RequestItemID" bson:"RequestItemID"`
	ARCID           int       `json:"ARCID" bson:"ARCID"`
	ARCWorkItemID   int       `json:"ARCWorkItemID" bson:"ARCWorkItemID"`
	ARCWorkID       int       `json:"ARCWorkID" bson:"ARCWorkID"`
	ARCWorkName     string    `json:"ARCWorkName" bson:"ARCWorkName"`
	Measure         string    `json:"Measure" bson:"Measure"`
	Name            string    `json:"Name" bson:"Name"`
	Qty             float64   `json:"Qty" bson:"Qty"`
	QtyOrdered      float64   `json:"QtyOrdered" bson:"QtyOrdered"`
	QtyDelivered    float64   `json:"QtyDelivered" bson:"QtyDelivered"`
	QtyPayed        float64   `json:"QtyPayed" bson:"QtyPayed"`
}

type Verification struct {
	ID              int       `json:"ID" bson:"ID"`
	InvoiceID       int       `json:"InvoiceID" bson:"InvoiceID"`
	RequestItemID   int       `json:"RequestItemID" bson:"RequestItemID"`
	RequestID       int       `json:"RequestID" bson:"RequestID"`
	ProjectID       int       `json:"ProjectID" bson:"ProjectID"`
	ProjectName     string    `json:"ProjectName" bson:"ProjectName"`
	ObjectID        int       `json:"ObjectID" bson:"ObjectID"`
	ObjectName      string    `json:"ObjectName" bson:"ObjectName"`
	CreatedDateTime time.Time `json:"CreatedDateTime" bson:"CreatedDateTime"`
	DivisionName    string    `json:"divisionName" bson:"divisionName"`
	VerifiedByID    int       `json:"VerifiedByID" bson:"VerifiedByID"`
	VerifiedByName  string    `json:"VerifiedByName" bson:"VerifiedByName"`
	Status          int       `json:"status" bson:"status"` // 0 - поступило 1 - в работе 2 - Отклонено 3 - Согласовано
}

type Request struct {
	ID                   int            `json:"ID" bson:"ID"`
	Type                 int            `json:"Type" bson:"Type"` // 0 - ? 1 - материалы, 2 - механизмы, 3 - трудозатраты, 4 - накладные расходы
	CreatedByID          int            `json:"CreatedBy" bson:"CreatedBy"`
	CreatedByName        string         `json:"CreatedByName" bson:"CreatedByName"`
	CreatedDateTime      time.Time      `json:"CreatedDateTime" bson:"CreatedDateTime"`
	LastModifiedByID     int            `json:"LastModifiedBy" bson:"LastModifiedBy"`
	LastModifiedByName   string         `json:"LastModifiedByName" bson:"LastModifiedByName"`
	LastModifiedDateTime time.Time      `json:"LastModifiedDateTime" bson:"LastModifiedDateTime"`
	Status               int            `json:"Status" bson:"Status"`
	ProjectID            int            `json:"ProjectID" bson:"ProjectID"`
	ProjectName          string         `json:"ProjectName" bson:"ProjectName"`
	ObjectID             int            `json:"ObjectID" bson:"ObjectID"`
	ObjectName           string         `json:"ObjectName" bson:"ObjectName"`
	ARCID                int            `json:"ARCID" bson:"ARCID"`
	ARC                  ARC            `json:"ARC" bson:"ARC"`
	Verifications        []Verification `json:"Verifications" bson:"Verifications"`
	RequestItems         []RequestItem  `json:"Completions" bson:"Completions"`
	WorkID               int            `json:"WorkID" bson:"WorkID"`
	WorkName             string         `json:"WorkName" bson:"WorkName"`
	ChatID               int            `json:"ChatID" bson:"ChatID"`
}

type Object struct {
	ID             int       `json:"ID" bson:"ID"`
	ProjectID      int       `json:"ProjectID" bson:"ProjectID"`
	ProjectName    string    `json:"ProjectName" bson:"ProjectName"`
	Address        string    `json:"Address" bson:"Address"`
	CompletionDate time.Time `json:"CompletionDate" bson:"CompletionDate"`
	Name           string    `json:"Name" bson:"Name"`
	ARCID          int       `json:"AVC_ID" bson:"AVC_ID"`
	ARC            ARC       `json:"ARC" bson:"ARC"`
}

type Project struct {
	ID             int       `json:"ID" bson:"ID"`
	Name           string    `json:"Name" bson:"Name"`
	Address        string    `json:"Address" bson:"Address"`
	CompletionDate time.Time `json:"CompletionDate" bson:"CompletionDate"`
	Objects        []Object  `json:"Objects" bson:"Objects"`
	Price          float64   `json:"Price" bson:"Price"`
}
