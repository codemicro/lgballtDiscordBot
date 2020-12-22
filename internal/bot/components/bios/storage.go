package bios

import (
	"encoding/json"
	"fmt"
	_ "github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/codemicro/lgballtDiscordBot/internal/logging"
	"io/ioutil"
	"sync"
)

const biosFile = "biosData.json"

var biosFileLock sync.RWMutex

type biosData struct {
	Lock     *sync.RWMutex                `json:"-"`
	Fields   []string                     `json:"bioFields"`
}

func loadBiosFile() (b biosData, err error) {
	biosFileLock.RLock()
	fCont, err := ioutil.ReadFile(biosFile)
	biosFileLock.RUnlock()

	if err != nil {
		logging.Warn(fmt.Sprintf("Could not open %s - assuming does not exist, creating from scratch", biosFile))
		err = nil
		b.Lock = new(sync.RWMutex)
		b.Fields = make([]string, 0)
		return
	}

	err = json.Unmarshal(fCont, &b)
	b.Lock = new(sync.RWMutex)
	return
}
