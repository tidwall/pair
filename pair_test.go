package pair

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestABC(t *testing.T) {
	item := New("hello", "world")
	assert.Equal(t, "hello", item.Key())
	assert.Equal(t, "world", item.Value())
	return
	item = New("", "world")
	assert.Equal(t, "", item.Key())
	assert.Equal(t, "world", item.Value())
	item = New("", "")
	assert.Equal(t, "", item.Key())
	assert.Equal(t, "", item.Value())
	rand.Seed(time.Now().UnixNano())
	N := 100000
	for i := 0; i < N; i++ {
		key := make([]byte, rand.Int()%1024)
		value := make([]byte, rand.Int()%1024)
		rand.Read(key)
		rand.Read(value)
		// push valid data, make sure matches
		item = New(string(key), string(value))
		ikey := item.Key()
		ivalue := item.Value()
		assert.Equal(t, ikey, string(key))
		assert.Equal(t, ivalue, string(value))
	}
}
