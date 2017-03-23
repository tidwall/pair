package pair

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

// Pair is a tightly packed key/value pair
type Pair struct {
	data unsafe.Pointer
}

// New returns a Pair
func New(key, value []byte) Pair {
	var khdr, vhdr byte
	var khdrsize, vhdrsize int
	if len(key) <= 0xFD {
		khdr, khdrsize = byte(len(key)), 0
	} else if len(key) <= 0xFFFF {
		khdr, khdrsize = 0xFE, 2
	} else if len(key) <= 0x7FFFFFFF {
		khdr, khdrsize = 0xFF, 4
	} else {
		panic("key is too large")
	}
	if len(value) <= 0xFD {
		vhdr, vhdrsize = byte(len(value)), 0
	} else if len(value) <= 0xFFFF {
		vhdr, vhdrsize = 0xFE, 2
	} else if len(value) <= 0x7FFFFFFF {
		vhdr, vhdrsize = 0xFF, 4
	} else {
		panic("value is too large")
	}
	slice := makenz(2 + khdrsize + vhdrsize + len(key) + len(value))
	slice[0] = khdr
	slice[1] = vhdr
	if khdrsize > 0 {
		if khdrsize == 2 {
			binary.LittleEndian.PutUint16(slice[2:], uint16(len(key)))
		} else {
			binary.LittleEndian.PutUint32(slice[2:], uint32(len(key)))
		}
	}
	if vhdrsize > 0 {
		if vhdrsize == 2 {
			binary.LittleEndian.PutUint16(slice[2+khdrsize:], uint16(len(value)))
		} else {
			binary.LittleEndian.PutUint32(slice[2+khdrsize:], uint32(len(value)))
		}
	}
	copy(slice[2+khdrsize+vhdrsize:], key)
	copy(slice[2+khdrsize+vhdrsize+len(key):], value)
	return Pair{unsafe.Pointer(&slice[0])}
}

const (
	size  = 0
	key   = 1
	value = 2
)

var alignSize = int(unsafe.Sizeof(uintptr(0)))
var maxInt = int(^uint(0) >> 1)

func (pair Pair) get(what int) ([]byte, int) {
	if uintptr(pair.data) == 0 {
		return nil, 0
	}
	slice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(pair.data),
		Len:  maxInt,
		Cap:  maxInt,
	}))
	khdr, vhdr := slice[0], slice[1]
	var khdrsize, vhdrsize int
	var ksize, vsize int
	if khdr == 0xFE {
		khdrsize = 2
	} else if khdr == 0xFF {
		khdrsize = 4
	}
	if vhdr == 0xFE {
		vhdrsize = 2
	} else if vhdr == 0xFF {
		vhdrsize = 4
	}
	if khdrsize == 0 {
		ksize = int(khdr)
	} else if khdrsize == 2 {
		ksize = int(binary.LittleEndian.Uint16(slice[2:]))
	} else {
		ksize = int(binary.LittleEndian.Uint32(slice[2:]))
	}
	kstart := 2 + khdrsize + vhdrsize
	if what == key {
		slice = slice[kstart : kstart+ksize : kstart+ksize]
		return slice, 0
	}
	if vhdrsize == 0 {
		vsize = int(vhdr)
	} else if vhdrsize == 2 {
		vsize = int(binary.LittleEndian.Uint16(slice[2+khdrsize:]))
	} else {
		vsize = int(binary.LittleEndian.Uint32(slice[2+khdrsize:]))
	}
	vstart := kstart + ksize
	if what == value {
		slice = slice[vstart : vstart+vsize : vstart+vsize]
		return slice, 0
	}
	sz := vstart + vsize
	if sz%alignSize != 0 {
		sz += alignSize - (sz % alignSize)
	}
	return nil, sz
}

// Key returns the key portion of the key
func (pair Pair) Key() []byte {
	s, _ := pair.get(key)
	return s
}

// Value returns the value
func (pair Pair) Value() []byte {
	s, _ := pair.get(value)
	return s
}

// Size returns the size of the in-memory allocation
func (pair Pair) Size() int {
	_, i := pair.get(size)
	return i
}

// Zero return true if the pair is unallocated
func (pair Pair) Zero() bool {
	return uintptr(pair.data) == 0
}

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size, typ uintptr, needzero bool) uintptr

// makenz returns a byte slice that is not zero filled. This can provide a big
// performance boost for large pairs.
func makenz(count int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: mallocgc(uintptr(count), 0, false),
		Len:  count,
		Cap:  count,
	}))
}
