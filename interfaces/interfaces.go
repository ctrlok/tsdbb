package interfaces

import "net/url"

type Req interface{}

type Client interface {
	Send(Req, []byte) error
}

type Basic interface {
	NewRequests(count int)
	Req(i int) Req
	NewClient(url *url.URL) (Client, error)
}
