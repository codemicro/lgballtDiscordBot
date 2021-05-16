package db

import "time"

type AnalyticsEvent struct {
	Time      time.Time `gorm:"primarykey"`
	EventType string
	Event     string
}

func (a *AnalyticsEvent) Create() error {
	return Conn.Create(a).Error
}