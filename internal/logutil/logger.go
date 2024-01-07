package logutil

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "METEO-AIRPORT", log.LstdFlags|log.Lshortfile)
var formatOutput = "\n\t%s %s\n"

var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

func Error(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Printf(formatOutput, red(bold("[ERROR]")), red(msg))
}

func Info(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Printf(formatOutput, green(bold("[INFO]")), green(msg))
}

func Warn(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Printf(formatOutput, yellow(bold("[WARN]")), yellow(msg))
}

func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}
