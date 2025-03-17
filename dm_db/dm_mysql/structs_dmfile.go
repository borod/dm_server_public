package dm_mysql

import "time"

type DMFile struct {
	ID          int64     `gorm:"primaryKey" json:"ID"`
	CT          time.Time `gorm:"type:datetime(6);default:current_timestamp(6);not null" json:"CT"`
	CreatedByID int64     `gorm:"not null" json:"CreatedByID" bson:"CreatedByID"`
	CreatedBy   *User     `gorm:"foreignKey:CreatedByID" json:"CreatedBy" bson:"CreatedBy"`
	Name        string    `json:"Name"`
}
