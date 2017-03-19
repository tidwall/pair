package pair

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func hdrSize(s string) int {
	if len(s) <= 0xFD {
		return 1
	} else if len(s) <= 0xFFFF {
		return 3
	} else if len(s) <= 0x7FFFFFFF {
		return 5
	}
	panic("out of range")
}

func expectedSizeForPair(pair Pair) int {
	key := pair.Key()
	value := pair.Value()
	sz := hdrSize(key) + hdrSize(value) + len(key) + len(value)
	if sz%alignSize != 0 {
		sz += alignSize - (sz % alignSize)
	}
	return sz
}

func TestBasic(t *testing.T) {
	mallocgc(10, 0, false)
	item := New("hello", "world")
	assert.Equal(t, "hello", item.Key())
	assert.Equal(t, "world", item.Value())
	assert.Equal(t, expectedSizeForPair(item), item.Size())
	item = New("", "world")
	assert.Equal(t, "", item.Key())
	assert.Equal(t, "world", item.Value())
	assert.Equal(t, expectedSizeForPair(item), item.Size())
	item = New("", "")
	assert.Equal(t, "", item.Key())
	assert.Equal(t, "", item.Value())
	assert.Equal(t, expectedSizeForPair(item), item.Size())
}

func TestRandom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100000; i++ {
		key := make([]byte, rand.Int()%300)
		//value := make([]byte, rand.Int()%300+(0xFFFF-150))
		value := make([]byte, rand.Int()%300)
		rand.Read(key)
		rand.Read(value)
		item := New(string(key), string(value))
		ikey := item.Key()
		ivalue := item.Value()
		assert.Equal(t, ikey, string(key))
		assert.Equal(t, ivalue, string(value))
		assert.Equal(t, expectedSizeForPair(item), item.Size())
	}
	for i := 0; i < 10000; i++ {
		key := make([]byte, rand.Int()%300+(0xFFFF-150))
		value := make([]byte, rand.Int()%300+(0xFFFF-150))
		rand.Read(key)
		rand.Read(value)
		item := New(string(key), string(value))
		ikey := item.Key()
		ivalue := item.Value()
		assert.Equal(t, ikey, string(key))
		assert.Equal(t, ivalue, string(value))
		assert.Equal(t, expectedSizeForPair(item), item.Size())
	}
}

func BenchmarkNew2(t *testing.B) {
	s := strings.Repeat("*", 1)
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}
func BenchmarkNew6(t *testing.B) {
	s := strings.Repeat("*", 3)
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}
func BenchmarkNew14(t *testing.B) {
	s := strings.Repeat("*", 7)
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}

func BenchmarkNew62(t *testing.B) {
	s := strings.Repeat("*", 31)
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}

func BenchmarkNew126(t *testing.B) {
	s := strings.Repeat("*", 63)
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}
func BenchmarkNew256(t *testing.B) {
	s := strings.Repeat("*", 128)
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}
func BenchmarkNew1024(t *testing.B) {
	s := strings.Repeat("*", 512)
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}
func BenchmarkNew0xFFFF(t *testing.B) {
	s := strings.Repeat("*", 0xFFFF/2)
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}

func BenchmarkGet62(t *testing.B) {
	s := strings.Repeat("*", 31)
	var pairs []Pair
	for i := 0; i < t.N; i++ {
		pairs = append(pairs, New(s, s))
	}
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		p := pairs[i]
		p.Key()
		p.Value()
	}
}

func BenchmarkGet1024(t *testing.B) {
	s := strings.Repeat("*", 512)
	var pairs []Pair
	for i := 0; i < t.N; i++ {
		pairs = append(pairs, New(s, s))
	}
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		p := pairs[i]
		p.Key()
		p.Value()
	}
}
