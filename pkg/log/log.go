package log

import (
	"fmt"
	"io"
	"log"
)

type LoggingMode int

const (
	DEBUG LoggingMode = iota
	INFO
	WARNING
	ERROR
	FATAL
)

type Logger struct {
	Mode LoggingMode
	log  *log.Logger
}

func NewLogger(outTo io.Writer, logMode LoggingMode, prefix string) (logger *Logger) {
	return &Logger{
		Mode: logMode,
		log:  log.New(outTo, prefix, log.LstdFlags|log.Lshortfile),
	}
}

// Debug logs only if LoggingMode is set to DEBUG
func (logger *Logger) Debug(msg ...interface{}) {
	if logger.Mode <= DEBUG {
		_ = logger.log.Output(2, "DEBUG:\t"+fmt.Sprint(msg...))
	}
}

// Debugf logs if LoggingMode is set to DEBUG or lower
func (logger *Logger) Debugf(format string, msg ...interface{}) {
	if logger.Mode <= DEBUG {
		_ = logger.log.Output(2, "DEBUG:\t"+fmt.Sprintf(format, msg...))
	}
}

// Info logs if LoggingMode is set to INFO or lower
func (logger *Logger) Info(msg ...interface{}) {
	if logger.Mode <= INFO {
		_ = logger.log.Output(2, "INFO:\t"+fmt.Sprint(msg...))
	}
}

// Infof logs if LoggingMode is set to INFO or lower
func (logger *Logger) Infof(format string, msg ...interface{}) {
	if logger.Mode <= INFO {
		_ = logger.log.Output(2, "INFO:\t"+fmt.Sprintf(format, msg...))
	}
}

// Warning logs if LoggingMode is set to WARNING or lower
func (logger *Logger) Warning(msg ...interface{}) {
	if logger.Mode <= WARNING {
		_ = logger.log.Output(2, "WARNING:\t"+fmt.Sprint(msg...))
	}
}

// Warningf logs if LoggingMode is set to WARNING or lower
func (logger *Logger) Warningf(format string, msg ...interface{}) {
	if logger.Mode <= WARNING {
		_ = logger.log.Output(2, "WARNING:\t"+fmt.Sprintf(format, msg...))
	}
}

// Error logs if LoggingMode is set to ERROR or lower
func (logger *Logger) Error(msg ...interface{}) {
	if logger.Mode <= ERROR {
		_ = logger.log.Output(2, "ERROR:\t"+fmt.Sprint(msg...))
	}
}

// Errorf logs if LoggingMode is set to ERROR or lower
func (logger *Logger) Errorf(format string, msg ...interface{}) {
	if logger.Mode <= ERROR {
		_ = logger.log.Output(2, "Error:\t"+fmt.Sprintf(format, msg...))
	}
}

// Fatal always logs when used
func (logger *Logger) Fatal(msg ...interface{}) {
	_ = logger.log.Output(2, "FATAL:\t"+fmt.Sprint(msg...))
}

// Fatalf always logs when used
func (logger *Logger) Fatalf(format string, msg ...interface{}) {
	_ = logger.log.Output(2, "FATAL:\t"+fmt.Sprintf(format, msg...))
}
