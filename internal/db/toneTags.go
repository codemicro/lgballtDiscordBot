package db

import (
	"errors"
	"gorm.io/gorm"
	"sort"
)

type ToneTag struct {
	Shorthand string `gorm:"primarykey"`
	Description string
}

type ToneTagSlice []ToneTag

func (ts ToneTagSlice) Less(i, j int) bool {
	return ts[i].Shorthand < ts[j].Shorthand
}

func (ts ToneTagSlice) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts ToneTagSlice) Len() int {
	return len(ts)
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

func GetAllToneTags() (ToneTagSlice, error) {
	var all ToneTagSlice
	err := Conn.Where(&ToneTag{}).Find(&all).Error
	if err != nil {
		return nil, err
	}
	sort.Sort(all)
	return all, nil
}