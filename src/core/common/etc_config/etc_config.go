package etc_config

import (
    "core/common/config"
    "core/common/util"
    "os"
)

// Get default config path "WD/etc/main.conf"
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

var Config *config.Config
var path string

// Used in Spider for initialization at first time.
func StartConf(configFilePath string) *config.Config {
    if configFilePath != "" && !util.IsFileExists(configFilePath) {
        panic("config path is not valiad:" + configFilePath)
    }

    path = configFilePath
    return Conf()
}

// Get config instance
func Conf() *config.Config {
    if Config == nil {
        if path == "" {
            path = configpath()
        }
        Config = config.NewConfig().Load(path)
    }
    return Config
}
