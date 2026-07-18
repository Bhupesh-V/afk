package pkg

import "io"

type NoBellWriter struct {
	io.Writer
}

func (w *NoBellWriter) Write(p []byte) (n int, err error) {
	// Filter out the ASCII Bell character (\a or 0x07)
	cleaned := make([]byte, 0, len(p))
	for _, b := range p {
		if b != '\a' {
			cleaned = append(cleaned, b)
		}
	}
	return w.Writer.Write(cleaned)
}

// Close implements io.WriteCloser since promptui requires it
func (w *NoBellWriter) Close() error {
	return nil
}
