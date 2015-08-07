package engine

import "io"

// OutputTooLarge is the error type returned when write capacty is exceeded.
type OutputTooLarge struct {
}

func (err *OutputTooLarge) Error() string {
	return "write capacity exceeded"
}

// LimitedWriter wraps a bytes.Buffer and limits the # of bytes you can write.
// Going over the limit will cause subsequent Write calls to return an error.
type LimitedWriter struct {
	io.Writer
	available int
}

// NewLimitedWriter returns a new LimitedWriter instance with the specified
// write limit.
func NewLimitedWriter(writer io.Writer, limit int) *LimitedWriter {
	return &LimitedWriter{writer, limit}
}

func (buf *LimitedWriter) Write(p []byte) (n int, err error) {
	shortWrite := false
	if len(p) > buf.available {
		shortWrite = true
		p = p[0:buf.available]
	}

	sz, err := buf.Writer.Write(p)
	buf.available = buf.available - sz

	if shortWrite == true && err == nil {
		err = &OutputTooLarge{}
	}
	return sz, err
}
