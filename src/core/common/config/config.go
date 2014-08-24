package config

import (
    "errors"
    "io/ioutil"
    "strconv"
    "strings"
    "time"
)

// 注意，针对 Config 的各种操作，如果需要并发，需要在外层进行同步
type Config struct {
    globalContent   map[string]string
    sectionContents map[string]map[string]string
    sections        []string
}

func NewConfig() *Config {
    return &Config{
        globalContent:   make(map[string]string),
        sectionContents: make(map[string]map[string]string),
    }
}

// 从配置文件加载配置内容
func (this *Config) Load(configFile string) *Config {
    stream, err := ioutil.ReadFile(configFile)
    if err != nil {
        panic("config read file error : " + configFile + "\n")
    }
    this.LoadString(string(stream))
    return this
}

// 把配置内容写入某配置文件
func (this *Config) Save(configFile string) error {
    // @todo: 如果外层加锁，本调用可能会阻塞外层一段时间。考虑用异步
    return ioutil.WriteFile(configFile, []byte(this.String()), 0777)
}

func (this *Config) Clear() {
    this.globalContent = make(map[string]string)
    this.sectionContents = make(map[string]map[string]string)
    this.sections = nil
}

func (this *Config) LoadString(s string) error {
    lines := strings.Split(s, "\n")
    section := ""
    // @todo: 两遍遍历，第一遍检查语法是否正确，第二遍执行更改
    for _, line := range lines {
        line = strings.Trim(line, emptyRunes)
        if line == "" || line[0] == '#' {
            continue
        }
        if line[0] == '[' {
            if lineLen := len(line); line[lineLen-1] == ']' {
                section = line[1 : lineLen-1]
                sectionAdded := false
                for _, oldSection := range this.sections {
                    if section == oldSection {
                        sectionAdded = true
                        break
                    }
                }
                if !sectionAdded {
                    this.sections = append(this.sections, section)
                }
                continue
            }
        }
        pair := strings.SplitN(line, "=", 2)
        if len(pair) != 2 {
            return errors.New("bad config file syntax")
        }
        key := strings.Trim(pair[0], emptyRunes)
        value := strings.Trim(pair[1], emptyRunes)
        if section == "" {
            this.globalContent[key] = value
        } else {
            if _, ok := this.sectionContents[section]; !ok {
                this.sectionContents[section] = make(map[string]string)
            }
            this.sectionContents[section][key] = value
        }
    }
    return nil
}

func (this *Config) String() string {
    s := ""
    for key, value := range this.globalContent {
        s += key + "=" + value + "\n"
    }
    for section, content := range this.sectionContents {
        s += "[" + section + "]\n"
        for key, value := range content {
            s += key + "=" + value + "\n"
        }
    }
    return s
}

func (this *Config) StringWithMeta() string {
    s := "__sections__=" + strings.Join(this.sections, ",") + "\n"
    return s + this.String()
}

func (this *Config) GlobalHas(key string) bool {
    if _, ok := this.globalContent[key]; ok {
        return true
    }
    return false
}

func (this *Config) GlobalGet(key string) string {
    return this.globalContent[key]
}

func (this *Config) GlobalSet(key string, value string) {
    this.globalContent[key] = value
}

func (this *Config) GlobalGetInt(key string) int {
    value := this.GlobalGet(key)
    if value == "" {
        return 0
    }
    result, err := strconv.Atoi(value)
    if err != nil {
        return 0
    }
    return result
}

func (this *Config) GlobalGetInt64(key string) int64 {
    value := this.GlobalGet(key)
    if value == "" {
        return 0
    }
    result, err := strconv.ParseInt(value, 10, 64)
    if err != nil {
        return 0
    }
    return result
}

func (this *Config) GlobalGetDuration(key string) time.Duration {
    return time.Duration(this.GlobalGetInt(key)) * time.Second
}

func (this *Config) GlobalGetDeadline(key string) time.Time {
    return time.Now().Add(time.Duration(this.GlobalGetInt(key)) * time.Second)
}

func (this *Config) GlobalGetSlice(key string, separator string) []string {
    result := []string{}
    value := this.GlobalGet(key)
    if value != "" {
        for _, part := range strings.Split(value, separator) {
            result = append(result, strings.Trim(part, emptyRunes))
        }
    }
    return result
}

func (this *Config) GlobalGetSliceInt(key string, separator string) []int {
    result := []int{}
    value := this.GlobalGetSlice(key, separator)
    for _, part := range value {
        int, err := strconv.Atoi(part)
        if err != nil {
            continue
        }
        result = append(result, int)
    }
    return result
}

func (this *Config) GlobalContent() map[string]string {
    return this.globalContent
}

func (this *Config) Sections() []string {
    return this.sections
}

func (this *Config) HasSection(section string) bool {
    if _, ok := this.sectionContents[section]; ok {
        return true
    }
    return false
}

func (this *Config) SectionHas(section string, key string) bool {
    if !this.HasSection(section) {
        return false
    }
    if _, ok := this.sectionContents[section][key]; ok {
        return true
    }
    return false
}

func (this *Config) SectionGet(section string, key string) string {
    if content, ok := this.sectionContents[section]; ok {
        return content[key]
    }
    return ""
}

func (this *Config) SectionSet(section string, key string, value string) {
    if content, ok := this.sectionContents[section]; ok {
        content[key] = value
    } else {
        content = make(map[string]string)
        content[key] = value
        this.sectionContents[section] = content
    }
}

func (this *Config) SectionGetInt(section string, key string) int {
    value := this.SectionGet(section, key)
    if value == "" {
        return 0
    }
    result, err := strconv.Atoi(value)
    if err != nil {
        return 0
    }
    return result
}

func (this *Config) SectionGetDuration(section string, key string) time.Duration {
    return time.Duration(this.SectionGetInt(section, key)) * time.Second
}

func (this *Config) SectionGetSlice(section string, key string, separator string) []string {
    result := []string{}
    value := this.SectionGet(section, key)
    if value != "" {
        for _, part := range strings.Split(value, separator) {
            result = append(result, strings.Trim(part, emptyRunes))
        }
    }
    return result
}

func (this *Config) SectionContent(section string) map[string]string {
    return this.sectionContents[section]
}

func (this *Config) SectionContents() map[string]map[string]string {
    return this.sectionContents
}

const emptyRunes = " \r\t\v"
