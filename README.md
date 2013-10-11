# lexicon [![Build Status](https://drone.io/github.com/PreetamJinka/lexicon/status.png)](https://drone.io/github.com/PreetamJinka/lexicon/latest)
A lexicographically ordered map for Go.

A *lexicon* is a synonym for dictionary.

It's a combination of a hashmap and an [ordered list](https://github.com/PreetamJinka/orderedlist).

I've tried to make it as generic as possible, but since we need to have some way of knowing how to order elements,
keys must implement `orderedlist.Comparable`.

```go
type Comparable interface {
	Compare(c Comparable) int
}
```

## Features
* Fast key-value existence checks and retrieval
* Range reads

## API Documentation

The GoDoc generated documentation is [here](http://godoc.org/github.com/PreetamJinka/lexicon).

## License
MIT
