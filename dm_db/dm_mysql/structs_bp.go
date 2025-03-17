package dm_mysql

import (
	"time"
)

// Платёжное поручение
type PaymentOrder struct {
	ID          int64      `json:"ID" bson:"ID"`
	CT          time.Time  `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID int64      `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy   *User      `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	Date        *time.Time `json:"Date" bson:"Date"`
}

// Платёж
type Payment struct {
	ID             int64         `json:"ID" bson:"ID"`
	CT             time.Time     `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	InvoiceID      int64         `gorm:"not null" json:"InvoiceID" bson:"InvoiceID"`
	Invoice        *Invoice      `gorm:"foreignKey:InvoiceID" json:"Invoice" bson:"Invoice"`
	CreatedByID    int64         `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy      *User         `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	Date           *time.Time    `json:"Date" bson:"Date"`
	Summ           float64       `gorm:"not null" json:"Summ" bson:"Summ"`
	PayerID        int64         `gorm:"not null" json:"PayerID" bson:"PayerID"`
	Payer          *Counteragent `gorm:"foreignKey:PayerID" json:"Payer" bson:"Payer"`
	AcceptorID     *int64        `gorm:"not null" json:"AcceptorID" bson:"AcceptorID"`
	Acceptor       *Counteragent `gorm:"foreignKey:AcceptorID" json:"Acceptor" bson:"Acceptor"`
	PaymentOrderID *int64        `json:"PaymentOrderID" bson:"PaymentOrderID"`
	PaymentOrder   *PaymentOrder `gorm:"foreignKey:PaymentOrderID" json:"PaymentOrder" bson:"PaymentOrder"`
}

