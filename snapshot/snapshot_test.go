package snapshot

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshotRead(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "name",
			input:  "input",
			expect: "expect",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filename := "/tmp/snapshot-112974791.bin"
			r, err := NewReader(filename)
			fmt.Println("Filename", filename)
			defer r.Close()

			assert.NoError(t, err)
			assert.Equal(t, r.Header.Version, uint32(1))

			for {
				section, err := r.Next()
				if err == io.EOF {
					break
				}
				assert.NoError(t, err)
				fmt.Println("Section", section.Name, "rows", section.RowCount, "bytes", section.BufferSize)
			}
			fmt.Println("Done")
		})
	}
}
