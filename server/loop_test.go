package server

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"time"

	"github.com/ctrlok/tsdbb/interfaces"
	"github.com/ctrlok/tsdbb/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	log.Logger = zap.NewNop()
	log.SLogger = log.Logger.Sugar()
}

type testCli struct {
	returnErr           bool
	errCount, succCount int
}

func (c *testCli) Send(interfaces.Req, []byte) error {
	if c.returnErr {
		c.errCount++
		return fmt.Errorf("error")
	}
	c.succCount++
	return nil
}

type testBasic struct {
	reqCount      int
	reqNumAsked   []int
	failCreateCli bool
	cliCount      int
	cliCountErr   int
}

func (b *testBasic) NewRequests(count int) {
	panic("implement me")
}

func (b *testBasic) Req(i int) interfaces.Req {
	b.reqCount++
	b.reqNumAsked = append(b.reqNumAsked, i)
	return ""
}

func (b *testBasic) NewClient(url *url.URL) (interfaces.Client, error) {
	b.cliCount++
	if b.failCreateCli {
		b.cliCountErr++
		return nil, fmt.Errorf("errorCli")
	}
	return &testCli{}, nil
}

func TestStartClient_Nothing(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	cli := testCli{returnErr: false}
	basic := testBasic{}

	bus <- busMessage{}
	close(bus)
	startClient(context.Background(), &cli, &basic, bus, st)
	assert.Equal(t, 0, cli.succCount)
	assert.Equal(t, 0, basic.reqCount)
}

func TestStartClient_SimpleSucc(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	cli := testCli{returnErr: false}
	basic := testBasic{}

	bus <- busMessage{context.Background(), 0, 1, []byte{}}
	close(bus)
	startClient(context.Background(), &cli, &basic, bus, st)
	assert.Equal(t, 1, cli.succCount)
	assert.Equal(t, 1, basic.reqCount)
	assert.Contains(t, basic.reqNumAsked, 0)
	assert.NotContains(t, basic.reqNumAsked, 1)
}

func TestStartClient_SimpleSuccThreeSend(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	cli := testCli{returnErr: false}
	basic := testBasic{}

	bus <- busMessage{context.Background(), 0, 3, []byte{}}
	close(bus)
	startClient(context.Background(), &cli, &basic, bus, st)
	assert.Equal(t, 3, cli.succCount)
	assert.Equal(t, 3, basic.reqCount)
	assert.Contains(t, basic.reqNumAsked, 0)
	assert.Contains(t, basic.reqNumAsked, 1)
	assert.Contains(t, basic.reqNumAsked, 2)
	assert.NotContains(t, basic.reqNumAsked, 3)
}

func TestStartClient_SimpleSuccTwoBus(t *testing.T) {
	bus := make(chan busMessage, 2)
	st := make(chan statMessage, 2)
	cli := testCli{returnErr: false}
	basic := testBasic{}

	bus <- busMessage{context.Background(), 0, 3, []byte{}}
	bus <- busMessage{context.Background(), 3, 3, []byte{}}
	close(bus)
	startClient(context.Background(), &cli, &basic, bus, st)
	assert.Equal(t, 6, cli.succCount)
	assert.Equal(t, 6, basic.reqCount)
	assert.Contains(t, basic.reqNumAsked, 0)
	assert.Contains(t, basic.reqNumAsked, 1)
	assert.Contains(t, basic.reqNumAsked, 5)
	assert.NotContains(t, basic.reqNumAsked, 6)
}

func TestStartClient_BreakContext(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	cli := testCli{returnErr: false}
	basic := testBasic{}
	ctx, cacnel := context.WithCancel(context.Background())
	cacnel()

	bus <- busMessage{ctx, 1, 200, []byte{}}
	close(bus)
	startClient(context.Background(), &cli, &basic, bus, st)
	assert.Equal(t, 99, cli.succCount)
	assert.Equal(t, 99, basic.reqCount)
	assert.Contains(t, basic.reqNumAsked, 1)
	assert.Contains(t, basic.reqNumAsked, 99)
	assert.NotContains(t, basic.reqNumAsked, 100)
}

func TestStartClient_StatOneSucc(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	cli := testCli{returnErr: false}
	basic := testBasic{}

	bus <- busMessage{context.Background(), 0, 1, []byte{}}
	close(bus)
	startClient(context.Background(), &cli, &basic, bus, st)
	statMsg := <-st
	assert.Equal(t, 1, statMsg.succ)
	assert.Equal(t, 0, statMsg.err)
}

func TestStartClient_StatManySucc(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	cli := testCli{returnErr: false}
	basic := testBasic{}

	bus <- busMessage{context.Background(), 0, 50, []byte{}}
	close(bus)
	startClient(context.Background(), &cli, &basic, bus, st)
	statMsg := <-st
	assert.Equal(t, 50, statMsg.succ)
	assert.Equal(t, 0, statMsg.err)
}

