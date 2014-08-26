// 日志类
package mlog

import (
    "core/common/etc_config"
    "log"
    "os"
    "strconv"
    "time"
)

func Log(v ...interface{}) {
    log.Panicln(v)
}

// 向标准错误输出输出执行过程
type strace struct {
    loginst *log.Logger

    isopen bool
}

var Strace *strace = newStrace()

func newStrace() *strace {
    return &strace{loginst: log.New(os.Stderr, "", log.Llongfile|log.LstdFlags), isopen: true}
}

func (this *strace) Println(str string) {
    if !this.isopen {
        return
    }
    this.loginst.Println(str)
}

func (this *strace) Open() {
    this.isopen = true
}

func (this *strace) Close() {
    this.isopen = false
}

// 文件日志
type filelog struct {
    loginst *log.Logger

    isopen bool
}

var Filelog *filelog

// Get Log instance
func LogInst() *filelog {
    if Filelog == nil {
        Filelog = newFilelog()
    }
    return Filelog
}

// 默认日志路径为WORKDIR/log/log.2011-01-01
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

    return &filelog{loginst: log.New(f, "", log.Llongfile|log.LstdFlags), isopen: isopen}
}

func (this *filelog) log(lable string, str string) {
    if !this.isopen {
        return
    }
    this.loginst.Println(lable + " " + str)
}

func (this *filelog) LogError(str string) {
    this.log("[ERROR]", str)
}

func (this *filelog) LogInfo(str string) {
    this.log("[INFO]", str)
}

func (this *filelog) Open() {
    this.isopen = true
}

func (this *filelog) Close() {
    this.isopen = false
}
