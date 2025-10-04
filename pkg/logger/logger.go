package logger

import (
	"log"
	"os"
)

var (
	infologger  = log.New(os.Stdout, "INFO_: ", log.Ltime|log.Ldate)
	errorlogger = log.New(os.Stdout, "ERROR: ", log.Ltime|log.Ldate)
	fatallogger = log.New(os.Stdout, "FATAL: ", log.Ltime|log.Ldate)
)

func Info(v ...any) {
	infologger.Println(v...)
}

func Error(v ...any) {
	errorlogger.Println(v...)
}

func Fatal(v ...any) {
	fatallogger.Fatalln(v...)
}
