Pair
====
<a href="https://godoc.org/github.com/tidwall/pair"><img src="https://img.shields.io/badge/api-reference-blue.svg?style=flat-square" alt="GoDoc"></a>

Pair is a Go package that provides a low memory key/value object that takes up one allocation.

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

Contact
-------
Josh Baker [@tidwall](http://twitter.com/tidwall)

License
-------
Pair source code is available under the MIT [License](/LICENSE).

