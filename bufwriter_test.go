package logrus_papertrail

import (
	"strings"
	"testing"
	"time"
)

type testWriter struct {
	buffer []byte
}

func (t *testWriter) Write(b []byte) (int, error) {
	t.buffer = append(t.buffer, b...)
	return len(b), nil
}

func TestBufwriterBypass(t *testing.T) {
	double := &testWriter{}
	testString := "this is a test"
	bw := newBufwriter(double, 0)
	bw.Write([]byte(testString))

	if len(double.buffer) == 0 {
		t.Error("Nothing was received!")
	}
	if string(double.buffer) != testString {
		t.Errorf("Corrupted data: '%s'", string(double.buffer))
	}
}

func TestBufwriterAsync(t *testing.T) {
	double := &testWriter{}
	testString := "this is a test"
	bw := newBufwriter(double, 2)

	expectedResult := strings.Repeat(testString, 10000)
	for i := 0; i < 10000; i++ {
		bw.Write([]byte(testString))
	}

	// Wait for channel to drain
	for {
		if len(bw.buffer) == 0 {
			break
		}
		<-time.After(time.Duration(1 * time.Millisecond))
	}

	if len(double.buffer) == 0 {
		t.Error("Nothing was received!")
	}
	if string(double.buffer) != expectedResult {
		t.Errorf("Corrupted data: '%s'", string(double.buffer))
	}
}
