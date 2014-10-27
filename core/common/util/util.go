// Package util contains some common functions of GO_SPIDER project.
package util

import (
    "os"
    "regexp"
)

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
