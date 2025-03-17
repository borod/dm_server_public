package dm_actions

import (
	mysql "dm_server/dm_db/dm_mysql"
	"time"
)

type ResidueResponse struct {
	Residue float64 `json:"Residue"`
	Total   float64 `json:"Total"`
}

type UserDivisionsResponse struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

type EnumResponse struct {
	Statuses []string `json:"statuses"`
}

type IncExclChatParticipantsResponse struct {
	ID     int64 `json:"ID"`
	Status int   `json:"Status"` //0 undefined 1 added, 2 Already (included/excluded), 3 does not exist
}

type CreateChatPayload struct {
	Name          string  `json:"Name"`
	Participants  []int64 `json:"Participants"`
	Justification string  `json:"Justification"`
}

type VerifyEntityPayload struct {
	InvoiceID     int64  `json:"InvoiceID"`
	RequestItemID int64  `json:"RequestItemID"`
	Justification string `json:"Justification"`
	Status        int    `json:"Status"` //0 - поступило 1 - в работе 2 - отклонено 3 - согласовано
}

type UploadMediaPayload struct {
	Name string `json:"Name"`
}

type ChatParticipantsResponse struct {
	ID       int64  `json:"ID"`
	Email    string `json:"Email"`
	Name     string `json:"Name"`
	Nickname string `json:"Nickname"`
}

type MessageResponse struct {
	ID          int64     `json:"ID"`
	CT          time.Time `json:"CT"`
	Content     string    `json:"Content"`
	AuthorID    int64     `json:"AuthorID"`
	ChatID      int64     `json:"ChatID"`
	Author      string    `json:"Author"`
	RepliedToID int64     `json:"RepliedToID"`
	ForwardedID int64     `json:"ForwardedID"`
	FileID      int64     `json:"FileID"`
}

type InitVerificationsPayload struct {
	InvoiceID     int64
	RequestItemID int64
}

type PageCreateRequestResponse struct {
	List []CRCard
}

type PageWorkMobileResponse struct {
	Projects     []PwmProject
	Objects      []PwmObject
	Requests     []PwmRequest
	RequestItems []PwmRequestItem
	Invoices     []*mysql.Invoice
	InvoiceItems []*mysql.InvoiceItem
}

type PageWorListZayavkaResponse struct {
	Requests []PwmRequest
}

type VerifierResponse struct {
	ID                 int64
	Name               string
	DivisionID         int64
	DivisionName       string
	VerificationStatus int
	Date               time.Time
}

type RequestedItemResponse struct {
	Name                string
	Qty                 float64
	IsAnalogue          bool
	AnalogueInitialName string
	RequiredAtDateTime  time.Time
	Measure             string
}

type PageViewingZakazResponse struct {
	CreatedDateTime time.Time
	CreatedByName   string
	ProjectName     string
	ObjectName      string
	WorkName        string
	ChatID          int64
	Verifiers       []VerifierResponse
	UserStatus      int
	Excess          bool
	ARCTotal        float64
	ARCResidue      float64
	RequestedItem   RequestedItemResponse
	Status          int
	Percentages     InfoPercentages
	Note            string
	CreatedByID     int64
	ObjectID        int64
	ARCWorkID       int64
	ARCID           int64
}

type CRCard struct {
	ObjectID    int64
	ObjectName  string
	ARCWorkName string
	ARCCT       time.Time
	Qty         float64
	Measure     string
	Address     string
	ID          int64
	CODE        string
}

type PwmProject struct {
	ID             int64
	Name           string
	Address        string
	CompletionDate time.Time
	Price          float64
}

type PwmObject struct {
	ID             int64
	Name           string
	ProjectName    string
	Address        string
	CompletionDate time.Time
}

type PwmRequest struct {
	ID           int64
	WorkName     string
	Status       int
	CT           time.Time
	ProjectName  string
	ObjectID     int64
	ObjectName   string
	ChatID       int64
	RequestItems []PwmRequestItem
}

type PwmRequestItem struct {
	ID                 int64
	ARCWorkItemID      int64
	ARCWorkItemName    string
	Excess             bool
	Name               string
	Qty                float64
	Measure            string
	Note               string
	RequiredAtDateTime time.Time
	ProjectName        string
	ObjectID           int64
	ObjectName         string
	Status             int
	Percentages        InfoPercentages
}

type InfoPercentages struct {
	QtyOrdered   float64
	QtyDelivered float64
	QtyPayed     float64
}

