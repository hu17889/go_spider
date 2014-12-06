// Package util contains some common functions of GO_SPIDER project.
package util

import (
    "os"
    "regexp"
    "strings"
)

// JsonpToJson modify jsonp string to json string
// Example: forbar({a:"1",b:2}) to {"a":"1","b":2}
func JsonpToJson(json string) string {
    start := strings.Index(json, "{")
    end := strings.LastIndex(json, "}")
    start1 := strings.Index(json, "[")
    if start1 > 0 && start > start1 {
        start = start1
        end = strings.LastIndex(json, "]")
    }
    if end > start && end != -1 && start != -1 {
        json = json[start : end+1]
    }
    json = strings.Replace(json, "\\'", "", -1)
    regDetail, _ := regexp.Compile("([^\\s\\:\\{\\,\\d\"]+|[a-z][a-z\\d]*)\\s*\\:")
    return regDetail.ReplaceAllString(json, "\"$1\":")
}

// The GetWDPath gets the work directory path.
func GetWDPath() string {
    wd := os.Getenv("GOPATH")
    if wd == "" {
        panic("GOPATH is not setted in env.")
    }
    return wd
}

// The IsDirExists judges path is directory or not.
func IsDirExists(path string) bool {
    fi, err := os.Stat(path)

    if err != nil {
        return os.IsExist(err)
    } else {
        return fi.IsDir()
    }

    panic("util isDirExists not reached")
}

// The IsFileExists judges path is file or not.
func IsFileExists(path string) bool {
    fi, err := os.Stat(path)

    if err != nil {
        return os.IsExist(err)
    } else {
        return !fi.IsDir()
    }

    panic("util isFileExists not reached")
}

// The IsNum judges string is number or not.
func IsNum(a string) bool {
    reg, _ := regexp.Compile("^\\d+$")
    return reg.MatchString(a)
}
