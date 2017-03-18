package pair

import (
	"reflect"
	"unsafe"
)

// Pair is a tightly packed key/value pair
type Pair struct {
	data uintptr
}

// New returns a Pair
func New(key, value string) Pair {
	var khdr, vhdr byte
	var khdrsize, vhdrsize int
	if len(key) < 0xFD {
		khdr, khdrsize = byte(len(key)), 0
	} else if len(key) <= 0xFFFF {
		khdr, khdrsize = 0xFE, 2
	} else if len(key) <= 0x7FFFFFFF {
		khdr, khdrsize = 0xFF, 4
	} else {
		panic("key is too large")
	}
	if len(value) < 0xFD {
		vhdr, vhdrsize = byte(len(value)), 0
	} else if len(value) <= 0xFFFF {
		vhdr, vhdrsize = 0xFE, 2
	} else if len(value) <= 0x7FFFFFFF {
		vhdr, vhdrsize = 0xFF, 4
	} else {
		panic("key is too large")
	}
	slice := make([]byte, 2+khdrsize+vhdrsize+len(key)+len(value))
	slice[0] = khdr
	slice[1] = vhdr
	if khdrsize > 0 {
		if khdrsize == 2 {
			*(*uint16)(unsafe.Pointer(&slice[2])) = uint16(len(key))
		} else {
			*(*uint32)(unsafe.Pointer(&slice[2])) = uint32(len(key))
		}
	}
	if vhdrsize > 0 {
		if vhdrsize == 2 {
			*(*uint16)(unsafe.Pointer(&slice[2+khdrsize])) = uint16(len(value))
		} else {
			*(*uint32)(unsafe.Pointer(&slice[2+khdrsize])) = uint32(len(value))
		}
	}
	copy(slice[2+khdrsize+vhdrsize:], key)
	copy(slice[2+khdrsize+vhdrsize+len(key):], value)
	return Pair{uintptr(unsafe.Pointer(&slice[0]))}
}

func (pair Pair) get(key bool) string {
	if pair.data == 0 {
		return ""
	}
	slice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: pair.data,
		Len:  int(^uint(0) >> 1),
		Cap:  int(^uint(0) >> 1),
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
		ksize = int(*(*uint16)(unsafe.Pointer(&slice[2])))
	} else {
		ksize = int(*(*uint32)(unsafe.Pointer(&slice[2])))
	}
	if key {
		slice = slice[2+khdrsize+vhdrsize : 2+khdrsize+vhdrsize+ksize]
		return *(*string)(unsafe.Pointer(&slice))
	}
	if vhdrsize == 0 {
		vsize = int(vhdr)
	} else if vhdrsize == 2 {
		vsize = int(*(*uint16)(unsafe.Pointer(&slice[2+khdrsize])))
	} else {
		vsize = int(*(*uint32)(unsafe.Pointer(&slice[2+khdrsize])))
	}
	slice = slice[2+khdrsize+vhdrsize+ksize : 2+khdrsize+vhdrsize+ksize+vsize]
	return *(*string)(unsafe.Pointer(&slice))
}

// Key returns the key portion of the key
func (pair Pair) Key() string {
	return pair.get(true)
}

// Value returns the value
func (pair Pair) Value() string {
	return pair.get(false)
}
