package db

import (
	"errors"
	"gorm.io/gorm"
	"strings"
)

type UserMute struct {
	UserId          string `gorm:"primarykey"`
	GuildId         string
	EndTime         int64
	RawRemovedRoles string
	RemovedRoles    []string `gorm:"-"`
}

func (um *UserMute) Get() (bool, error) {
	err := Conn.Where(um).First(um).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	um.RemovedRoles = strings.Split(um.RawRemovedRoles, ",")
	return true, nil
}

func (um *UserMute) Save() error {
	um.RawRemovedRoles = strings.Join(um.RemovedRoles, ",")
	return Conn.Save(um).Error
}

func (um *UserMute) Create() error {
	um.RawRemovedRoles = strings.Join(um.RemovedRoles, ",")
	return Conn.Create(um).Error
}

func (um *UserMute) Delete() error {
	return Conn.Where(um).Delete(um).Error
}

func GetAllUserMutes() ([]UserMute, error) {
	var ums []UserMute
	err := Conn.Find(&ums).Error
	for i := range ums {
		ums[i].RemovedRoles = strings.Split(ums[i].RawRemovedRoles, ",")
	}
	return ums, err
}
