Pair
====
<a href="https://godoc.org/github.com/tidwall/pair"><img src="https://img.shields.io/badge/api-reference-blue.svg?style=flat-square" alt="GoDoc"></a>

Pair is a Go package that provides a low memory key/value object that takes up one allocation. It's useful for in-memory key/value stores and data structures where memory space is a concern.


Data structure
--------------

The allocation is a single packed block of bytes with the following format:

| ValueSize uint32 | KeySize uint32 | Value []byte | Key []byte |
|------------------|----------------|--------------|------------|

Using
-----

To start using Pair, install Go and run `go get`:

```
$ go get -u github.com/tidwall/pair
```

Create a new Pair:

```go
item := pair.New([]byte("user:2054:name"), []byte("Alice Tripplehorn"))
```

Access the Pair data:
```go
item.Key() []byte    // returns the key portion of the pair.
item.Value() []byte  // returns the value portion of the pair.
item.Size() int      // returns the exact in-memory size of the item.
item.Zero() bool     // returns true if the pair is unallocated.
```

Unsafe pointer access:
```go
item.Pointer() unsafe.Pointer             // returns the base pointer
pair.FromPointer(ptr unsafe.Pointer) Pair // returns a Pair with provided base pointer
```

Contact
-------
Josh Baker [@tidwall](http://twitter.com/tidwall)

License
-------
Pair source code is available under the MIT [License](/LICENSE).

