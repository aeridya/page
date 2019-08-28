package page

import (
	"fmt"
	"net/http/httputil"

	"github.com/aeridya/core"
)

type Paging interface {
	LoadPage() error
	Get(resp *core.Response)
	Put(resp *core.Response)
	Post(resp *core.Response)
	Delete(resp *core.Response)
	Options(resp *core.Response)
	Head(resp *core.Response)
	Unsupported(resp *core.Response)
}

type Page struct {
	Route string
	Title string

	options []string
}

func (p *Page) LoadPage() error {
	return nil
}

func (p *Page) Get(resp *core.Response) {
	p.undefined(resp)
}

func (p *Page) Put(resp *core.Response) {
	p.undefined(resp)
}

func (p *Page) Post(resp *core.Response) {
	p.undefined(resp)
}

func (p *Page) Delete(resp *core.Response) {
	p.undefined(resp)
}

func (p *Page) OnOptions(opts ...string) {
	p.options = make([]string, len(opts)+2)
	for i := range p.options {
		if i < len(opts) {
			p.options[i] = opts[i]
		} else {
			p.options[i] = "HEAD"
			p.options[i+1] = "OPTIONS"
			return
		}
	}
}

func (p *Page) Head(resp *core.Response) {
	requestDump, err := httputil.DumpRequest(resp.R, false)
	if err != nil {
		resp.Status = 500
		resp.Err = err
		resp.Data = resp
		return
	}
	resp.Status = 200
	fmt.Fprintf(resp.W, "%s\n", string(requestDump))
}

func (p *Page) Options(resp *core.Response) {
	if p.options == nil {
		p.undefined(resp)
		return
	}
	resp.Good(200)
	fmt.Fprintf(resp.W, "%s\n", p.options)
}

func (p *Page) undefined(resp *core.Response) {
	resp.Bad(400, "Bad Request "+resp.R.Method+" to "+resp.R.URL.Path)
	fmt.Fprintf(resp.W, "Error: %d\n%s\n", resp.Status, resp.Err)
	resp.Data = resp
}

func (p *Page) Unsupported(resp *core.Response) {
	resp.Bad(418, "Unsupported Request "+resp.R.Method+" to "+resp.R.URL.Path)
	fmt.Fprintf(resp.W, "Error: %d\n%s\nI'M A TEAPOT!\n", resp.Status, resp.Err)
	resp.Data = resp
}

func ServePage(resp *core.Response, p Paging) {
	switch resp.R.Method {
	case "GET":
		p.Get(resp)
	case "PUT":
		p.Put(resp)
	case "POST":
		p.Post(resp)
	case "DELETE":
		p.Delete(resp)
	case "OPTIONS":
		p.Options(resp)
	case "HEAD":
		p.Head(resp)
	default:
		p.Unsupported(resp)
	}
}