func TestStartClient_StatOneErr(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	cli := testCli{returnErr: true}
	basic := testBasic{}

	bus <- busMessage{context.Background(), 0, 1, []byte{}}
	close(bus)
	startClient(context.Background(), &cli, &basic, bus, st)
	statMsg := <-st
	assert.Equal(t, 0, statMsg.succ)
	assert.Equal(t, 1, statMsg.err)
}

func TestStartClient_StatWithBreakContext(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	cli := testCli{returnErr: false}
	basic := testBasic{}
	ctx, cacnel := context.WithCancel(context.Background())
	cacnel()

	bus <- busMessage{ctx, 1, 200, []byte{}}
	close(bus)
	startClient(context.Background(), &cli, &basic, bus, st)
	statMsg := <-st
	assert.Equal(t, 99, statMsg.succ)
	assert.Equal(t, 0, statMsg.err)
}

func TestStartClients_FailLen(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	ctx := context.Background()
	basic := testBasic{}
	opts := Options{Servers: []string{}}

	err := startClients(ctx, &basic, opts, bus, st)
	assert.Error(t, err, "Should return error if no servers...")
	assert.Zero(t, basic.cliCount)
}

func TestStartClients_FailUrlParse(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	ctx := context.Background()
	basic := testBasic{}
	opts := Options{Servers: []string{"12e:12e;asd"}, Parallel: 1}

	err := startClients(ctx, &basic, opts, bus, st)
	assert.Error(t, err, "Fail with parse bad url")
	assert.IsType(t, &url.Error{}, err, "Fail with parsing bad url")
	assert.Zero(t, basic.cliCount)
}

func TestStartClients_FailCreatingCli(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	ctx := context.Background()
	basic := testBasic{failCreateCli: true}
	opts := Options{Servers: []string{"http://example.com:8080/"}, Parallel: 1}

	err := startClients(ctx, &basic, opts, bus, st)
	assert.Error(t, err, "Should fail on cli creation")
	assert.Equal(t, "errorCli", err.Error(), "Should return error with cli creation")
	assert.NotZero(t, basic.cliCount)
	assert.Equal(t, 1, basic.cliCountErr)
}

func TestStartClients_SuccStartOneClient(t *testing.T) {
	bus := make(chan busMessage, 1)
	st := make(chan statMessage, 1)
	ctx := context.Background()
	basic := testBasic{}
	opts := Options{Servers: []string{"http://example.com:8080/"}, Parallel: 1}
	bus <- busMessage{}
	close(bus)

	err := startClients(ctx, &basic, opts, bus, st)
	assert.NoError(t, err)
	assert.NotZero(t, basic.cliCount)
	time.Sleep(time.Millisecond)
	assert.Zero(t, len(bus))
}

func TestSplitArray_EmptyCount(t *testing.T) {
	ctx := context.Background()
	ar := splitArray(ctx, 0, 100, time.Now())
	assert.Empty(t, ar)

	ar = splitArray(ctx, -100, 100, time.Now())
	assert.Empty(t, ar)
}

func TestSplitArray_SimpleSucc(t *testing.T) {
	ctx := context.Background()
	ar := splitArray(ctx, 100, 1, time.Now())
	assert.NotEmpty(t, ar)
	assert.Len(t, ar, 1)
}

func TestSplitArray_Len(t *testing.T) {
	ctx := context.Background()
	ar := splitArray(ctx, 100, 2, time.Now())
	assert.NotEmpty(t, ar)
	assert.Len(t, ar, 2)
	ar = splitArray(ctx, 101, 2, time.Now())
	assert.NotEmpty(t, ar)
	assert.Len(t, ar, 3)
}

func TestSplitArray_Ctx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ar := splitArray(ctx, 101, 2, time.Now())
	for _, a := range ar {
		assert.Error(t, a.ctx.Err())
	}
}

func TestSplitArray_Time(t *testing.T) {
	ctx := context.Background()
	timeNow, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 3, 2013 at 7:54pm (PST)")
	ar := splitArray(ctx, 101, 2, timeNow)
	for _, a := range ar {
		assert.Equal(t, []byte{0x31, 0x33, 0x35, 0x39, 0x39, 0x32, 0x31, 0x32, 0x34, 0x30}, a.time)
	}
}

func TestSplitArray_StartNSimpleTwo(t *testing.T) {
	ctx := context.Background()
	ar := splitArray(ctx, 100, 2, time.Now())
	assert.Equal(t, 0, ar[0].start)
	assert.Equal(t, 50, ar[0].N)
	assert.Equal(t, 50, ar[1].start)
	assert.Equal(t, 50, ar[1].N)
}

func TestSplitArray_StartNSimpleTree(t *testing.T) {
	ctx := context.Background()
	ar := splitArray(ctx, 101, 2, time.Now())
	assert.Equal(t, 0, ar[0].start)
	assert.Equal(t, 50, ar[0].N)
	assert.Equal(t, 50, ar[1].start)
	assert.Equal(t, 50, ar[1].N)
	assert.Equal(t, 100, ar[2].start)
	assert.Equal(t, 1, ar[2].N)
}