// Счёт
// FP - Final Price
type Invoice struct {
	ID                   int64           `json:"ID" bson:"ID"`
	CT                   time.Time       `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID          int64           `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy            *User           `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	LastModifiedByID     int64           `gorm:"not null" json:"LastModifiedByID" bson:"LastModifiedByID"`
	LastModifiedBy       *User           `gorm:"foreignKey:LastModifiedByID" json:"LastModifiedBy" bson:"LastModifiedBy"`
	LastModifiedDateTime *time.Time      `gorm:"type:datetime(6);default:current_timestamp(6);" json:"LastModifiedDateTime" bson:"LastModifiedDateTime"`
	NumberStr            string          `json:"NumberStr" bson:"NumberStr"`
	Date                 *time.Time      `gorm:"type:datetime(6);default:current_timestamp(6);" json:"Date" bson:"Date"`
	SupplierID           int64           `json:"SupplierID" bson:"SupplierID"`
	Supplier             *Counteragent   `gorm:"foreignKey:SupplierID" json:"Supplier" bson:"Supplier"`
	CustomerID           int64           `json:"CustomerID" bson:"CustomerID"`
	Customer             *Counteragent   `gorm:"foreignKey:CustomerID" json:"Customer" bson:"Customer"`
	PayerID              int64           `json:"PayerID" bson:"PayerID"`
	Payer                *Counteragent   `gorm:"foreignKey:PayerID" json:"Payer" bson:"Payer"`
	Status               int             `json:"Status" bson:"Status"`
	ProjectID            int64           `json:"ProjectID" bson:"ProjectID"`
	Project              *Project        `gorm:"foreignKey:ProjectID" json:"Project" bson:"ProjectID"`
	ObjectID             int64           `json:"ObjectID" bson:"ObjectID"`
	Object               *Object         `gorm:"foreignKey:ObjectID" json:"Object" bson:"Object"`
	FPWithoutNDS         *float64        `json:"FPWithoutNDS" bson:"FPWithoutNDS"`
	FPWithNDS            *float64        `json:"FPWithNDS" bson:"FPWithNDS"`
	FPNDS                *float64        `json:"FPNDS" bson:"FPNDS"`
	FPWithDiscount       *float64        `json:"FPWithDiscount" bson:"FPWithDiscount"`
	TotalWeight          *float64        `json:"TotalWeight" bson:"TotalWeight"`
	MeasureWeight        string          `json:"MeasureWeight" bson:"MeasureWeight"`
	Verifications        []*Verification `gorm:"foreignKey:ObjID" json:"Verifications" bson:"Verifications"`
	InvoiceItems         []*InvoiceItem  `gorm:"foreignKey:InvoiceID" json:"InvoiceItems" bson:"InvoiceItems"`
	Payments             []*Payment      `gorm:"foreignKey:InvoiceID" json:"Payments" bson:"Payments"`
	ChatID               *int64          `json:"ChatID" bson:"ChatID"`
	Chat                 *Chat           `gorm:"foreignKey:ChatID" json:"Chat" bson:"Chat"`
	Requests             []Request       `gorm:"many2many:invoice_request" json:"Requests" bson:"Requests"`
	RequestItems         []RequestItem   `gorm:"many2many:invoice_requestitem" json:"RequestItems" bson:"RequestItems"`
}

// строка в счёте
// UP Unit Price
// TC Total Cost
type InvoiceItem struct {
	ID                 int64         `json:"ID" bson:"ID"`
	CT                 time.Time     `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID        int64         `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy          *User         `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	InvoiceID          int64         `json:"InvoiceID" bson:"InvoiceID"`
	Invoice            *Invoice      `gorm:"foreignKey:InvoiceID" json:"Invoice" bson:"Invoice"`
	ProjectID          int64         `json:"ProjectID" bson:"ProjectID"`
	Project            *Project      `gorm:"foreignKey:ProjectID" json:"Project" bson:"Project"`
	ObjectID           int64         `json:"ObjectID" bson:"ObjectID"`
	Object             *Object       `gorm:"foreignKey:ObjectID" json:"Object" bson:"Object"`
	ARCID              *int64        `json:"ARCID" bson:"ARCID"`
	ARC                *ARC          `gorm:"foreignKey:ARCID" json:"ARC" bson:"ARC"`
	ARCWorkItemID      *int64        `json:"ARCWorkItemID" bson:"ARCWorkItemID"`
	ARCWorkItem        *ARCWorkItem  `gorm:"foreignKey:ARCWorkItemID" json:"ARCWorkItem" bson:"ARCWorkItem"`
	RequestID          int64         `json:"RequestID" bson:"RequestID"`
	Request            *Request      `gorm:"foreignKey:RequestID" json:"Request" bson:"Request"`
	RequestItemID      int64         `json:"RequestItemID" bson:"RequestItemID"`
	RequestItem        *RequestItem  `gorm:"foreignKey:RequestItemID" json:"RequestItem" bson:"RequestItem"`
	Name               string        `json:"Name" bson:"Name"`
	Measure            string        `json:"Measure" bson:"Measure"`
	Qty                float64       `json:"Qty" bson:"Qty"`
	UPWithoutNDS       float64       `json:"UPWithoutNDS" bson:"UPWithoutNDS"`
	UPWithNDS          float64       `json:"UPWithNDS" bson:"UPWithNDS"`
	UPWithoutDiscount  float64       `json:"UPWithoutDiscount" bson:"UPWithoutDiscount"`
	UPWithDiscount     float64       `json:"UPWithDiscount" bson:"UPWithDiscount"`
	TCWithoutNDS       float64       `json:"TCWithoutNDS" bson:"TCWithoutNDS"`
	TCWithNDS          float64       `json:"TCWithNDS" bson:"TCWithNDS"`
	TCWithoutDiscount  float64       `json:"TCWithoutDiscount" bson:"TCWithoutDiscount"`
	TCWithDiscount     float64       `json:"TCWithDiscount" bson:"TCWithDiscount"`
	RequiredAtDateTime *time.Time    `json:"RequiredAtDateTime" bson:"RequiredAtDateTime"`
	Requests           []Request     `gorm:"many2many:invoice_request" json:"Requests" bson:"Requests"`
	RequestItems       []RequestItem `gorm:"many2many:invoice_requestitem" json:"RequestItems" bson:"RequestItems"`
	DiscountType       string        `json:"DiscountType" bson:"DiscountType"`
	DiscountValue      *float64      `json:"DiscountValue" bson:"DiscountValue"`
}

