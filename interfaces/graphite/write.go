package graphite

import (
	"bufio"
	"io"
	"net/url"

	"io/ioutil"
	"net"

	"github.com/ctrlok/tsdbb-2/interfaces"
)

type Req [12]byte

type Basic struct {
	DevNull bool
	reqs    []Req
}

func (b *Basic) NewRequests(count int) {
	b.reqs = make([]Req, count)
	for n := range b.reqs[0] {
		b.reqs[0][n] = 48
	}
	for n := 1; n < count; n++ {
		var plus byte = 1
		b.reqs[n][0] = 46
		for k := 11; k > 0; k-- {
			// Every third symbol is a dot
			if k%3 == 0 {
				b.reqs[n][k] = 46
				continue
			}
			if plus == 1 {
				if b.reqs[n-1][k] == 57 {
					b.reqs[n][k] = 48
				} else {
					b.reqs[n][k] = b.reqs[n-1][k] + plus
					plus = 0
				}
			} else {
				b.reqs[n][k] = b.reqs[n-1][k]
			}
		}
	}
}

func (b *Basic) Req(i int) interfaces.Req {
	return &b.reqs[i]
}

func (b *Basic) NewClient(uri *url.URL) (interfaces.Client, error) {
	client := Client{}
	if b.DevNull {
		// client.f, _ = os.OpenFile("/tmp/metricTEst", os.O_RDWR, 0755)
		client.f = ioutil.Discard
		client.w = bufio.NewWriterSize(client.f, 4*1024)
		return &client, nil
	}
	addr, err := net.ResolveTCPAddr(uri.Scheme, uri.Host)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP(uri.Scheme, nil, addr)
	if err != nil {
		return nil, err
	}
	client.f = conn
	client.w = bufio.NewWriter(client.f)
	return &client, err
}

type Client struct {
	w      *bufio.Writer
	f      io.Writer
	prefix []byte

	host string
}

var helpNum = []byte{32, 49, 32}

func (c *Client) Send(req interfaces.Req, time []byte) (err error) {
	m := req.(*Req)
	_, err = c.w.Write(c.prefix)
	if err != nil {
		return err
	}
	_, err = c.w.Write(m[0:])
	if err != nil {
		return err
	}
	_, err = c.w.Write(helpNum)
	if err != nil {
		return err
	}
	_, err = c.w.Write(time)
	return err
}
