// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"errors"
	"fmt"
	. "github.com/kimiazhu/golib/stack"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	Global Logger
)

func init() {
	// auto load config from default position
	Global = NewDefaultLogger(DEBUG)
	file, _ := exec.LookPath(os.Args[0])
	dir := filepath.Dir(file)
	if _, err := os.Stat("log4go.xml"); !os.IsNotExist(err) {
		Global.LoadConfiguration("log4go.xml")
	} else if _, err := os.Stat(filepath.Join(dir, "/log4go.xml")); !os.IsNotExist(err) {
		Global.LoadConfiguration(filepath.Join(dir, "log4go.xml"))
	} else if _, err := os.Stat(filepath.Join(dir, "/conf/log4go.xml")); !os.IsNotExist(err) {
		Global.LoadConfiguration(filepath.Join(dir, "/conf/log4go.xml"))
	} else {
		//fmt.Fprintf(os.Stderr, "log4go config not found, exec dir is: %s, u need to load it by yourself.\n", dir)
	}
}

// setup by config string, not a config file
func Setup(config []byte) {
	Global.Config(config)
}

// Wrapper for (*Logger).LoadConfiguration
func LoadConfiguration(filename string) {
	Global.LoadConfiguration(filename)
}

// Wrapper for (*Logger).AddFilter
func AddFilter(name string, lvl Level, writer LogWriter) {
	Global.AddFilter(name, lvl, writer)
}

// Wrapper for (*Logger).Close (closes and removes all logwriters)
func Close() {
	Global.Close()
}

func Crash(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(CRITICAL, strings.Repeat(" %v", len(args))[1:], args...)
	}
	panic(args)
}

// Logs the given message and crashes the program
func Crashf(format string, args ...interface{}) {
	Global.intLogf(CRITICAL, format, args...)
	Global.Close() // so that hopefully the messages get logged
	panic(fmt.Sprintf(format, args...))
}

// Compatibility with `log`
func Exit(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func Exitf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func Stderr(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func Stderrf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
}

// Compatibility with `log`
func Stdout(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(INFO, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func Stdoutf(format string, args ...interface{}) {
	Global.intLogf(INFO, format, args...)
}

// Send a log message manually
// Wrapper for (*Logger).Log
func Log(lvl Level, source, message string) {
	Global.Log(lvl, source, message)
}

// Send a formatted log message easily
// Wrapper for (*Logger).Logf
func Logf(lvl Level, format string, args ...interface{}) {
	Global.intLogf(lvl, format, args...)
}

// Send a closure log message
// Wrapper for (*Logger).Logc
func Logc(lvl Level, closure func() string) {
	Global.intLogc(lvl, closure)
}

// Utility for finest log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Finest
func Finest(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINEST
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for fine log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Fine
func Fine(arg0 interface{}, args ...interface{}) {
	const (
		lvl = FINE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for debug log messages
// When given a string as the first argument, this behaves like Logf but with the DEBUG log level (e.g. the first argument is interpreted as a format for the latter arguments)
// When given a closure of type func()string, this logs the string returned by the closure iff it will be logged.  The closure runs at most one time.
// When given anything else, the log message will be each of the arguments formatted with %v and separated by spaces (ala Sprint).
// Wrapper for (*Logger).Debug
func Debug(arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for trace log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Trace
func Trace(arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for info log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func Info(arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for Access log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func Access(arg0 interface{}, args ...interface{}) {
	const (
		lvl = ACCESS
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Utility for warn log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Warn
func Warn(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = WARNING
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Utility for error log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Error
func Error(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = ERROR
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Utility for critical log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Critical. This method will log the call stack
func Critical(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = CRITICAL
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg := fmt.Sprintf("%s\n%s", fmt.Sprintf(first, args...), CallStack(3))
		Global.intLogf(lvl, msg)
		//Global.intLogf(lvl, "%s", CallStack(3))
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s\n%s", str, CallStack(3))
		//Global.intLogf(lvl, "%s", CallStack(3))
		return errors.New(str)
	case func(interface{}) string:
		str := first(args[0])
		Global.intLogf(lvl, "%s\n%s", str, CallStack(3))
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		msg := fmt.Sprintf("%s\n%s", fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...), CallStack(3))
		Global.intLogf(lvl, msg)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

// Recover used to log the stack when panic occur.
// usage: defer log4go.Recover("this is a msg: %v", "msg")
// or:
//      defer log4go.Recover(func(err interface{}) string {
//          // ... your code here, return the error message
//          return fmt.Sprintf("recover..v1=%v;v2=%v;err=%v", 1, 2, err)
//      })
func Recover(arg0 interface{}, args ...interface{}) {
	if err := recover(); err != nil {
		switch a := arg0.(type) {
		case func(interface{}) string:
			// the recovered err will pass to this func
			Critical(arg0, append([]interface{}{err}, args)...)
		case string:
			Critical(a+"\n%v", append(args, err)...)
		default:
			Critical(arg0, append(args, err)...)
		}
	}
}

func IsFinestEnabled() bool {
	return isLevelEnabled(FINEST)
}

func IsFineEnabled() bool {
	return isLevelEnabled(FINE)
}

func IsDebugEnabled() bool {
	return isLevelEnabled(DEBUG)
}

func IsTraceEnabled() bool {
	return isLevelEnabled(TRACE)
}

func IsInfoEnabled() bool {
	return isLevelEnabled(INFO)
}

func IsWarnEnabled() bool {
	return isLevelEnabled(WARNING)
}

func IsErrorEnabled() bool {
	return isLevelEnabled(ERROR)
}

func isLevelEnabled(lvl Level) bool {
	enabled := false
	for _, filt := range Global {
		if lvl >= filt.Level {
			// return true if any filt matched
			enabled = true
			break
		}
	}
	return enabled
}