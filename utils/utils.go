package utils

import (
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var raindropSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	raindropSourceDir = regexp.MustCompile(`utils.utils\.go`).ReplaceAllString(file, "")
}

func FileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		println(raindropSourceDir)
		if ok && (!strings.HasPrefix(file, raindropSourceDir) || strings.HasSuffix(file, "_test.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}
