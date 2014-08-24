//
package page_items

import (
    "core/common/request"
)

type PageItems struct {
    req *request.Request

    items map[string]string

    // Whether send ResultItems to scheduler or not.
    skip bool
}

func NewPageItems(r *request.Request) *PageItems {
    return &PageItems{req: r}
}

func (this *PageItems) AddItem(key string, item string) {
    this.items[key] = item
}

func (this *PageItems) GetItem(key string) string {
    return this.items[key]
}

func (this *PageItems) GetAll() map[string]string {
    return this.items
}

func (this *PageItems) SetSkip(skip bool) *PageItems {
    this.skip = skip
    return this
}

func (this *PageItems) Skip() bool {
    return this.skip
}
