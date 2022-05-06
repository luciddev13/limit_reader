package limit_reader

import "io"

// A limitedReader reads from R but limits the amount of
// data returned to just N bytes. Each call to Read
// updates N to reflect the new amount remaining.
// Read returns ReaderBoundsExceededError when N <= 0 and there is
// still more bytes to read from R or EOF when the underlying R returns EOF.
type limitedReader struct {
	r        io.Reader
	n        int64
	complete error
}

// New returns a Reader that reads from r
// but stops with ReaderBoundsExceededError after n bytes.
// The underlying implementation is a *limitedReader.
func New(r io.Reader, n int64) io.Reader {
	return &limitedReader{r, n, nil}
}

func (l *limitedReader) Read(p []byte) (int, error) {
	if l.complete != nil {
		return 0, l.complete
	}

	if int64(len(p)) > l.n {
		p = p[0:l.n]
	}
	n, err := l.r.Read(p)
	l.n -= int64(n)

	if err != nil {
		l.complete = err
	} else {
		if len(p) == 0 {
			err = ReaderBoundsExceededError{}
		}
	}
	return n, err
}
