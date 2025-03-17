package dm_mysql

import "time"

type APIKey struct {
	ID      int64  `gorm:"primaryKey" json:"ID" bson:"ID"`
	UserID  int64  `json:"UserID" bson:"UserID"`
	API     int64  `gorm:"primaryKey" json:"API" bson:"API"`
	Key     string `json:"Key" bson:"Key"`
	PartyID string `json:"PartyID" bson:"PartyID"`
}

type User struct {
	ID            int64            `gorm:"primaryKey" json:"ID" bson:"ID"`
	CT            time.Time        `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT" bson:"CT"`
	Email         string           `gorm:"type:NVARCHAR(320);not null;unique" json:"Email" bson:"Email"`
	Password      string           `gorm:"type:text;not null" json:"-" bson:"-"`
	StatusID      int64            `json:"StatusID" bson:"StatusID"`
	Tokens        []Token          `gorm:"foreignKey:UserID" json:"Tokens" bson:"Tokens"`
	APIKeys       []APIKey         `gorm:"foreignKey:UserID" json:"APIKeys" bson:"APIKeys"`
	Codes         []Code           `gorm:"foreignKey:UserID" json:"Codes" bson:"Codes"`
	Name          string           `gorm:"type:text" json:"Name" bson:"Name"`
	Nickname      string           `gorm:"type:text" json:"Nickname" bson:"Nickname"`
	Phone         string           `gorm:"type:text" json:"Phone" bson:"Phone"`
	Projects      []Project        `gorm:"many2many:user_project;" json:"Projects" bson:"Projects"`
	Objects       []Object         `gorm:"many2many:user_object;" json:"Objects" bson:"Objects"`
	Chats         []Chat           `gorm:"many2many:user_chat;" json:"Chats" bson:"Chats"`
	Requests      []Request        `gorm:"many2many:request_user" json:"Requests" bson:"Requests"`
	RequestItems  []RequestItem    `gorm:"many2many:requestitem_user" json:"RequestItems" bson:"RequestItems"`
	Invoices      []Invoice        `gorm:"many2many:invoice_user" json:"Invoices" bson:"Invoices"`
	Verifications []Verification   `gorm:"many2many:verification_user" json:"Verifications" bson:"Verifications"`
	ChatAccess    []UserChatAccess `gorm:"many2many:chataccess_user" json:"UserChatAccess" bson:"UserChatAccess"`
	Divisions     []Division       `gorm:"many2many:division_user;" json:"Divisions" bson:"Divisions"`
}

type Token struct {
	ID            int64     `gorm:"primaryKey" json:"ID" bson:"ID"`
	CT            time.Time `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT" bson:"CT"`
	UserID        int64     `json:"UserID" bson:"UserID"`
	Authorization string    `gorm:"type:text" json:"-" bson:"-"`
	Refresh       string    `gorm:"type:text" json:"Refresh" bson:"Refresh"`
}

type Code struct {
	ID     int64     `gorm:"primaryKey" json:"ID" bson:"ID"`
	CT     time.Time `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT" bson:"CT"`
	UserID int64     `json:"UserID" bson:"UserID"`
	Code   string    `gorm:"type:text" json:"Code" bson:"Code"`
}

type Division struct {
	ID           int64   `gorm:"primaryKey" json:"ID" bson:"ID"`
	Name         string  `gorm:"type:text" json:"Name" bson:"Name"`
	ManagerID    int64   `gorm:"not null" json:"UserID" bson:"UserID"`
	Manager      *User   `gorm:"foreignKey:ManagerID" json:"Manager" bson:"-"`
	Participants []*User `gorm:"many2many:division_user;" json:"Participants" bson:"-"`
}
