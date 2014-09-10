package pipeline

import (
    "github.com/hu17889/go_spider/core/common/com_interfaces"
    "github.com/hu17889/go_spider/core/common/page_items"
    "os"
)

type PipelineFile struct {
    pFile *os.File

    path string
}

func NewPipelineFile(path string) *PipelineFile {
    pFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
    if err != nil {
        panic("File '" + path + "' in PipelineFile open failed.")
    }
    return &PipelineFile{path: path, pFile: pFile}
}

func (this *PipelineFile) Process(items *page_items.PageItems, t com_interfaces.Task) {
    this.pFile.WriteString("----------------------------------------------------------------------------------------------\n")
    this.pFile.WriteString("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
    this.pFile.WriteString("Crawled result : \n")
    for key, value := range items.GetAll() {
        this.pFile.WriteString(key + "\t:\t" + value + "\n")
    }
}
