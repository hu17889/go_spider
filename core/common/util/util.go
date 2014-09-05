// Package util contains some common functions of GO_SPIDER project.
package util

import (
    "os"
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
