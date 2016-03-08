// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type xmlProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type xmlFilter struct {
	Enabled  string        `xml:"enabled,attr"`
	Tag      string        `xml:"tag"`
	Level    string        `xml:"level"`
	Type     string        `xml:"type"`
	Property []xmlProperty `xml:"property"`
	Exclude  []string      `xml:"exclude"`
}

type xmlLoggerConfig struct {
	Filter []xmlFilter `xml:"filter"`
}

// jsonConfig used to setup a single filelog.
// this usually for a client app.
type jsonConfig struct {
	Level    string `json:"level"`
	Filename string `json:"filename"`
	Format   string `json:"format"`
	Maxlines string `json:"maxlines"`
	Maxsize  string `json:"maxsize"`
	Excludes string `json:"excludes"`
}

// Setup Can set the file logger by a json config string.
// It is shorthand for config the filelog without an xml
// config file.(...For NOW).
//
// This is usually used for a client app which always write
// a single log file.
//
// config example:
//	{
//  "level":"DEBUG",
//	"filename":"log/all.log",
//  "format":"[%D %T] [%L] (%S) %M",
//	"maxlines":"100K",
//	"maxsize":"100M",
//  "excludes":"github.com/example,github.com/example2"
//	}
// NOTE: rotate and daily properties will allways true
func (log Logger) SetupFileLog(cnf *jsonConfig) {
	lvl, bad := convertLevel(cnf.Level, cnf.Filename)
	excludes := strings.Split(cnf.Excludes, ",")
	props := []xmlProperty{{"filename",cnf.Filename}, {"maxlines", cnf.Maxlines}, {"maxsize", cnf.Maxsize}, {"daily", "true"}, {"rotate", "true"}}
	filt, good := xmlToFileLogWriter(cnf.Filename, excludes, props, true)
	if bad || !good {
		os.Exit(1)
	}
	log["file"] = &Filter{lvl, filt, excludes}
}

// Load XML configuration; see examples/example.xml for documentation
func (log Logger) LoadConfiguration(filename string) {
	fmt.Fprintf(os.Stdout, "Load log4go configuration: %s\n", filename)
	log.Close()

	// Open the configuration file
	fd, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not open %q for reading: %s\n", filename, err)
		os.Exit(1)
	}

	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not read %q: %s\n", filename, err)
		os.Exit(1)
	}

	xc := new(xmlLoggerConfig)
	if err := xml.Unmarshal(contents, xc); err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not parse XML configuration in %q: %s\n", filename, err)
		os.Exit(1)
	}

	for _, xmlfilt := range xc.Filter {
		var filt LogWriter
		var lvl Level
		bad, good, enabled := false, true, false

		// Check required children
		if len(xmlfilt.Enabled) == 0 {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required attribute %s for filter missing in %s\n", "enabled", filename)
			bad = true
		} else {
			enabled = xmlfilt.Enabled != "false"
		}
		if len(xmlfilt.Tag) == 0 {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required child <%s> for filter missing in %s\n", "tag", filename)
			bad = true
		}
		if len(xmlfilt.Type) == 0 {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required child <%s> for filter missing in %s\n", "type", filename)
			bad = true
		}
		if len(xmlfilt.Level) == 0 {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required child <%s> for filter missing in %s\n", "level", filename)
			bad = true
		}

		lvl, bad = convertLevel(xmlfilt.Level, filename)

		// Just so all of the required attributes are errored at the same time if missing
		if bad {
			os.Exit(1)
		}

		switch xmlfilt.Type {
		case "console":
			filt, good = xmlToConsoleLogWriter(filename, xmlfilt.Exclude, xmlfilt.Property, enabled)
		case "file":
			filt, good = xmlToFileLogWriter(filename, xmlfilt.Exclude, xmlfilt.Property, enabled)
		case "xml":
			filt, good = xmlToXMLLogWriter(filename, xmlfilt.Exclude, xmlfilt.Property, enabled)
		case "socket":
			filt, good = xmlToSocketLogWriter(filename, xmlfilt.Exclude, xmlfilt.Property, enabled)
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not load XML configuration in %s: unknown filter type \"%s\"\n", filename, xmlfilt.Type)
			os.Exit(1)
		}

		// Just so all of the required params are errored at the same time if wrong
		if !good {
			os.Exit(1)
		}

		// If we're disabled (syntax and correctness checks only), don't add to logger
		if !enabled {
			continue
		}

		log[xmlfilt.Tag] = &Filter{lvl, filt, xmlfilt.Exclude}
	}
}

func convertLevel(level, filename string) (lvl Level, bad bool) {
	switch level {
	case "ACCESS":
		lvl = ACCESS
	case "FINEST":
		lvl = FINEST
	case "FINE":
		lvl = FINE
	case "DEBUG":
		lvl = DEBUG
	case "TRACE":
		lvl = TRACE
	case "INFO":
		lvl = INFO
	case "WARNING":
		lvl = WARNING
	case "ERROR":
		lvl = ERROR
	case "CRITICAL":
		lvl = CRITICAL
	default:
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required child <%s> for filter has unknown value in %s: %s\n", "level", filename, level)
		bad = true
	}
	return
}

