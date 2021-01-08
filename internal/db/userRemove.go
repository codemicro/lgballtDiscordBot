package db

import (
	"errors"
	"gorm.io/gorm"
)

const (
	BannedKeyword = "banned"
	KickedKeyword = "kicked"
)

type UserRemove struct {
	UserId string
	Reason string
	Action string
}

type userBan struct {
	UserId string `gorm:"primarykey"`
	Reason string
}

type userKick struct {
	UserId string `gorm:"primarykey"`
	Reason string
}

func (r *UserRemove) Get() (bool, error) {

	ub := new(userBan)
	uk := new(userKick)
	ub.UserId = r.UserId
	uk.UserId = r.UserId

	var ubFound, ukFound bool

	// Check the user kicks first. Bans take precedence over kicks.

	err := Conn.Where(uk).First(uk).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
	} else {
		ukFound = true
	}

	err = Conn.Where(ub).First(ub).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
	} else {
		ubFound = true
	}

	if ubFound {
		r.Action = BannedKeyword
		r.Reason = ub.Reason
	} else if ukFound {
		r.Action = KickedKeyword
		r.Reason = uk.Reason
	}

	return ubFound || ukFound, nil
}

func (r *UserRemove) Save() error {
	var obj interface{}

	if r.Action == BannedKeyword {
		obj = &userBan{
			UserId: r.UserId,
			Reason: r.Reason,
		}
	} else if r.Action == KickedKeyword {
		obj = &userKick{
			UserId: r.UserId,
			Reason: r.Reason,
		}
	}

	return Conn.Save(obj).Error
}

func (r *UserRemove) Create() error {
	var obj interface{}

	if r.Action == BannedKeyword {
		obj = &userBan{
			UserId: r.UserId,
			Reason: r.Reason,
		}
	} else if r.Action == KickedKeyword {
		obj = &userKick{
			UserId: r.UserId,
			Reason: r.Reason,
		}
	}

	return Conn.Create(obj).Error
}

func (r *UserRemove) Delete() error {

	var obj interface{}

	if r.Action == BannedKeyword {
		obj = &userBan{
			UserId: r.UserId,
			Reason: r.Reason,
		}
	} else if r.Action == KickedKeyword {
		obj = &userKick{
			UserId: r.UserId,
			Reason: r.Reason,
		}
	}

	return Conn.Where(obj).Delete(obj).Error
}
