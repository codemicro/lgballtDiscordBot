package logging

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const logFileName = "log.log"

func toFile(e string) {

	var fileExists bool
	{
		info, err := os.Stat(logFileName)
		if os.IsNotExist(err) {
			fileExists = false
		} else {
			fileExists = !info.IsDir()
		}
	}

	if !fileExists {
		_, _ = os.Create(logFileName)
	}

	cTime := strconv.FormatInt(time.Now().Unix(), 10)
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	if _, err := file.WriteString(fmt.Sprintf("%s %s\n", cTime, e)); err != nil {
		panic(err)
	}
	file.Close()
}

func Error(err error, message ...string) {
	if len(message) == 0 {
		ErrorString(err.Error())
	} else {
		ErrorString(fmt.Sprintf("%s\n - %s", message[0], err))
	}
}

func ErrorString(e string) {
	e = "ERROR: " + e
	_, _ = fmt.Fprintln(os.Stderr, e)
	toFile(e)
}

func Warn(e string) {
	e = "WARNING: " + e
	_, _ = fmt.Fprintln(os.Stderr, e)
	toFile(e)
}

func Info(e string) {
	e = "INFO: " + e
	_, _ = fmt.Fprintln(os.Stderr, e)
	toFile(e)
}