// Заказ
type RequestItem struct {
	// LastModifiedByID     string         `json:"LastModifiedByID" bson:"LastModifiedByID"`
	// LastModifiedByName   string         `json:"LastModifiedByName" bson:"LastModifiedByName"`
	// LastModifiedDateTime time.Time      `json:"LastModifiedDateTime" bson:"LastModifiedDateTime"`
	// QtyOrdered         float64        `json:"QtyOrdered" bson:"QtyOrdered"`
	// QtyDelivered       float64        `json:"QtyDelivered" bson:"QtyDelivered"`
	// QtyPayed           float64        `json:"QtyPayed" bson:"QtyPayed"`
	ID                 int64           `json:"ID" bson:"ID"`
	CT                 time.Time       `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID        int64           `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy          *User           `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	Type               int             `json:"Type" bson:"Type"`     // 0 - ? 1 - материалы, 2 - механизмы, 3 - трудозатраты, 4 - накладные расходы
	Status             int             `json:"Status" bson:"Status"` //0 — отменён, 1 — черновик,  2 — создан, 3 — в работе, 4 — закупка, 5 — оплата, 6 — завершён
	RequestID          int64           `json:"RequestID" bson:"RequestID"`
	Request            *Request        `gorm:"foreignKey:RequestID" json:"Request" bson:"Request"`
	ProjectID          int64           `json:"ProjectID" bson:"ProjectID"`
	Project            *Project        `gorm:"foreignKey:ProjectID" json:"Project" bson:"Project"`
	ObjectID           int64           `json:"ObjectID" bson:"ObjectID"`
	Object             *Object         `gorm:"foreignKey:ObjectID" json:"Object" bson:"Object"`
	ARCID              *int64          `json:"ARCID" bson:"ARCID"`
	ARC                *ARC            `gorm:"foreignKey:ARCID" json:"ARC" bson:"ARC"`
	ARCWorkItemID      *int64          `json:"ARCWorkItemID" bson:"ARCWorkItemID"`
	ARCWorkItem        *ARCWorkItem    `gorm:"foreignKey:ARCWorkItemID" json:"ARCWorkItem" bson:"ARCWorkItem"`
	ARCWorkID          *int64          `json:"ARCWorkID" bson:"ARCWorkID"`
	ARCWork            *ARCWork        `gorm:"foreignKey:ARCWorkID" json:"ARCWork" bson:"ARCWork"`
	Measure            string          `json:"Measure" bson:"Measure"`
	Name               string          `json:"Name" bson:"Name"`
	Qty                float64         `json:"Qty" bson:"Qty"`
	IsAnalogueSet      bool            `json:"IsAnalogueSet" bson:"IsAnalogueSet"`
	AnalogueID         *int64          `json:"AnalogueID" bson:"AnalogueID"`
	Analogue           *Analogue       `gorm:"foreignKey:AnalogueID" json:"Analogue" bson:"Analogue"`
	Note               string          `json:"Note" bson:"Note"`
	RequiredAtDateTime *time.Time      `json:"RequiredAtDateTime" bson:"RequiredAtDateTime"`
	Verifications      []*Verification `gorm:"foreignKey:ObjID" json:"Verifications" bson:"Verifications"`
	Invoices           []*Invoice      `gorm:"many2many:invoice_request" json:"Invoices" bson:"Invoices"`
	InvoiceItems       []*InvoiceItem  `gorm:"many2many:invoiceitem_request" json:"InvoiceItems" bson:"InvoiceItems"`
}

// Аналог
type Analogue struct {
	// QtyOrdered      float64      `json:"QtyOrdered" bson:"QtyOrdered"`
	// QtyDelivered    float64      `json:"QtyDelivered" bson:"QtyDelivered"`
	// QtyPayed        float64      `json:"QtyPayed" bson:"QtyPayed"`
	ID              int64        `json:"ID" bson:"ID"`
	CT              time.Time    `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID     int64        `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy       *User        `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	CreatedDateTime *time.Time   `json:"CreatedDateTime" bson:"CreatedDateTime"`
	RequestID       int64        `json:"RequestID" bson:"RequestID"`
	Request         *Request     `gorm:"foreignKey:RequestID" json:"Request" bson:"Request"`
	RequestItemID   *int64       `json:"RequestItemID" bson:"RequestItemID"`
	RequestItem     *RequestItem `gorm:"foreignKey:RequestItemID" json:"RequestItem" bson:"RequestItem"`
	ProjectID       int64        `json:"ProjectID" bson:"ProjectID"`
	Project         *Project     `gorm:"foreignKey:ProjectID" json:"Project" bson:"Project"`
	ObjectID        int64        `json:"ObjectID" bson:"ObjectID"`
	Object          *Object      `gorm:"foreignKey:ObjectID" json:"Object" bson:"Object"`
	ARCID           *int64       `json:"ARCID" bson:"ARCID"`
	ARC             *ARC         `gorm:"foreignKey:ARCID" json:"ARC" bson:"ARC"`
	ARCWorkItemID   *int64       `json:"ARCWorkItemID" bson:"ARCWorkItemID"`
	ARCWorkItem     *ARCWorkItem `gorm:"foreignKey:ARCWorkItemID" json:"ARCWorkItem" bson:"ARCWorkItem"`
	ARCWorkID       int64        `json:"ARCWorkID" bson:"ARCWorkID"`
	ARCWork         *ARCWork     `gorm:"foreignKey:ARCWorkID" json:"ARCWork" bson:"ARCWork"`
	Measure         string       `json:"Measure" bson:"Measure"`
	Name            string       `json:"Name" bson:"Name"`
	Qty             float64      `json:"Qty" bson:"Qty"`
}

