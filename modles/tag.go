package modles

import "time"

type Tag struct {
	Id         int       `json:"id" gorm:"AUTO_INCREMENT;NOT NULL"`
	TagName    string    `json:"tag_name" gorm:"NOT NULL"`
	Status     int       `json:"status" gorm:"NOT NULL"`
	CreateTime time.Time `json:"create_time" gorm:"NOT NULL"`
}
