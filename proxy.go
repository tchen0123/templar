package templar

import (
	"io"
	"net/http"
	"time"
)

type Proxy struct {
	client Client
	stats  Stats
}

func NewProxy(cl Client, stats Stats) *Proxy {
	return &Proxy{cl, stats}
}

type copyResonder struct {
	w http.ResponseWriter
}

func (c *copyResonder) Send(res *http.Response) io.Writer {
	for k, v := range res.Header {
		c.w.Header()[k] = v
	}

	c.w.WriteHeader(res.StatusCode)

	return c.w
}

func (p *Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	start := time.Now()

	p.stats.StartRequest(req)

	p.client.Forward(&copyResonder{res}, req)

	p.stats.Emit(req, time.Since(start))
}
