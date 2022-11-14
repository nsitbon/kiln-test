package reader

import (
	"fmt"
	"io"
)

func Summarize(body io.Reader, limit uint) string {
	buf := make([]byte, limit)

	if _, err := body.Read(buf); err != nil && err != io.EOF {
		return fmt.Sprintf("failed to summarize io.Reader: %s", err.Error())
	}

	return string(buf)
}

func DrainAndClose(rc io.ReadCloser) {
	if rc != nil {
		_, _ = io.Copy(io.Discard, rc)
		_ = rc.Close()
	}
}
