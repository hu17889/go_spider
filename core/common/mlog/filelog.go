package mlog

import (
    "github.com/hu17889/core/common/etc_config"
    "log"
    "os"
    "strconv"
    "time"
)

// Filelog represents an active object that logs on file to record error or other useful info.
// The filelog info is output to os.Stderr.
// The loginst is an point of logger in Std-Packages.
// The isopen is a label represents whether open filelog or not.
type filelog struct {
    plog

    loginst *log.Logger
}

var flog *filelog

// LogInst get the singleton filelog object.
func LogInst() *filelog {
    if flog == nil {
        flog = newFilelog()
    }
    return flog
}

// The newFilelog returns initialized filelog object.
// The default file path is "WORKDIR/log/log.2011-01-01".
func newFilelog() *filelog {
    logconf := etc_config.Conf().SectionContent("log")
    var isopen bool = false
    if value, ok := logconf["isopen"]; ok && value == "true" {
        isopen = true
    }

    var logpath string
    if value, ok := logconf["logpath"]; ok {
        logpath = value + "/"
    } else {
        file, _ := os.Getwd()
        logpath = file + "/log/"
    }

    year, month, day := time.Now().Date()
    filename := "log." + strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
    err := os.MkdirAll(logpath, 0755)
    if err != nil {
        panic("logpath error : " + logpath + "\n")
    }

    f, err := os.OpenFile(logpath+filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        panic("log file open error : " + logpath + filename + "\n")
    }

    pfilelog := &filelog{}
    pfilelog.loginst = log.New(f, "", log.LstdFlags)
    pfilelog.isopen = isopen
    return pfilelog
}

func (this *filelog) log(lable string, str string) {
    if !this.isopen {
        return
    }
    file, line := this.getCaller()
    this.loginst.Printf("%s:%d: %s %s\n", file, line, lable, str)
}

// LogError logs error info.
func (this *filelog) LogError(str string) {
    this.log("[ERROR]", str)
}

// LogError logs normal info.
func (this *filelog) LogInfo(str string) {
    this.log("[INFO]", str)
}
