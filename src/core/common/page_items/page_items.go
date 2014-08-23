// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package page_items

import (
    "core/common/request"
)

type PageItems struct {
    req *request.Request

    items map[string]interface{}
}

func NewPageItems(r *request.Request) *PageItems {
    return &PageItems{req: r}
}

func (this *PageItems) AddItem(key string, item interface{}) {
    this.items[key] = item
}

func (this *PageItems) GetItem(key string) interface{} {
    return items[key]
}

func (this *PageItems) GetAll() map[string]interface{} {
    return items
}
