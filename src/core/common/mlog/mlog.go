package mlog

import (
    "log"
)

func Log(v ...interface{}) {
    log.Panicln(v)
}
