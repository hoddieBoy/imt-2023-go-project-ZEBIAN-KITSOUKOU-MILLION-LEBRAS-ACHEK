package logutil

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"runtime"
)

var logger = log.New(os.Stdout, "METEO-AIRPORT", log.LstdFlags)

var (
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	reset  = color.New(color.Reset).SprintFunc()
)

// logMessage formats the log message with the file name and line number
func logMessage(message, level string) string {
	_, file, line, ok := runtime.Caller(2)
	if ok && (level == red(bold("[Error] ")) || level == yellow(bold("[Warn] "))) {
		return fmt.Sprintf("\n\t%s%s:%d %s%s", level, file, line, message, reset())
	}
	return fmt.Sprintf("\n\t%s%s%s", level, message, reset())
}

// Error logs an error message with red color
func Error(format string, v ...interface{}) {
	logger.Print(logMessage(red(fmt.Sprintf(format, v...)), red(bold("[Error] "))))
}

// Warn logs a warning message with yellow color
func Warn(format string, v ...interface{}) {
	logger.Print(logMessage(yellow(fmt.Sprintf(format, v...)), yellow(bold("[Warn] "))))
}

// Info logs an info message with bold formatting
func Info(format string, v ...interface{}) {
	logger.Print(logMessage(green(fmt.Sprintf(format, v...)), green(bold("[Info] "))))
}