type LzRequestItems struct {
	Name    string
	Qty     float64
	Measure string
	Exceed  bool
}

type InvoiceItemResponse struct {
	Name               string
	RequiredAtDateTime time.Time
	Measure            string
	TotalPrice         float64
	QtyInvoice         float64
	QtyRequest         float64
	InvoiceItemID      int64
	RequestItemID      int64
	DiscountType       string
	ARCID              int64
	UPWithoutNDS       float64
	UPWithNDS          float64
	UPWithDiscount     float64
	DiscountValue      float64
}

// FP - Final Price
type PageViewingInvoiceResponse struct {
	CT             time.Time
	CreatedByName  string
	ProjectName    string
	ObjectID       int64
	ObjectName     string
	ChatID         int64
	Date           time.Time
	Verifiers      []VerifierResponse
	Items          []InvoiceItemResponse
	Verifications  []*mysql.Verification
	Status         int     // 0 — отменён, 1 — черновик,  2 — создан, 3 — в работе, 4 — закупка, 5 — оплата, 6 — завершён
	UserStatus     int     // 0 - поступил 1 - в работе 2 - Отклонён 3 - Согласован
	FPWithoutNDS   float64 `json:"FPWithoutNDS" bson:"FPWithoutNDS"`
	FPWithNDS      float64 `json:"FPWithNDS" bson:"FPWithNDS"`
	FPNDS          float64 `json:"FPNDS" bson:"FPNDS"`
	FPWithDiscount float64 `json:"FPWithDiscount" bson:"FPWithDiscount"`
	Payments       []*mysql.Payment
	PayerID        int64
	PayerName      string
	SupplierID     int64
	SupplierName   string
	CustomerID     int64
	CustomerName   string
	InvoiceNumber  string
	TotalWeight    float64
	MeasureWeight  string
	CreatedByID    int64
}

type RequestedItemViewZayavkaResponse struct {
	ID                  int64
	Name                string
	Qty                 float64
	Measure             string
	RequiredAtDateTime  time.Time
	AnalogueInitialName string
	Status              int
	Excess              bool
	Percentages         InfoPercentages
	Note                string
	Residue             ResidueResponse
	AnalogueID          int64
	ARCWorkItemID       int64
	ARCWorkItemName     string
}

type PageViewingZayavkaResponse struct {
	CT           time.Time
	WorkName     string
	ProjectName  string
	ObjectID     int64
	ObjectName   string
	RequestItems []RequestedItemViewZayavkaResponse
	Status       int
	ARCWorkID    int64
	ARCID        int64
}

type PageListZakazZakupkaResponse struct {
	ID                       int64
	RequiredAtDateTime       time.Time
	Name                     string
	Residue                  ResidueResponse
	Measure                  string
	SummMoneyTotalExpended   float64
	SummMoneyResidueToExpend float64
	SummMoneyExpectedViaARC  float64
	SummMoneyTotalPrediction float64
	ARCWorkItemID            int64
	RequestID                int64
	Qty                      float64
	Note                     string
}

type PageListProviderResponse struct {
	ID          int64
	Address     string
	INN         string
	ServiceType string
}

type PageChatDetilsResponse struct {
	Name              string
	CountParticipants int
	CountMessages     int
	Messages          []MessageResponse
}

type PaymentCounteragentDetails struct {
	Name                 string
	INN                  string
	KPP                  string
	BankName             string
	BankBIK              string
	BankAccountNumber    string
	CorrespondentAccount string
}

type ViewingOnePaymentResponse struct {
	ID          int64
	PaymentDate time.Time
	Summ        float64
	Payer       PaymentCounteragentDetails
	Supplier    PaymentCounteragentDetails
}

type PliResult struct {
	ID               int64
	CounteragentID   int64
	CounteragentName string
	InvoiceNumber    string
	InvoiceDate      time.Time
	CT               time.Time
	FPWithNDS        float64
	ObjectID         int64
	ObjectName       string
	Status           int
	ChatID           int64
}

type PageListStaffResponse struct {
	ID            int64
	Name          string
	Email         string
	Nickname      string
	Divisions     []UserDivisionsResponse
	RedmineUserID string
	Phone         string
	CT            time.Time
}

type AddedVerifierResponse struct {
	Verifier PageListStaffResponse
	ChatID   int64
	Date     time.Time
}

type UserChatResponse struct {
	ID            int64
	Name          string
	Description   string
	CountMessages int
}
