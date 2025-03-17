package dm_mysql

import "time"

type Message struct {
	ID          int64     `gorm:"primaryKey" json:"ID"`
	CT          time.Time `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	Content     string    `gorm:"type:text;not null" json:"Content"`
	AuthorID    int64     `gorm:"not null" json:"AuthorID"`
	ChatID      int64     `gorm:"not null" json:"ChatID"`
	Author      *User     `gorm:"foreignKey:AuthorID" json:"Author"`
	RepliedToID *int64    `json:"RepliedToID"`
	ForwardedID *int64    `json:"ForwardedID"`
	FileID      *int64    `json:"FileID"`
	File        *DMFile   `gorm:"foreignKey:FileID" json:"File"`
}

type Chat struct {
	ID           int64     `gorm:"primaryKey" json:"ID"`
	CT           time.Time `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	Name         string    `gorm:"type:NVARCHAR(255);not null" json:"Name"`
	Description  string    `gorm:"type:text" json:"Description"`
	Messages     []Message `gorm:"foreignKey:ChatID" json:"Messages"`
	Participants []User    `gorm:"many2many:user_chat;" json:"Participants"`
	OwnerID      int64     `gorm:"not null" json:"OwnerID"`
	Owner        *User     `gorm:"foreignKey:OwnerID" json:"Owner"`
}

type ChatGroup struct {
	ID      int64     `gorm:"primaryKey" json:"ID"`
	CT      time.Time `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	Name    string    `gorm:"type:NVARCHAR(255);not null" json:"Name"`
	Chats   []Chat    `gorm:"many2many:chat_group_chats;" json:"Chats"`
	OwnerID int64     `gorm:"not null" json:"OwnerID"`
	Owner   *User     `gorm:"foreignKey:OwnerID" json:"Owner"`
}

type AccessRights struct {
	ID     int64 `gorm:"primaryKey" json:"ID"`
	Create bool  `json:"Create"`
	Read   bool  `json:"Read"`
	Update bool  `json:"Update"`
	Delete bool  `json:"Delete"`
	Verify bool  `json:"Verify"`
	Owner  bool  `json:"Owner"`
}

type UserChatAccess struct {
	ID             int64         `gorm:"primaryKey" json:"ID"`
	UserID         int64         `json:"-"`
	User           *User         `gorm:"foreignKey:UserID" json:"User"`
	ChatID         int64         `json:"-"`
	Chat           *Chat         `gorm:"foreignKey:ChatID" json:"Chat"`
	AccessRightsID int64         `json:"AccessRightsID"`
	AccessRights   *AccessRights `gorm:"foreignKey:AccessRightsID" json:"AccessRights"`
}
