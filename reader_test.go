package limit_reader

import (
	"errors"
	"io"
	"strings"
	"testing"
)

var text = "this is a string that is 49 characters in length."
var textLen = int64(len(text))

func TestLimitedReader_ReadAll_LongBuffer(t *testing.T) {
	doReadAllTest(t, textLen*2)
}

func TestLimitedReader_Read_LongBuffer(t *testing.T) {
	doReadTest(t, textLen*2)
}

func TestLimitedReader_ReadAll_LongBufferBy1(t *testing.T) {
	doReadAllTest(t, textLen+1)
}

func TestLimitedReader_Read_LongBufferBy1(t *testing.T) {
	doReadTest(t, textLen+1)
}

func TestLimitedReader_ReadAll_ShortBuffer(t *testing.T) {
	doReadAllTest(t, textLen/2)
}

func TestLimitedReader_Read_ShortBuffer(t *testing.T) {
	doReadTest(t, textLen/2)
}

func TestLimitedReader_ReadAll_ShortBufferBy1(t *testing.T) {
	doReadAllTest(t, textLen-1)
}

func TestLimitedReader_Read_ShortBufferBy1(t *testing.T) {
	doReadTest(t, textLen-1)
}

func TestLimitedReader_ReadAll_ExactBuffer(t *testing.T) {
	doReadAllTest(t, textLen)
}

func TestLimitedReader_Read_ExactBuffer(t *testing.T) {
	doReadTest(t, textLen)
}

func doReadTest(t *testing.T, readLimit int64) {
	r := strings.NewReader(text)
	lr := New(r, readLimit)
	rbuf := make([]byte, 10)
	resultBuf := make([]byte, 0)
	totalBytesRead := 0
	var err error
	var n int
	for err == nil {
		n, err = lr.Read(rbuf)
		totalBytesRead += n
		resultBuf = append(resultBuf, rbuf[:n]...)
	}

	if readLimit >= int64(len(text)) {
		if int64(len(resultBuf)) != textLen {
			t.Errorf("expected %d bytes read, got: %d", textLen, len(resultBuf))
		}
		compareBuffers(t, text, string(resultBuf))
		if !errors.Is(err, io.EOF) {
			t.Errorf("expected no error, got: %v", err)
		}
	} else {
		if int64(len(resultBuf)) != readLimit {
			t.Errorf("expected %d bytes read, got: %d", readLimit, len(resultBuf))
		}
		compareBuffers(t, text[:readLimit], string(resultBuf))

		if err == nil {
			t.Errorf("expected error, received nil")
		} else if !errors.Is(err, ReaderBoundsExceededError{}) {
			t.Errorf("expected a ReaderBoundsExceededError, received: %T", err)
		}
	}
}

func doReadAllTest(t *testing.T, readLimit int64) {
	r := strings.NewReader(text)
	lr := New(r, readLimit)
	resultBuf, err := io.ReadAll(lr)

	if readLimit >= int64(len(text)) {
		if int64(len(resultBuf)) != textLen {
			t.Errorf("expected %d bytes read, got: %d", textLen, len(resultBuf))
		}
		compareBuffers(t, text, string(resultBuf))
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	} else {
		if int64(len(resultBuf)) != readLimit {
			t.Errorf("expected %d bytes read, got: %d", readLimit, len(resultBuf))
		}
		compareBuffers(t, text[:readLimit], string(resultBuf))

		if err == nil {
			t.Errorf("expected error, received nil")
		} else if !errors.Is(err, ReaderBoundsExceededError{}) {
			t.Errorf("expected a ReaderBoundsExceededError, received: %T", err)
		}
	}
}

func compareBuffers(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("buffers do not match, expected: \"%s\", received: \"%s\"", expected, actual)
	}
}
