package self

import (
	"net/url"

	"github.com/ctrlok/tsdbb-2/interfaces"
)

type Req struct{}

type Client struct{}

func (c *Client) Send(interfaces.Req, []byte) error {
	return nil
}

type Basic struct {
	req *Req
}

func (b *Basic) NewRequests(count int) {
	b.req = &Req{}
}

func (b *Basic) Req(i int) interfaces.Req {
	return b.req
}

func (b *Basic) NewClient(url *url.URL) (interfaces.Client, error) {
	return &Client{}, nil
}
