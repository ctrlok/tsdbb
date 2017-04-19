package graphite

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClientNull(t *testing.T) {
	b := Basic{DevNull: true}
	cli, err := b.NewClient(&url.URL{})
	assert.NoError(t, err)
	assert.Equal(t, ioutil.Discard, cli.(*Client).f)
}

func TestNewClientConn(t *testing.T) {
	returnChan := make(chan []byte, 1)
	go func() {
		l, err := net.Listen("tcp", ":13313")
		if err != nil {
			t.Skipf("Can't start listening server")
		}
		defer l.Close()
		conn, _ := l.Accept()
		b, _ := bufio.NewReader(conn).ReadBytes(10)
		returnChan <- b
	}()
	time.Sleep(time.Millisecond)

	b := Basic{}
	uri, _ := url.Parse("tcp://127.0.0.1:13313")
	cli, err := b.NewClient(uri)
	assert.NoError(t, err)
	assert.Empty(t, returnChan)
	err = cli.Send(&Req{55, 55, 55, 55, 55, 55, 55, 55, 55}, []byte{55, 55})
	assert.NoError(t, err)
	cli.(*Client).w.(*bufio.Writer).Flush()
	time.Sleep(time.Millisecond)
	assert.Equal(t, 1, len(returnChan))
	buf := <-returnChan
	assert.Equal(t, "777777777\x00\x00\x00 1 77\n", fmt.Sprintf("%s", buf))
}

func TestBewClientErr_NotExist(t *testing.T) {
	b := Basic{}
	uri, _ := url.Parse("tcp://asd1qd3322.csd2")
	cli, err := b.NewClient(uri)
	assert.Nil(t, cli)
	assert.Error(t, err)
	assert.IsType(t, &net.AddrError{}, err)
}

func TestBewClientErr_DialProblem(t *testing.T) {
	b := Basic{}
	uri, _ := url.Parse("tcp://localhost:1")
	cli, err := b.NewClient(uri)
	assert.Nil(t, cli)
	assert.Error(t, err)
	assert.IsType(t, &net.OpError{}, err)
}

func TestNewRequestsSimple(t *testing.T) {
	b := Basic{}
	b.NewRequests(5)
	assert.NotEmpty(t, b.reqs)
	assert.Equal(t, 5, len(b.reqs))
}

func TestNewRequestsSimpleBigger(t *testing.T) {
	b := Basic{}
	b.NewRequests(10000)
	assert.NotEmpty(t, b.reqs)
	assert.Equal(t, 10000, len(b.reqs))
	assert.Equal(t, ".00.00.99.99", fmt.Sprintf("%s", b.reqs[9999]))
}

func TestNewRequestsSimpleBiggerBigger(t *testing.T) {
	b := Basic{}
	b.NewRequests(2000000)
	assert.NotEmpty(t, b.reqs)
	assert.Equal(t, 2000000, len(b.reqs))
	assert.Equal(t, ".01.99.99.99", fmt.Sprintf("%s", b.reqs[1999999]))
}
