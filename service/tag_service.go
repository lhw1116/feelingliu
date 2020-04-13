package service

import "feelingliu/modles"

type Tag struct {
	ID      int    `json:"id" db:"id"`
	TagName string `json:"tag_name" db:"tag_name" binding:"required,max=16"`
}

func (t Tag) GetAll() ([]Tag, error) {
	tags := make([]Tag, 0)
	db := modles.DB.Find(&tags)
	return tags, db.Error
}