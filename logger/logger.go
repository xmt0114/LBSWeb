package logger

import (
	seelog "github.com/cihub/seelog"
	"log"
)

var Logger seelog.LoggerInterface

func init() {
	DisableLog()
	loadLoggerConfig()
}

func DisableLog() {
	Logger = seelog.Disabled
}

func UseLogger(newLogger seelog.LoggerInterface) {
	Logger = newLogger
}

func loadLoggerConfig()  {
	logger, err := seelog.LoggerFromConfigAsFile("./conf/logger.xml")
	if err != nil {
		log.Println("logger load failed")
		return
	}
	log.Println("logger load succ")
	UseLogger(logger)
}