type Verification struct {
	ID           int64     `json:"ID" bson:"ID"`
	CT           time.Time `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID  int64     `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy    *User     `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	DivisionID   int64     `json:"DivisionID" bson:"DivisionID"`
	Division     *Division `gorm:"foreignKey:DivisionID" json:"Division" bson:"Division"`
	VerifiedByID int64     `json:"VerifiedByID" bson:"VerifiedByID"`
	VerifiedBy   *User     `gorm:"foreignKey:VerifiedByID" json:"VerifiedBy" bson:"VerifiedBy"`
	Status       int       `json:"status" bson:"status"` // 0 - поступило 1 - в работе 2 - Отклонено 3 - Согласовано
	ObjID        int64     `json:"ObjID" bson:"ObjID"`
	ObjType      int       `json:"ObjType" bson:"ObjType"`
}

type Request struct {
	// LastModifiedByID     int            `json:"LastModifiedBy" bson:"LastModifiedBy"`
	// LastModifiedByName   string         `json:"LastModifiedByName" bson:"LastModifiedByName"`
	// LastModifiedDateTime time.Time      `json:"LastModifiedDateTime" bson:"LastModifiedDateTime"`
	ID           int64          `json:"ID" bson:"ID"`
	CT           time.Time      `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID  int64          `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy    *User          `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	Type         int            `json:"Type" bson:"Type"` // 0 - ? 1 - материалы, 2 - механизмы, 3 - трудозатраты
	Status       int            `json:"Status" bson:"Status"`
	ProjectID    int64          `json:"ProjectID" bson:"ProjectID"`
	Project      *Project       `gorm:"foreignKey:ProjectID" json:"Project" bson:"Project"`
	ObjectID     int64          `json:"ObjectID" bson:"ObjectID"`
	Object       *Object        `gorm:"foreignKey:ObjectID" json:"Object" bson:"Object"`
	ARCID        *int64         `json:"ARCID" bson:"ARCID"`
	ARC          *ARC           `gorm:"foreignKey:ARCID" json:"ARC" bson:"ARC"`
	Invoices     []*Invoice     `gorm:"many2many:invoice_request" json:"Invoices" bson:"Invoices"`
	InvoiceItems []*InvoiceItem `gorm:"many2many:invoiceitem_request" json:"InvoiceItems" bson:"InvoiceItems"`
	RequestItems []*RequestItem `gorm:"foreignKey:RequestID" json:"RequestItems" bson:"RequestItems"`
	ARCWorkID    *int64         `json:"ARCWorkID" bson:"WorkID"`
	ARCWork      *ARCWork       `gorm:"foreignKey:ARCWorkID" json:"Work" bson:"Work"`
	ChatID       *int64         `json:"ChatID" bson:"ChatID"`
	Chat         *Chat          `gorm:"foreignKey:ChatID" json:"Chat" bson:"Chat"`
}

