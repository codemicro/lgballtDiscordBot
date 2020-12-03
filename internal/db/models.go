package db

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
)

type UserBio struct {
	UserId string `gorm:"primarykey"`
	RawBioData string
	BioData map[string]string `gorm:"-"`
}

func marshalBioData(raw map[string]string) (string, error) {
	b, err := json.Marshal(raw)
	return string(b), err
}

func (bio *UserBio) Populate(userId string) (found bool, err error) {

	bio.UserId = userId
	conn := Conn

	err = conn.Take(&bio).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return found, nil
		} else {
			return
		}
	}

	found = true

	err = json.Unmarshal([]byte(bio.RawBioData), &bio.BioData)

	return
}

func (bio *UserBio) Save() error {
	rb, err := marshalBioData(bio.BioData)
	if err != nil {
		return err
	}
	bio.RawBioData = rb
	return Conn.Save(bio).Error
}

func (bio *UserBio) Create() error {
	rb, err := marshalBioData(bio.BioData)
	if err != nil {
		return err
	}
	bio.RawBioData = rb
	return Conn.Create(bio).Error
}

func (bio *UserBio) Delete() error {
	return Conn.Delete(bio).Error
}