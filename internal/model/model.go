package model

import "time"

// CreateModel 内嵌 model
type CreateModel struct {
	Creator    string    `gorm:"type:varchar(100);not null;default ''"`
	CreateTime time.Time `gorm:"autoCreateTime"` // 在创建时自动生成时间
}

// ModifyModel 内嵌 model
type ModifyModel struct {
	Modifier   string    `gorm:"type:varchar(100);not null;default ''"`
	modifyTime time.Time `gorm:"autoUpdateTime"` // 在更新记录时自动生成时间
}

// User 用户
type User struct {
	CreateModel
	ModifyModel
	ID       int    `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	Gender   string `gorm:"column:gender"`
	Age      int    `gorm:"column:age"`
	PassWord string `gorm:"column:password"`
	NickName string `gorm:"column:nickname"`
}

func (t *User) TableName() string {
	return "t_user"
}
