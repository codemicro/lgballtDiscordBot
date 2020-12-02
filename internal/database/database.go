package database

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var Conn *gorm.DB

func init() {
	var err error
	Conn, err = gorm.Open(sqlite.Open(config.Config.DbFileName), &gorm.Config{})
	if err != nil {
		logging.Error(err, fmt.Sprintf("Unable to open database file %s", config.Config.DbFileName))
		os.Exit(1)
	}
}
