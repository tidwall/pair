package pair

import (
	"math/rand"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func hdrSize(s []byte) int {
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
	item := New([]byte("hello"), []byte("world"))
	assert.Equal(t, "hello", string(item.Key()))
	assert.Equal(t, "world", string(item.Value()))
	assert.Equal(t, expectedSizeForPair(item), item.Size())
	item = New([]byte(""), []byte("world"))
	assert.Equal(t, "", string(item.Key()))
	assert.Equal(t, "world", string(item.Value()))
	assert.Equal(t, expectedSizeForPair(item), item.Size())
	item = New([]byte(""), nil)
	assert.Equal(t, "", string(item.Key()))
	assert.Equal(t, "", string(item.Value()))
	assert.Equal(t, expectedSizeForPair(item), item.Size())
	assert.Equal(t, false, item.Zero())
	assert.Equal(t, true, (Pair{}).Zero())

}

func TestRandom(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var items []Pair
	var keys []string
	var values []string
	var sizes []int
	for i := 0; i < 100000; i++ {
		key := make([]byte, rand.Int()%300)
		rand.Read(key)
		keys = append(keys, string(key))
		value := make([]byte, rand.Int()%300)
		rand.Read(value)
		values = append(values, string(value))
		item := New([]byte(keys[i]), []byte(values[i]))
		items = append(items, item)
		sizes = append(sizes, item.Size())
		ikey := item.Key()
		ivalue := item.Value()
		assert.Equal(t, string(key), string(ikey))
		assert.Equal(t, string(value), string(ivalue))
		assert.Equal(t, expectedSizeForPair(item), item.Size())
	}
	// check for memory span errors
	var ikeys [][]byte
	var ivalues [][]byte
	runtime.GC()
	for i := 0; i < len(items); i++ {
		item := items[i]
		key := keys[i]
		value := values[i]
		ikey := item.Key()
		ikeys = append(ikeys, ikey)
		ivalue := item.Value()
		ivalues = append(ivalues, ivalue)
		assert.Equal(t, string(key), string(ikey))
		assert.Equal(t, string(value), string(ivalue))
		assert.Equal(t, expectedSizeForPair(item), item.Size())
	}
	items = nil
	// check for memory span errors
	runtime.GC()
	for i := 0; i < len(items); i++ {
		key := keys[i]
		value := values[i]
		ikey := ikeys[i]
		ivalue := ivalues[i]
		assert.Equal(t, string(key), string(ikey))
		assert.Equal(t, string(value), string(ivalue))
		assert.Equal(t, key, ikey)
		assert.Equal(t, value, ivalue)
	}
	// larger items
	for i := 0; i < 10000; i++ {
		key := make([]byte, rand.Int()%300+(0xFFFF-150))
		value := make([]byte, rand.Int()%300+(0xFFFF-150))
		rand.Read(key)
		rand.Read(value)
		item := New(key, value)
		ikey := item.Key()
		ivalue := item.Value()
		assert.Equal(t, string(key), string(ikey))
		assert.Equal(t, string(value), string(ivalue))
		assert.Equal(t, expectedSizeForPair(item), item.Size())
	}
}

func BenchmarkNew2(t *testing.B) {
	s := []byte(strings.Repeat("*", 1))
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}
func BenchmarkNew6(t *testing.B) {
	s := []byte(strings.Repeat("*", 3))
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}
func BenchmarkNew14(t *testing.B) {
	s := []byte(strings.Repeat("*", 7))
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}

func BenchmarkNew62(t *testing.B) {
	s := []byte(strings.Repeat("*", 31))
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}

func BenchmarkNew126(t *testing.B) {
	s := []byte(strings.Repeat("*", 63))
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}
func BenchmarkNew256(t *testing.B) {
	s := []byte(strings.Repeat("*", 128))
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}
func BenchmarkNew1024(t *testing.B) {
	s := []byte(strings.Repeat("*", 512))
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}
func BenchmarkNew0xFFFF(t *testing.B) {
	s := []byte(strings.Repeat("*", 0xFFFF/2))
	t.ReportAllocs()
	t.ReportAllocs()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		New(s, s)
	}
}

func BenchmarkGet62(t *testing.B) {
	s := []byte(strings.Repeat("*", 31))
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
	s := []byte(strings.Repeat("*", 512))
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
