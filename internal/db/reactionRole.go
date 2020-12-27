package db

import (
	"errors"
	"gorm.io/gorm"
)

type ReactionRole struct {
	MessageId string
	RoleId    string
	Emoji     string
}

func (r *ReactionRole) Get() (bool, error) {
	err := Conn.Model(&ReactionRole{}).Take(r).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (r *ReactionRole) Save() error {
	// I don't think this function will ever really get much use, but it's here anyway.
	// It also doesn't work because it needs where clauses
	return Conn.Save(r).Error
}

func (r *ReactionRole) Create() error {
	return Conn.Create(r).Error
}

func (r *ReactionRole) Delete() error {
	return Conn.Where(r).Delete(r).Error
}

func GetAllReactionRolesForMessage(messageId string) ([]ReactionRole, error) {

	var all []ReactionRole

	err := Conn.Where(&ReactionRole{MessageId: messageId}).Find(&all).Error
	if err != nil {
		return []ReactionRole{}, err
	}

	return all, nil
}
