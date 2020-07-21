package snapshot

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			filename := "/tmp/0125111385-07750c59b24ed52d2dbf2048b67b58e9c9bd53ff5cc4550277718c1d5d800f73-snapshot.bin" // mainnet
			//filename := "/tmp/0003212331-0031042b02b2cf711fee6e1e24da94101fa6c1ea9ece568d5f13232473429db1-snapshot.bin" // kylin
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
				fmt.Println("Section", section.Name, "rows", section.RowCount, "bytes", section.BufferSize, "offset", section.Offset)

				if strings.Contains(section.Name, "contract") {
					require.NoError(t, section.Process(func(o interface{}) error {
						switch obj := o.(type) {
						case *TableIDObject:
							fmt.Println("Table ID", obj.Code, obj.Scope, obj.TableName)
						case *KeyValueObject:
							fmt.Println("KV", obj.PrimKey, obj.Value)
						default:
							fmt.Printf("Ignoring row %T\n", obj)
						}
						return nil

					}))
				}
			}
		})
	}
}
