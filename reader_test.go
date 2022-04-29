package limit_reader

import (
	"io"
	"strings"
	"testing"
)

var text = "this is a string that is 49 characters in length."
var textLen = int64(len(text))

func TestLimitedReadAllReader_Read_Short(t *testing.T) {
	r := strings.NewReader(text)
	lr := New(r, textLen*2)
	buf, err := io.ReadAll(lr)
	if err != nil {
		t.Errorf("expected no error, received: %v", err)
		return
	}
	compareBuffers(t, text, string(buf))
}

func TestLimitedReadAllReader_Read_ShortBy1(t *testing.T) {
	r := strings.NewReader(text)
	lr := New(r, textLen+1)
	buf, err := io.ReadAll(lr)
	if err != nil {
		t.Errorf("expected no error, received: %v", err)
		return
	}
	compareBuffers(t, text, string(buf))
}

func TestLimitedReadAllReader_Read_Long(t *testing.T) {
	r := strings.NewReader(text)
	lr := New(r, textLen/2)
	_, err := io.ReadAll(lr)
	if err == nil {
		t.Errorf("expected error, received nil")
		return
	}
	if _, ok := err.(ReaderBoundsExceededError); !ok {
		t.Errorf("expected a ReaderBoundsExceededError, received: %T", err)
	}
}

func TestLimitedReadAllReader_Read_LongBy1(t *testing.T) {
	r := strings.NewReader(text)
	lr := New(r, textLen-1)
	_, err := io.ReadAll(lr)
	if err == nil {
		t.Errorf("expected error, received nil")
		return
	}
	if _, ok := err.(ReaderBoundsExceededError); !ok {
		t.Errorf("expected a ReaderBoundsExceededError, received: %T", err)
	}
}

func TestLimitedReadAllReader_Read_Exact(t *testing.T) {
	r := strings.NewReader(text)
	lr := New(r, textLen)
	buf, err := io.ReadAll(lr)
	if err != nil {
		t.Errorf("expected no error, received: %v", err)
		return
	}
	compareBuffers(t, text, string(buf))
}

func compareBuffers(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("buffers do not match, expected: %s, received: %s", expected, actual)
	}
}
