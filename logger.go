package tbot

import (
	"log"
)

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Printf(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)

	Debug(args ...any)
	Info(args ...any)
	Print(args ...any)
	Warn(args ...any)
	Error(args ...any)
}

type nopLogger struct{}

func (nopLogger) Debugf(format string, args ...any) {}
func (nopLogger) Infof(format string, args ...any)  {}
func (nopLogger) Printf(format string, args ...any) {}
func (nopLogger) Warnf(format string, args ...any)  {}
func (nopLogger) Errorf(format string, args ...any) {}
func (nopLogger) Debug(args ...any)                 {}
func (nopLogger) Info(args ...any)                  {}
func (nopLogger) Print(args ...any)                 {}
func (nopLogger) Warn(args ...any)                  {}
func (nopLogger) Error(args ...any)                 {}

type BasicLogger struct{}

func (BasicLogger) Debugf(format string, args ...any) { log.Printf(format, args...) }
func (BasicLogger) Infof(format string, args ...any)  { log.Printf(format, args...) }
func (BasicLogger) Printf(format string, args ...any) { log.Printf(format, args...) }
func (BasicLogger) Warnf(format string, args ...any)  { log.Printf(format, args...) }
func (BasicLogger) Errorf(format string, args ...any) { log.Printf(format, args...) }
func (BasicLogger) Debug(args ...any)                 { log.Print(args...) }
func (BasicLogger) Info(args ...any)                  { log.Print(args...) }
func (BasicLogger) Print(args ...any)                 { log.Print(args...) }
func (BasicLogger) Warn(args ...any)                  { log.Print(args...) }
func (BasicLogger) Error(args ...any)                 { log.Print(args...) }
