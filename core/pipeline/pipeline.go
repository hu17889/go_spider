// Pipeline is the persistent and offline process part of crawler.
package pipeline

import (
    "core/common/com_interfaces"
    "core/common/page_items"
)

// The interface Pipeline can be implemented to customize ways of persistent.
type Pipeline interface {
    Process(items *page_items.PageItems, t com_interfaces.Task)
}

// The result will not to be persisted, and record in process's memory.
type CollectPipeline interface {
    Pipeline

    GetCollected() []*page_items.PageItems
}
