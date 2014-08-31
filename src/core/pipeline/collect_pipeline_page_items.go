package pipeline

import (
    "core/common/page_items"
    "core/common/task"
)

type CollectPipelinePageItems struct {
    collector []*page_items.PageItems
}

func NewCollectPipelinePageItems() *CollectPipelinePageItems {
    collector := make([]*page_items.PageItems, 0)
    return &CollectPipelinePageItems{collector: collector}
}

func (this *CollectPipelinePageItems) Process(items *page_items.PageItems, t task.Task) {
    this.collector = append(this.collector, items)
}

func (this *CollectPipelinePageItems) GetCollected() []*page_items.PageItems {
    return this.collector
}
