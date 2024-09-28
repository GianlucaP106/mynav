package core

import (
	"log"
	"os"
)

var logger *log.Logger

func Init(path string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panicln(err)
	}
	logger = log.New(f, "", 0)
}

func Logger() *log.Logger {
	return logger
}
