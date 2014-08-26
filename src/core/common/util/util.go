package util

import (
    "os"
)

func IsDirExists(path string) bool {
    fi, err := os.Stat(path)

    if err != nil {
        return os.IsExist(err)
    } else {
        return fi.IsDir()
    }

    panic("util isDirExists not reached")
}

func IsFileExists(path string) bool {
    fi, err := os.Stat(path)

    if err != nil {
        return os.IsExist(err)
    } else {
        return !fi.IsDir()
    }

    panic("util isFileExists not reached")
}