func TestCheckCount_Equal(t *testing.T) {
	control := &ControlMessages{100, 100}
	i := checkCount(100, control)
	assert.Equal(t, 100, i)
}

func TestCheckCount_BiggerBigStep(t *testing.T) {
	control := &ControlMessages{105, 100}
	i := checkCount(100, control)
	assert.Equal(t, 105, i)
}

func TestCheckCount_BiggerLittleStep(t *testing.T) {
	control := &ControlMessages{205, 100}
	i := checkCount(100, control)
	assert.Equal(t, 200, i)
}

func TestCheckCount_LessBigStep(t *testing.T) {
	control := &ControlMessages{105, 100}
	i := checkCount(200, control)
	assert.Equal(t, 105, i)
}

func TestCheckCount_LessLittleStep(t *testing.T) {
	control := &ControlMessages{105, 100}
	i := checkCount(300, control)
	assert.Equal(t, 200, i)
}

func TestStartGenerator_Empty(t *testing.T) {
	ctx := context.Background()
	opts := Options{Parallel: 1, Servers: []string{""}}
	controlChan := make(chan ControlMessages, 1)
	bus := make(chan busMessage, 1)
	tickChan := make(chan time.Time, 1)
	close(tickChan)
	startGenerator(ctx, opts, controlChan, bus, tickChan)
	assert.Empty(t, bus)
}

func TestStartGenerator_Simple(t *testing.T) {
	ctx := context.Background()
	opts := Options{Parallel: 1, Servers: []string{""}, StartCount: 1000}
	controlChan := make(chan ControlMessages, 1)
	bus := make(chan busMessage, 1)
	tickChan := make(chan time.Time, 1)
	tickChan <- time.Time{}
	close(tickChan)

	startGenerator(ctx, opts, controlChan, bus, tickChan)
	assert.Equal(t, 1, len(bus))
	c := <-bus
	assert.Equal(t, []byte{0x2d, 0x36, 0x32, 0x31, 0x33, 0x35, 0x35, 0x39, 0x36, 0x38, 0x30, 0x30},
		c.time, "should be equal time.Time{}")
}

func TestStartGenerator_SimpleTwoTicks(t *testing.T) {
	ctx := context.Background()
	opts := Options{Parallel: 1, Servers: []string{""}, StartCount: 1000}
	controlChan := make(chan ControlMessages, 1)
	bus := make(chan busMessage, 2)
	tickChan := make(chan time.Time, 2)
	tickChan <- time.Time{}
	tickChan <- time.Time{}
	close(tickChan)

	startGenerator(ctx, opts, controlChan, bus, tickChan)
	assert.Equal(t, 2, len(bus))
}

func TestStartGenerator_ControlSimple(t *testing.T) {
	ctx := context.Background()
	opts := Options{Parallel: 1, Servers: []string{""}, StartCount: 10}
	controlChan := make(chan ControlMessages, 1)
	bus := make(chan busMessage, 2)
	tickChan := make(chan time.Time, 2)
	tickChan <- time.Time{}
	close(tickChan)
	controlChan <- ControlMessages{20, 1}

	startGenerator(ctx, opts, controlChan, bus, tickChan)
	assert.Equal(t, 1, len(bus))
	assert.Empty(t, controlChan)
	m := <-bus
	assert.Equal(t, 11, m.N)
}

func TestStartGenerator_ControlSimpleTwoTicksOneControl(t *testing.T) {
	ctx := context.Background()
	opts := Options{Parallel: 1, Servers: []string{""}, StartCount: 10}
	controlChan := make(chan ControlMessages, 1)
	bus := make(chan busMessage, 2)
	tickChan := make(chan time.Time, 2)
	tickChan <- time.Time{}
	tickChan <- time.Time{}
	close(tickChan)
	controlChan <- ControlMessages{20, 1}

	startGenerator(ctx, opts, controlChan, bus, tickChan)
	assert.Equal(t, 2, len(bus))
	assert.Empty(t, controlChan)
	m := <-bus
	assert.Equal(t, 11, m.N)
	m = <-bus
	assert.Equal(t, 12, m.N)
}

func TestStartGenerator_ControlSimpleTwoTicksTwoControl(t *testing.T) {
	ctx := context.Background()
	opts := Options{Parallel: 1, Servers: []string{""}, StartCount: 10}
	controlChan := make(chan ControlMessages, 2)
	bus := make(chan busMessage, 2)
	tickChan := make(chan time.Time, 2)
	tickChan <- time.Time{}
	tickChan <- time.Time{}
	close(tickChan)
	controlChan <- ControlMessages{20, 1}
	controlChan <- ControlMessages{8, 10}

	startGenerator(ctx, opts, controlChan, bus, tickChan)
	assert.Equal(t, 2, len(bus))
	assert.Empty(t, controlChan)
	m := <-bus
	assert.Equal(t, 11, m.N)
	m = <-bus
	assert.Equal(t, 8, m.N)
}
