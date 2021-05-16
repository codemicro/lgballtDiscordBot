package analytics

import (
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"time"
)

func newAnalyticsEvent(eventType, event string) *db.AnalyticsEvent {
	return &db.AnalyticsEvent{
		Time:      time.Now(),
		EventType: eventType,
		Event:     event,
	}
}

const (
	commandUseEventType = "commandUse"
	pluralkitRequestEventType = "pkRequest"
)

func ReportCommandUse(commandName string) {
	err := newAnalyticsEvent(commandUseEventType, commandName).Create()
	if err != nil {
		logging.Warn("Analytics event create: " + err.Error())
	}
}

func ReportPluralKitRequest(requestName string) {
	err := newAnalyticsEvent(pluralkitRequestEventType, requestName).Create()
	if err != nil {
		logging.Warn("Analytics event create: " + err.Error())
	}
}