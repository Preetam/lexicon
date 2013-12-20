# lexicon [![Build Status](https://drone.io/github.com/PreetamJinka/lexicon/status.png)](https://drone.io/github.com/PreetamJinka/lexicon/latest)
An ordered map for Go.

A *lexicon* is a synonym for dictionary.

It's a combination of a hashmap and an [ordered list](https://github.com/PreetamJinka/orderedlist).

You'll have to pass in a compare function to New(). Here's an example that compares strings. Notice I'm recovering from panics.

```go
func compareStrings(a, b interface{}) (result int) {
	defer func() {
		if r := recover(); r != nil {
			// Log it?
		}
	}()

	aStr := a.(string)
	bStr := b.(string)

	if aStr > bStr {
		result = 1
	}

	if aStr < bStr {
		result -1
	}

	return
}
```

## Features
* Fast key-value existence checks and retrieval
* Range reads

## API Documentation

The GoDoc generated documentation is [here](http://godoc.org/github.com/PreetamJinka/lexicon).

## License
MIT
