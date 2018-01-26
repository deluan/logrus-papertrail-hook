package logrus_papertrail

import (
	"io"
)

type bufwriter struct {
	buffer chan []byte
	w      io.Writer
}

func (bw bufwriter) Write(p []byte) (int, error) {
	if bw.buffer == nil {
		return bw.w.Write(p)
	}
	bw.buffer <- p
	return len(p), nil
}

func newBufwriter(writer io.Writer, n int) bufwriter {
	b := bufwriter{w: writer}
	if n > 0 {
		b.buffer = make(chan []byte, n)
		go func() {
			for p := range b.buffer {
				b.w.Write(p)
			}
		}()
	}
	return b
}
