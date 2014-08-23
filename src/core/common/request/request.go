// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package request

type Request struct {
    url      string
    respType string
}

// respType is "json" or "thml"
func NewRequest(url string, respType string) *Request {
    return &Request{url, respType}
}

func (this *Request) GetUrl() string {
    return this.url
}

func (this *Request) GetResponceType() string {
    return this.respType
}
