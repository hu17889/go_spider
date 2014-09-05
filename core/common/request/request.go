// Package request implements request entity contains url and other relevant informaion.
package request

// Request represents object waiting for being crawled.
type Request struct {
    url      string
    respType string
}

// NewRequest returns initialized Request object.
// The respType is "json" or "html"
func NewRequest(url string, respType string) *Request {
    return &Request{url, respType}
}

func (this *Request) GetUrl() string {
    return this.url
}

func (this *Request) GetResponceType() string {
    return this.respType
}
