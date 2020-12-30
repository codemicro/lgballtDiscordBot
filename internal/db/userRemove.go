package db

import (
	"errors"
	"gorm.io/gorm"
)

type UserRemove struct {
	UserId string `gorm:"primarykey"`
	Reason    string
	Action string
}

func (r *UserRemove) Get() (bool, error) {
	err := Conn.Where(r).First(r).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (r *UserRemove) Save() error {
	return Conn.Save(r).Error
}

func (r *UserRemove) Create() error {
	return Conn.Create(r).Error
}

func (r *UserRemove) Delete() error {
	return Conn.Where(r).Delete(r).Error
}
