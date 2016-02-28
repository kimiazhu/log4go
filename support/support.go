package support

import "time"

type support interface {
	StatTimes(filepath string) (atime, ctime, mtime time.Time, err error)
}

var _support support

// GetStatTime returns the times properties corresponding to the given filepath
// NOTE: the atime under windows system may not correct, it maybe the same with
// ctime. (2016-02-26 golang version 1.5.3)
func GetStatTime(filepath string) (atime, ctime, mtime time.Time, err error) {
	return _support.StatTimes(filepath)
}
