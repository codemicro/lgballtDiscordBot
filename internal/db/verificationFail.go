package db

import (
	"errors"
	"gorm.io/gorm"
)

type VerificationFail struct {
	UserId      string `gorm:"primarykey"`
	MessageLink string
}

func (vf *VerificationFail) Get() (bool, error) {
	err := Conn.Where(vf).First(vf).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (vf *VerificationFail) Save() error {
	return Conn.Save(vf).Error
}

func (vf *VerificationFail) Create() error {
	return Conn.Create(vf).Error
}

func (vf *VerificationFail) Delete() error {
	return Conn.Where(vf).Delete(vf).Error
}
