package pair

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

const maxInt = int(^uint(0) >> 1)

// Pair is a tightly packed key/value pair
type Pair struct {
	ptr unsafe.Pointer
}

// New returns a Pair
func New(key, value []byte) Pair {
	slice := makenz(8 + len(value) + len(key))
	binary.LittleEndian.PutUint32(slice, uint32(len(value)))
	binary.LittleEndian.PutUint32(slice[4:], uint32(len(key)))
	copy(slice[8:], value)
	copy(slice[8+len(value):], key)
	return Pair{unsafe.Pointer(&slice[0])}
}
func (pair Pair) getSlice() []byte {
	if uintptr(pair.ptr) == 0 {
		return nil
	}
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(pair.ptr),
		Len:  maxInt,
		Cap:  maxInt,
	}))
}

// Value returns the value
func (pair Pair) Value() []byte {
	slice := pair.getSlice()
	if slice == nil {
		return nil
	}
	valuesz := binary.LittleEndian.Uint32(slice)
	return slice[8 : 8+valuesz : 8+valuesz]
}

// Key returns the key portion of the key
func (pair Pair) Key() []byte {
	slice := pair.getSlice()
	if slice == nil {
		return nil
	}
	valuesz := binary.LittleEndian.Uint32(slice)
	keysz := binary.LittleEndian.Uint32(slice[4:])
	return slice[8+valuesz : 8+valuesz+keysz : 8+valuesz+keysz]
}

// Size returns the size of the allocation
func (pair Pair) Size() int {
	slice := pair.getSlice()
	if slice == nil {
		return 0
	}
	valuesz := binary.LittleEndian.Uint32(slice)
	keysz := binary.LittleEndian.Uint32(slice[4:])
	return int(8 + valuesz + keysz)
}

// Zero return true if the pair is unallocated
func (pair Pair) Zero() bool {
	return uintptr(pair.ptr) == 0
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

// Pointer returns the underlying pointer
func (pair Pair) Pointer() unsafe.Pointer {
	return pair.ptr
}

// FromPointer returns a pair that uses the memory at the pointer.
func FromPointer(ptr unsafe.Pointer) Pair {
	return Pair{ptr}
}
