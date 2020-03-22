package modles

import "time"

//  用户表的创建
type User struct {
	Id         int       `gorm:"AUTO_INCREMENT;NOT NULL"`
	Username   string    `gorm:"NOT NULL"`
	Password   string    `gorm:"NOT NULL"`
	CreateTime time.Time `gorm:"NOT NULL"`
}

type GetAuth struct {
	Token    string `json:"token"`
}
