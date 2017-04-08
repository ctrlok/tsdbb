package graphite

import (
	"strconv"
	"testing"
)

// func TestGenerateMetrics(t *testing.T) {
// 	tsdb := &TSDB{}
// 	for i := 100; i < 101; i++ {
// 		assert.Equal(t, i, len(tsdb.GenerateMetrics(i).(*PregeneratedMetrics).metrics))
// 	}
// }

type tmpPregenerated1 struct {
	m [3000000]tmpMetric1
}

type tmpMetric1 struct {
	n1, n2, n3, n4 int64
}

func generate1(i int) *tmpPregenerated1 {
	p := tmpPregenerated1{}
	var n1, n2, n3, n4 int64 = 0, 0, 0, 0
	for n := 0; n < i; n++ {
		if n4 > 999 {
			n4, n3 = 0, n3+1
		}
		if n3 > 99 {
			n3, n2 = 0, n2+1
		}
		if n2 > 99 {
			n2, n1 = 0, n1+1
		}
		p.m[n].n1, p.m[n].n2, p.m[n].n3, p.m[n].n4 = n1, n2, n3, n4
		n4++
	}
	return &p
}

func send1(i int, p *tmpPregenerated1) ([]byte, []byte, []byte, []byte, []byte, []byte, []byte, []byte) {
	m := p.m[i]
	part1 := []byte("string")
	byte1 := []byte(strconv.FormatInt(m.n1, 10))
	part2 := []byte(".string")
	byte2 := []byte(strconv.FormatInt(m.n2, 10))
	part3 := []byte(".string")
	byte3 := []byte(strconv.FormatInt(m.n3, 10))
	part4 := []byte(".string")
	byte4 := []byte(strconv.FormatInt(m.n4, 10))

	return part1, part2, part3, part4, byte1, byte2, byte3, byte4

}

func BenchmarkT(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = []byte("some.string")
	}
}

func BenchmarkGenerate1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generate1(3000000 - 1)
	}
}

var gen1 = generate1(3000000)

func BenchmarkSend1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		send1(2, gen1)
	}
}

type tmpPregenerated2 struct {
	m [3000000]tmpMetric2
}

type tmpMetric2 struct {
	n1, n2, n3, n4 []byte
}

func generate2(i int) *tmpPregenerated2 {
	p := tmpPregenerated2{}
	var n1, n2, n3, n4 int = 0, 0, 0, 0
	for n := 0; n < i; n++ {
		if n4 > 999 {
			n4, n3 = 0, n3+1
		}
		if n3 > 99 {
			n3, n2 = 0, n2+1
		}
		if n2 > 99 {
			n2, n1 = 0, n1+1
		}
		p.m[n].n1 = []byte(strconv.Itoa(n1))
		p.m[n].n2 = []byte(strconv.Itoa(n2))
		p.m[n].n3 = []byte(strconv.Itoa(n3))
		p.m[n].n4 = []byte(strconv.Itoa(n4))
		n4++
	}
	return &p
}

func send2(i int, p *tmpPregenerated2) ([]byte, []byte, []byte, []byte, []byte, []byte, []byte, []byte) {
	m := p.m[i]
	part1 := []byte("string")
	part2 := []byte(".string")
	part3 := []byte(".string")
	part4 := []byte(".string")

	return part1, part2, part3, part4, m.n1, m.n2, m.n3, m.n4

}

func BenchmarkGenerate2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generate2(3000000)
	}
}

var gen2 = generate2(3000000)

func BenchmarkSend2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		send2(1, gen2)
	}
}
