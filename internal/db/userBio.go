package db

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
)

type UserBio struct {
	UserId     string `gorm:"primarykey"`
	RawBioData string
	BioData    map[string]string `gorm:"-"`
}

func marshalBioData(raw map[string]string) (string, error) {
	b, err := json.Marshal(raw)
	return string(b), err
}

func (bio *UserBio) Populate(userId string) (found bool, err error) {
	if found, err = bio.PopulateRaw(userId); err != nil {
		return false, err
	} else if !found {
		return false, nil
	}
	return true, json.Unmarshal([]byte(bio.RawBioData), &bio.BioData)
}

func (bio *UserBio) PopulateRaw(userId string) (found bool, err error) {
	bio.UserId = userId
	err = Conn.Take(&bio).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return found, nil
		} else {
			return
		}
	}
	found = true
	return
}

func (bio *UserBio) Save() error {
	rb, err := marshalBioData(bio.BioData)
	if err != nil {
		return err
	}
	bio.RawBioData = rb
	return bio.SaveRaw()
}

func (bio *UserBio) SaveRaw() error {
	return Conn.Save(bio).Error
}

func (bio *UserBio) Create() error {
	rb, err := marshalBioData(bio.BioData)
	if err != nil {
		return err
	}
	bio.RawBioData = rb
	return bio.CreateRaw()
}

func (bio *UserBio) CreateRaw() error {
	return Conn.Create(bio).Error
}

func (bio *UserBio) Delete() error {
	return Conn.Delete(bio).Error
}
