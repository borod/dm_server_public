package dm_mysql

import (
	"time"
)

// Структура для сущности "category"
type Category struct {
	ID   int    `gorm:"primaryKey;autoIncrement;not null"`
	Name string `gorm:"not null"`
	Code string
}

// Структура для сущности "obj"
type Obj struct {
	ID            int          `gorm:"primaryKey;autoIncrement;not null"`
	Ct            time.Time    `gorm:"autoCreateTime"`
	Ut            time.Time    `gorm:"autoCreateTime"`
	CreatorID     *int         `gorm:"not null"`
	Creator       *Obj         `gorm:"foreignKey:CreatorID"`
	ManagerID     *int         `gorm:"not null"`
	Manager       *Obj         `gorm:"foreignKey:ManagerID"`
	CategoryID    int          `gorm:"not null"`
	Category      *Category    `gorm:"foreignKey:CategoryID"`
	PermissionsID int          `gorm:"not null"`
	Permissions   *Permissions `gorm:"foreignKey:PermissionsID"`
	Name          string       `gorm:"not null;default:Category.Name"`
	Meta          string       `json:"Meta" bson:"Meta"`
	Children      []*Obj       `gorm:"many2many:obj_children;joinForeignKey:ParentID;joinReferences:ChildID" json:"-"`
	Parents       []*Obj       `gorm:"many2many:obj_children;joinForeignKey:ChildID;joinReferences:ParentID" json:"-"`
}

// Структура для сущности "measure"
type Measure struct {
	ID           int    `gorm:"primaryKey;autoIncrement;not null"`
	Code         string `gorm:"type:TEXT"`
	FullNameRus  string `gorm:"type:TEXT"`
	ShortNameRus string `gorm:"type:TEXT"`
	ShortNameInt string `gorm:"type:TEXT"`
	CodeNameRus  string `gorm:"type:TEXT"`
	CodeNameInt  string `gorm:"type:TEXT"`
}

// Структура для сущности "asset" (номенклатура в счетах, сметах, ведомостях)
type Asset struct {
	ID        int      `gorm:"primaryKey;autoIncrement;not null"`
	ObjID     int      `gorm:"primaryKey;not null"`
	Obj       *Obj     `gorm:"foreignKey:ObjID"`
	Qty       float64  `gorm:"type:DECIMAL(20,4);not null"`
	MeasureID int      `gorm:"not null"`
	Measure   *Measure `gorm:"foreignKey:MeasureID"`
}

// Структура для сущности "permissions"
type Permissions struct {
	ID     int  `gorm:"primaryKey;autoIncrement;not null"`
	Create bool `gorm:"not null"`
	Read   bool `gorm:"not null"`
	Update bool `gorm:"not null"`
	Delete bool `gorm:"not null"`
	Verify bool `gorm:"not null"`
	Assign bool `gorm:"not null"`
}
