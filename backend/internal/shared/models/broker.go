package models

import "time"

type Broker struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code       string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Permission string    `gorm:"type:text" json:"permission"`
	Group      string    `gorm:"type:varchar(100);index" json:"group"`
	Color      string    `gorm:"type:varchar(50)" json:"color"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Broker) TableName() string {
	return "brokers"
}
