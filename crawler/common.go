package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"time"
)

type Settings struct {
	StartTime time.Time
	SleepTime int

	Logger struct {
		Info    *log.Logger
		Warning *log.Logger
		Error   *log.Logger
	}
}

func InitLogger(infoHandle, warningHandle, errorHandle io.Writer) {
	Configuration.Logger.Info = log.New(infoHandle, "INFO:\t", log.Ldate|log.Ltime|log.Lshortfile)
	Configuration.Logger.Warning = log.New(warningHandle, "WARNING:\t", log.Ldate|log.Ltime|log.Lshortfile)
	Configuration.Logger.Error = log.New(errorHandle, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
}

var Configuration Settings

func CalcHash(content string) string {
	hasher := md5.New()
	hasher.Write([]byte(content))
	return hex.EncodeToString(hasher.Sum(nil))
}
