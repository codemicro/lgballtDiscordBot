package bios

import (
	"encoding/json"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"io/ioutil"
	"sync"
)

const biosFile = "biosData.json"
var biosFileLock sync.RWMutex

type biosData struct {
	Lock *sync.RWMutex                     `json:"-"`
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