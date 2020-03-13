package modles

import "time"

type Article struct {
	ID         int       `json:"id" gorm:"AUTO_INCREMENT;NOT NULL"`
	TagId      int       `json:"tag_id" gorm:"index"`
	Tag        string    `json:"tag" gorm:"NOT NULL"`
	Title      string    `json:"title" gorm:"NOT NULL"`
	Desc       string    `json:"desc" gorm:"NOT NULL"`
	Content    string    `json:"content" gorm:"NOT NULL"`
	CreateTime time.Time `json:"create_time" gorm:"NOT NULL"`
	Status     int       `json:"status" gorm:"NOT NULL"`
}
