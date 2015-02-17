package wfilter

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/getlantern/testify/assert"
)

func TestLines(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	i := int32(0)
	w := Lines(buf, func(w io.Writer, line string) (int, error) {
		j := atomic.AddInt32(&i, 1)
		if !strings.HasPrefix(line, "C") {
			return fmt.Fprintf(w, "%d %s", j, line)
		}
		return 0, nil
	})

	fmt.Fprintln(w, "A")
	fmt.Fprintln(w, "B")
	fmt.Fprintln(w, "C")
	fmt.Fprintln(w, "D")

	assert.Equal(t, expected, string(buf.Bytes()))
}

func TestLongLine(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := Lines(buf, func(w io.Writer, line string) (int, error) {
		return w.Write([]byte(line))
	})

	for i := 0; i <= MaxLineLength/10; i++ {
		w.Write([]byte("1234567890"))
	}

	fmt.Fprintln(w, "An actual line")
	assert.Equal(t, "An actual line\n", string(buf.Bytes()))
}

var expected = `1 A
2 B
4 D
`
