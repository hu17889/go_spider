// Package pipeline is the persistent and offline process part of crawler.
package pipeline

import (
    "github.com/hu17889/go_spider/core/common/com_interfaces"
    "github.com/hu17889/go_spider/core/common/page_items"
)

// The interface Pipeline can be implemented to customize ways of persistent.
type Pipeline interface {
    // The Process implements result persistent.
    // The items has the result be crawled.
    // The t has informations of this crawl task.
    Process(items *page_items.PageItems, t com_interfaces.Task)
}

// The interface CollectPipeline recommend result in process's memory temporarily.
type CollectPipeline interface {
    Pipeline

    // The GetCollected returns result saved in in process's memory temporarily.
    GetCollected() []*page_items.PageItems
}
