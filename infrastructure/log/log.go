package log

import (
	// golang package
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	// external package
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

const (
	maxLogFile                = 100
	maxLogFileSizeInMegabytes = 50
	maxLogFileAgeInDays       = 180
	logFilePath               = "log/output.log"
)

var zerologInstance *zerolog.Logger

// InitLogger is function to init logger using zerolog
func InitLogger() {
	var writers []io.Writer
	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})

	logFilePath := strings.Split(logFilePath, "/")
	writers = append(writers, &lumberjack.Logger{
		Filename:   path.Join(logFilePath...),
		MaxBackups: maxLogFile,                // files
		MaxSize:    maxLogFileSizeInMegabytes, // megabytes
		MaxAge:     maxLogFileAgeInDays,       // days
	})

	mw := io.MultiWriter(writers...)
	logger := zerolog.New(mw).With().Timestamp().Int("pid", os.Getpid()).Logger()

	zerologInstance = &logger
}

// Error write log with level error with given parameter
func Error(err error, params interface{}, msg string) {
	if zerologInstance == nil {
		log.Println(err, params, msg)
		return
	}

	jsonParams, _ := json.Marshal(params)
	_, file, no, _ := runtime.Caller(1)
	zerologInstance.Err(err).RawJSON("params", jsonParams).Str("caller", fmt.Sprintf("%s:%d", file, no)).Msg(msg)
}

// Info write log with level info with given parameter
func Info(params interface{}, msg string) {
	if zerologInstance == nil {
		log.Println(params, msg)
		return
	}

	jsonParams, _ := json.Marshal(params)
	_, file, no, _ := runtime.Caller(1)
	zerologInstance.Info().RawJSON("params", jsonParams).Str("caller", fmt.Sprintf("%s:%d", file, no)).Msg(msg)
}

// Fatal write log with level fatal with given parameter
func Fatal(err error, params interface{}, msg string) {
	if zerologInstance == nil {
		log.Fatal(err, params, msg)
		return
	}

	jsonParams, _ := json.Marshal(params)
	_, file, no, _ := runtime.Caller(1)
	zerologInstance.Fatal().RawJSON("params", jsonParams).Str("error", err.Error()).Str("caller", fmt.Sprintf("%s:%d", file, no)).Msg(msg)
}

// Fatal write log with level fatal with given parameter
func Fatalf(err error, format string, v ...interface{}) {
	if zerologInstance == nil {
		log.Fatalf(format, v...)
		return
	}

	_, file, no, _ := runtime.Caller(1)
	zerologInstance.Fatal().Str("error", err.Error()).Str("caller", fmt.Sprintf("%s:%d", file, no)).Msg(fmt.Sprintf(format, v...))
}
