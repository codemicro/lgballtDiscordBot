package bios

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	_ "github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"gorm.io/gorm"
	"io/ioutil"
	"sync"
)

func getUserBioData(userId string) (found bool, bio db.UserBio, err error) {

	bio.UserId = userId
	conn := db.Conn

	err = conn.Take(&bio).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return found, bio, nil
		} else {
			return
		}
	}

	found = true

	err = json.Unmarshal([]byte(bio.RawBioData), &bio.BioData)
	return
}

const biosFile = "biosData.json"

var biosFileLock sync.RWMutex

type biosData struct {
	Lock     *sync.RWMutex                `json:"-"`
	Fields   []string                     `json:"bioFields"`
	UserBios map[string]map[string]string `json:"userBios"`
}

func loadBiosFile() (b biosData, err error) {
	biosFileLock.RLock()
	fCont, err := ioutil.ReadFile(biosFile)
	biosFileLock.RUnlock()

	if err != nil {
		logging.Warn(fmt.Sprintf("Could not open %s - assuming does not exist, creating from scratch", biosFile))
		err = nil
		b.Lock = new(sync.RWMutex)
		b.UserBios = make(map[string]map[string]string)
		b.Fields = make([]string, 0)
		return
	}

	err = json.Unmarshal(fCont, &b)
	b.Lock = new(sync.RWMutex)
	return
}

func saveBiosFile(b biosData) (err error) {

	b.Lock.RLock()

	jCont, err := json.Marshal(&b)
	if err != nil {
		return
	}

	biosFileLock.Lock()
	err = ioutil.WriteFile(biosFile, jCont, 0666)
	biosFileLock.Unlock()

	b.Lock.RUnlock()

	return
}
