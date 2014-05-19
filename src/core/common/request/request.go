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
