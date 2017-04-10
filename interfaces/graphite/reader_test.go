package graphite

import (
	"bufio"
	"io/ioutil"
	"testing"
)

func BenchmarkWriteBuff(b *testing.B) {

	ch := make(chan []byte, 1)
	l := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	f := ioutil.Discard
	w := bufio.NewWriterSize(f, 4*1024)
	for n := 0; n < b.N; n++ {
		ch <- l
		m := <-ch
		w.Write(m)
	}
	close(ch)
}

func BenchmarkWriteByte(b *testing.B) {
	ch2 := make(chan [10]byte, 1)
	l := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	f := ioutil.Discard
	w := bufio.NewWriterSize(f, 4*1024)
	for n := 0; n < b.N; n++ {
		ch2 <- l
		k := <-ch2
		for i := range k {
			w.WriteByte(l[i])
		}
	}
}

func BenchmarkWriteBytePoint(b *testing.B) {
	ch2 := make(chan *[10]byte, 1)
	l := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	f := ioutil.Discard
	w := bufio.NewWriterSize(f, 4*1024)
	for n := 0; n < b.N; n++ {
		ch2 <- &l
		k := <-ch2
		for i := range k {
			w.WriteByte(l[i])
		}
	}
}

func BenchmarkWriteByteAr(b *testing.B) {
	ch2 := make(chan [10]byte, 1)
	l := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	f := ioutil.Discard
	w := bufio.NewWriterSize(f, 4*1024)
	for n := 0; n < b.N; n++ {
		ch2 <- l
		k := <-ch2
		w.Write(k[0:])
	}
}

func BenchmarkWriteBuffPoint(b *testing.B) {

	ch := make(chan *[]byte, 1)
	l := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	f := ioutil.Discard
	w := bufio.NewWriter(f)
	for n := 0; n < b.N; n++ {
		ch <- &l
		m := <-ch
		w.Write(*m)
	}
	close(ch)
}

func BenchmarkWriteString(b *testing.B) {

	ch := make(chan string, 1)
	l := "1234567890"
	f := ioutil.Discard
	w := bufio.NewWriterSize(f, 4*1024)
	for n := 0; n < b.N; n++ {
		ch <- l
		m := <-ch
		w.WriteString(m)
	}
	close(ch)
}
