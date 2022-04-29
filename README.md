# Golang Limit Reader

This was inspired from the Golang builtin `io.LimitReader`. See the Go source 
[here](https://cs.opensource.google/go/go/+/refs/tags/go1.18.1:src/io/io.go;l=455;drc=refs%2Ftags%2Fgo1.18.1)

This is a small utility wrapper for a `Reader` interface in Golang that will limit the number of bytes attempted to 
be read from the supplied `Reader`. If the number of bytes is greater than the supplied limit a 
`ReaderBoundsExceededError` will be returned. The original simply returned an `EOF` when the number of bytes
read exceeded the limit. This allows the calling code to determine if the number of bytes read was within
the bounds specified or if the read was truncated.

This allows for more deterministic code and robust error handling.

## Example Usage

```
r := strings.NewReader(someText)
lr := New(r, 512)
buf, err := io.ReadAll(lr)
if err != nil {
	// Handle error case here
	// io.ReadAll will never return an EOF as an error so we know in this example
	// that the original Reader had more bytes than we are willing to process
} else {
	// Handle the success case here
}
```

```
r := strings.NewReader(someText)
lr := New(r, 512)
buf := make([]byte, 512)
nu,, err := lr.Read(buf)
if err != nil {
	if _, ok := err.(ReaderBoundsExceededError); ok {
		// Here we need to see if the error is a ReaderBoundsExceededError to determine the
		// Original Reader had more bytes then we are willing to process
		
		// Handle too much data error
	} else { 
		// Handle other cases (it may be an EOF, in which case we were able to read
		// all the bytes from the original Reader)
    }
} else {
	// Handle the success case here
}
```



 