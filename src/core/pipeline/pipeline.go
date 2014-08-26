// Pipeline is the persistent and offline process part of crawler.
package pipeline

import (
    "core/common/page_items"
)

// The interface Pipeline can be implemented to customize ways of persistent.
type Pipeline interface {
    Process(items *page_items.PageItems)
}
