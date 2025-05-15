package utils

import (
	"fmt"
	"log"
	"os"
)

var (
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
	FatalLogger *log.Logger
)

func TuiLoggerStart() error {
	debug := os.Getenv("MCPHEE11_TUI_DEBUG")
	if debug != "true" {
		return nil
	}
	logFile, err := os.OpenFile("mcphee11-tui-log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	FatalLogger = log.New(logFile, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
	return nil
}

func TuiLogger(level, message string) {
	debug := os.Getenv("MCPHEE11_TUI_DEBUG")
	if debug != "true" {
		if level == "Fatal" || level == "Error" {
			fmt.Printf("%s: %s\n", level, message)
		}
		return
	}
	if level == "Fatal" {
		FatalLogger.Println(message)
		os.Exit(1)
	}
	if level == "Error" {
		ErrorLogger.Println(message)
	}
	if level == "Info" {
		InfoLogger.Println(message)
	}
}