func xmlToConsoleLogWriter(filename string, excludes []string, props []xmlProperty, enabled bool) (*ConsoleLogWriter, bool) {
	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Warning: Unknown property \"%s\" for console filter in %s\n", prop.Name, filename)
		}
	}

	// If it's disabled, we're just checking syntax
	if !enabled {
		return nil, true
	}

	return NewConsoleLogWriter(), true
}

// Parse a number with K/M/G suffixes based on thousands (1000) or 2^10 (1024)
func strToNumSuffix(str string, mult int) int {
	num := 1
	if len(str) > 1 {
		switch str[len(str)-1] {
		case 'G', 'g':
			num *= mult
			fallthrough
		case 'M', 'm':
			num *= mult
			fallthrough
		case 'K', 'k':
			num *= mult
			str = str[0 : len(str)-1]
		}
	}
	parsed, _ := strconv.Atoi(str)
	return parsed * num
}

func xmlToFileLogWriter(filename string, excludes []string, props []xmlProperty, enabled bool) (*FileLogWriter, bool) {
	file := ""
	format := "[%D %T] [%L] (%S) %M"
	maxlines := 0
	maxsize := 0
	daily := false
	rotate := false

	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		case "filename":
			abspath, _ := exec.LookPath(os.Args[0])
			dir := filepath.Dir(abspath)
			file = filepath.Join(dir, strings.Trim(prop.Value, " \r\n"))
			if _, err := os.Stat(filepath.Dir(file)); os.IsNotExist(err) {
				os.MkdirAll(filepath.Dir(file), os.ModeDir|os.ModePerm)
			}
		case "format":
			format = strings.Trim(prop.Value, " \r\n")
		case "maxlines":
			maxlines = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1000)
		case "maxsize":
			maxsize = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1024)
		case "daily":
			daily = strings.Trim(prop.Value, " \r\n") != "false"
		case "rotate":
			rotate = strings.Trim(prop.Value, " \r\n") != "false"
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Warning: Unknown property \"%s\" for file filter in %s\n", prop.Name, filename)
		}
	}

	// Check properties
	if len(file) == 0 {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for file filter missing in %s\n", "filename", filename)
		return nil, false
	}

	// If it's disabled, we're just checking syntax
	if !enabled {
		return nil, true
	}

	flw := NewFileLogWriter(file, rotate)
	flw.SetFormat(format)
	flw.SetRotateLines(maxlines)
	flw.SetRotateSize(maxsize)
	flw.SetRotateDaily(daily)
	return flw, true
}

func xmlToXMLLogWriter(filename string, excludes []string, props []xmlProperty, enabled bool) (*FileLogWriter, bool) {
	file := ""
	maxrecords := 0
	maxsize := 0
	daily := false
	rotate := false

	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		case "filename":
			abspath, _ := exec.LookPath(os.Args[0])
			dir := filepath.Dir(abspath)
			file = filepath.Join(dir, strings.Trim(prop.Value, " \r\n"))
		case "maxrecords":
			maxrecords = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1000)
		case "maxsize":
			maxsize = strToNumSuffix(strings.Trim(prop.Value, " \r\n"), 1024)
		case "daily":
			daily = strings.Trim(prop.Value, " \r\n") != "false"
		case "rotate":
			rotate = strings.Trim(prop.Value, " \r\n") != "false"
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Warning: Unknown property \"%s\" for xml filter in %s\n", prop.Name, filename)
		}
	}

	// Check properties
	if len(file) == 0 {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for xml filter missing in %s\n", "filename", filename)
		return nil, false
	}

	// If it's disabled, we're just checking syntax
	if !enabled {
		return nil, true
	}

	xlw := NewXMLLogWriter(file, rotate)
	xlw.SetRotateLines(maxrecords)
	xlw.SetRotateSize(maxsize)
	xlw.SetRotateDaily(daily)
	return xlw, true
}

func xmlToSocketLogWriter(filename string, exclude []string, props []xmlProperty, enabled bool) (SocketLogWriter, bool) {
	endpoint := ""
	protocol := "udp"

	// Parse properties
	for _, prop := range props {
		switch prop.Name {
		case "endpoint":
			endpoint = strings.Trim(prop.Value, " \r\n")
		case "protocol":
			protocol = strings.Trim(prop.Value, " \r\n")
		default:
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Warning: Unknown property \"%s\" for file filter in %s\n", prop.Name, filename)
		}
	}

	// Check properties
	if len(endpoint) == 0 {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for file filter missing in %s\n", "endpoint", filename)
		return nil, false
	}

	// If it's disabled, we're just checking syntax
	if !enabled {
		return nil, true
	}

	return NewSocketLogWriter(protocol, endpoint), true
}
