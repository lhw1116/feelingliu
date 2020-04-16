package service

import (
	"feelingliu/modles"
	"fmt"
)

type Tag struct {
	ID      int    `json:"id" db:"id"`
	TagName string `json:"tag_name" db:"tag_name" binding:"required,max=16"`
}

func (t Tag) GetAll() ([]Tag, error) {
	tags := make([]Tag, 0)
	db := modles.DB.Find(&tags)
	return tags, db.Error
}

func (t *Tag) Create() (Tag, error) {
	var one Tag = Tag{TagName: t.TagName}
	db := modles.DB.Create(&one)
	if db.Error != nil {
		return Tag{}, db.Error
	}
	return one, nil
}

func (t Tag) GetOne() (Tag, error) {
	var tag Tag
	db := modles.DB.Where("id = ?", t.ID).Find(&tag)
	return tag, db.Error
}

func (t *Tag) Delete() error {
	db := modles.DB.Delete(&Tag{ID: t.ID})
	return db.Error
}

func (t *Tag) Edit() error {
	var tag Tag = Tag{ID: t.ID}
	db := modles.DB.Model(&tag).Update("tag_name", t.TagName)
	fmt.Println("err: ", db.Error)

	return db.Error
}
