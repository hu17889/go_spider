package etc_config

import (
    "core/common/config"
    "os"
)

func configpath() string {
    wd, _ := os.Getwd()
    logpath := wd + "/etc/"
    filename := "main.conf"
    err := os.MkdirAll(logpath, 0755)
    if err != nil {
        panic("logpath error : " + logpath + "\n")
    }
    return logpath + filename
}

var Config *config.Config = config.NewConfig().Load(configpath())