type Counteragent struct {
	ID                   int64     `gorm:"primaryKey" json:"ID" bson:"ID"`
	CT                   time.Time `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID          int64     `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy            *User     `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	Name                 string    `gorm:"type:text;not null" json:"Name" bson:"Name"`
	Shortname            string    `gorm:"type:text;not null" json:"Shortname" bson:"Shortname"`
	ServiceType          string    `gorm:"type:text" json:"ServiceType" bson:"ServiceType"`
	LegalAddress         string    `gorm:"type:text" json:"LegalAddress" bson:"LegalAddress"`
	PhysicalAddress      string    `gorm:"type:text" json:"PhysicalAddress" bson:"PhysicalAddress"`
	Iinn                 string    `gorm:"type:text;not null" json:"Iinn" bson:"Iinn"`
	Ogrn                 string    `gorm:"type:text;not null" json:"Ogrn" bson:"Ogrn"`
	Kpp                  string    `gorm:"type:text;not null" json:"Kpp" bson:"Kpp"`
	BankAccountNumber    string    `gorm:"type:text;not null" json:"BankAccountNumber" bson:"BankAccountNumber"`
	CorrespondentAccount string    `gorm:"type:text" json:"CorrespondentAccount" bson:"CorrespondentAccount"`
	BankBIK              string    `gorm:"type:text;not null" json:"BankBIK" bson:"BankBIK"`
	BankName             string    `gorm:"type:text" json:"BankName" bson:"BankName"`
	BankAddress          string    `gorm:"type:text" json:"BankAddress" bson:"BankAddress"`
	PhoneNumber          string    `gorm:"type:text" json:"PhoneNumber" bson:"PhoneNumber"`
	ContactEmail         string    `gorm:"type:text" json:"ContactEmail" bson:"ContactEmail"`
	Website              string    `gorm:"type:text" json:"Website" bson:"Website"`
	CEO                  string    `gorm:"type:text" json:"CEO" bson:"CEO"`
	Acting_basis         string    `gorm:"type:text" json:"Acting_basis" bson:"Acting_basis"`
	Deputy_CEO           string    `gorm:"type:text" json:"Deputy_CEO" bson:"Deputy_CEO"`
	Chief_accountant     string    `gorm:"type:text" json:"Chief_accountant" bson:"Chief_accountant"`
}

type ARCWorkItem struct {
	ID           int64     `gorm:"primaryKey" json:"ID" bson:"ID"`
	CT           time.Time `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID  int64     `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy    *User     `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	Code         string    `gorm:"type:text" json:"Code" bson:"Code"`
	ARCWorkID    int64     `gorm:"not null" json:"ARCWorkID" bson:"ARCWorkID"`
	ARCWork      *ARCWork  `gorm:"foreignKey:ARCWorkID" json:"ARCWork" bson:"ARCWork"`
	ARCID        int64     `json:"ARCID" bson:"ARCID"`
	ARC          *ARC      `gorm:"foreignKey:ARCID" json:"ARC" bson:"ARC"`
	Name         string    `gorm:"type:text" json:"Name" bson:"Name"`
	Measure      string    `gorm:"type:text" json:"Measure" bson:"Measure"`
	Qty          float64   `json:"Qty" bson:"Qty"`
	CPWMU        string    `gorm:"type:text" json:"CPW" bson:"CPW"`
	Note         string    `gorm:"type:text" json:"Note" bson:"Note"`
	PPUWithNDS   float64   `json:"PPUWithNDS" bson:"PPUWithNDS"`
	PriceWithNDS float64   `json:"PriceWithNDS" bson:"PriceWithNDS"`
}

