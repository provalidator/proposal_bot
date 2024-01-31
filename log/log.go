package log

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Loggers struct {
	Trace *log.Logger
	Warn  *log.Logger
	Info  *log.Logger
	Error *log.Logger
}

var Logger Loggers

func LogInit() {
	// Open or create a log file (overwrite existing logs)
	logDir := os.Getenv("ROOT_PATH") + "/log"
	logFileName := getLogFileName()
	logFile, err := os.OpenFile(logDir+"/"+logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("Cannot create log file:", err)
	}

	// Log Level and Format Settings
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(logFile)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Setting Output to Standard Output and Log Files
	logWriter := io.MultiWriter(os.Stdout, logFile)
	Logger.Trace = log.New(logWriter, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Info = log.New(logWriter, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Warn = log.New(logWriter, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Error = log.New(logWriter, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func getLogFileName() string {
	// Create log file name based on the current date (for example:23-10-19)
	today := time.Now()
	return today.Format("2006-01-02") + ".log"
}
