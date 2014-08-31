package pipeline

import (
    "core/common/page_items"
    "core/common/task"
)

type PipelineConsole struct {
}

func NewPipelineConsole() *PipelineConsole {
    return &PipelineConsole{}
}

func (this *PipelineConsole) Process(items *page_items.PageItems, t task.Task) {
    println("----------------------------------------------------------------------------------------------")
    println("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
    println("Crawled result : ")
    for key, value := range items.GetAll() {
        println(key + "\t:\t" + value)
    }
}