type ARCWork struct {
	ID           int64          `gorm:"primaryKey" json:"ID" bson:"ID"`
	CT           time.Time      `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID  int64          `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy    *User          `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	Code         string         `gorm:"type:text" json:"Code" bson:"Code"`
	Name         string         `gorm:"type:text" json:"Name" bson:"Name"`
	ACRWorkItems []*ARCWorkItem `gorm:"foreignKey:ARCWorkID" json:"ACRWorkItems" bson:"ACRWorkItems"`
	Qty          float64        `json:"Qty" bson:"Qty"`
	Measure      string         `gorm:"type:text" json:"Measure" bson:"Measure"`
	ARCID        int64          `json:"ARCID" bson:"ARCID"`
	ARC          *ARC           `gorm:"foreignKey:ARCID" json:"ARC" bson:"ARC"`
	PPUWithNDS   float64        `json:"PPUWithNDS" bson:"PPUWithNDS"`
	PriceWithNDS float64        `json:"PriceWithNDS" bson:"PriceWithNDS"`
}

type ARC struct {
	ID                   int64      `gorm:"primaryKey" json:"ID" bson:"ID"`
	CT                   time.Time  `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID          int64      `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy            *User      `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	LastModifiedByName   string     `gorm:"type:text" json:"LastModifiedByName" bson:"LastModifiedByName"`
	LastModifiedDateTime *time.Time `gorm:"type:datetime(6);default:current_timestamp(6);" json:"LastModifiedDateTime" bson:"LastModifiedDateTime"`
	ProjectID            int64      `json:"ProjectID" bson:"ProjectID"`
	Project              *Project   `gorm:"foreignKey:ProjectID" json:"Project" bson:"Project"`
	ObjectID             int64      `json:"ObjectID" bson:"ObjectID"`
	Object               *Object    `gorm:"foreignKey:ObjectID" json:"Object" bson:"Object"`
	Name                 string     `gorm:"type:text" json:"Name" bson:"Name"`
	ACRWorks             []*ARCWork `gorm:"foreignKey:ARCID" json:"ACRWorks" bson:"ACRWorks"`
}

type Object struct {
	ID             int64          `gorm:"primaryKey" json:"ID" bson:"ID"`
	CT             time.Time      `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID    int64          `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy      *User          `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	ProjectID      *int64         `json:"ProjectID" bson:"ProjectID"`
	Project        *Project       `gorm:"foreignKey:ProjectID" json:"Project" bson:"Project"`
	ProjectName    string         `json:"ProjectName" bson:"ProjectName"`
	Address        string         `json:"Address" bson:"Address"`
	CompletionDate *time.Time     `json:"CompletionDate" bson:"CompletionDate"`
	Name           string         `json:"Name" bson:"Name"`
	ARCID          int64          `json:"ARCID" bson:"ARCID"`
	ARC            *ARC           `gorm:"foreignKey:ID" json:"ARC" bson:"ARC"`
	Requests       []Request      `gorm:"foreignKey:ObjectID" json:"Requests" bson:"Requests"`
	RequestItems   []RequestItem  `gorm:"foreignKey:ObjectID" json:"RequestItems" bson:"Requests"`
	Invoices       []Invoice      `gorm:"foreignKey:ObjectID" json:"Invoices" bson:"Invoices"`
	InvoiceItems   []InvoiceItem  `gorm:"foreignKey:ObjectID" json:"InvoiceItems" bson:"InvoiceItems"`
	Verifications  []Verification `gorm:"foreignKey:ObjID" json:"Verifications" bson:"Verifications"`
	ChatID         *int           `json:"ChatID" bson:"ChatID"`
	Chat           *Chat          `gorm:"foreignKey:ChatID" json:"Chat" bson:"Chat"`
}

type Project struct {
	ID             int64      `gorm:"primaryKey" json:"ID" bson:"ID"`
	CT             time.Time  `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID    int64      `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy      *User      `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	Name           string     `json:"Name" bson:"Name"`
	Address        string     `json:"Address" bson:"Address"`
	CompletionDate *time.Time `json:"CompletionDate" bson:"CompletionDate"`
	Objects        []*Object  `gorm:"foreignKey:ProjectID" json:"Objects" bson:"Objects"`
	Price          float64    `json:"Price" bson:"Price"`
	ChatID         *int64     `json:"ChatID" bson:"ChatID"`
	Chat           *Chat      `gorm:"foreignKey:ChatID" json:"Chat" bson:"Chat"`
}
