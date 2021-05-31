package db

import (
	"errors"
	"gorm.io/gorm"
)

type ToneTag struct {
	Shorthand string `gorm:"primarykey"`
	Description string
}

func (t *ToneTag) Get() (bool, error) {
	err := Conn.Model(&ToneTag{}).Take(t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (t *ToneTag) Save() error {
	return Conn.Where(&ToneTag{Shorthand: t.Shorthand}).Save(t).Error
}

func (t *ToneTag) Create() error {
	return Conn.Create(t).Error
}

func (t *ToneTag) Delete() error {
	return Conn.Where(t).Delete(t).Error
}

func GetAllToneTags() ([]ToneTag, error) {
	var all []ToneTag
	err := Conn.Where(&ToneTag{}).Find(&all).Error
	if err != nil {
		return nil, err
	}
	return all, nil
}