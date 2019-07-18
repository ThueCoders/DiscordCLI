package logger

import (
    "os"
    "log"
)

var (
    Log *log.Logger
)

func Init(f *os.File) {
	Log = log.New(f, "", log.LstdFlags|log.Lshortfile)
}
