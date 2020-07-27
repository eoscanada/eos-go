package snapshot

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSnapshotRead(t *testing.T) {

	//if os.Getenv("READ_SNAPSHOT_FILE") != "true" {
	//	t.Skipf("Environment varaible 'READ_SNAPSHOT_FILE' not set to true")
	//	return
	//}

	logger, _ := zap.NewDevelopment()
	tests := []struct {
		name     string
		testFile string
	}{
		{name: "eos-local dev", testFile: "eos-jdev_0000000638.bin"},
		{name: "eos-dev1", testFile: "eos-dev1_0004841949.bin"},
		{name: "Battlefield - b8d703ed1", testFile: "battlefield-snapshot.bin"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFile := testData(test.testFile)

			if !fileExists(testFile) {
				logger.Error("test file not found", zap.String("testfile", testFile))
				return
			}

			r, err := NewReader(testFile)
			require.NoError(t, err)
			defer r.Close()

			assert.NoError(t, err)
			assert.Equal(t, r.Header.Version, uint32(1))

			for {
				section, err := r.Next()
				if err == io.EOF {
					break
				}
				assert.NoError(t, err)
				logger.Info("new section",
					zap.String("section_name", string(section.Name)),
					zap.Uint64("row_count", section.RowCount),
					zap.Uint64("bytes_count", section.BufferSize),
					zap.Uint64("bytes_count", section.Offset),
				)
				switch section.Name {
				case SectionNameAccountObject:
					require.NoError(t, section.Process(func(o interface{}) error {
						acc, ok := o.(AccountObject)
						if !ok {
							return fmt.Errorf("process account object: unexpected object type: %T", o)
						}
						logger.Info("new account object",
							zap.String("name", string(acc.Name)),
						)
						return nil
					}))
				case SectionNameContractTables:
					var currentTable *TableIDObject
					require.NoError(t, section.Process(func(o interface{}) error {
						switch obj := o.(type) {
						case *TableIDObject:
							logger.Info("Table ID", zap.Reflect("table_id", obj))
							currentTable = obj
						case *KeyValueObject:
							if currentTable.Count < 20 {
								logger.Info("Key Value Object", zap.Reflect("kv", obj))
							}
						case *Index64Object:
							logger.Info("Index64Object", zap.Reflect("index_64_object", obj))
						case *Index128Object:
							logger.Info("Index128Object", zap.Reflect("index_128_object", obj))
						case *Index256Object:
							logger.Info("Index256Object", zap.Reflect("index_256_object", obj))
						case *IndexDoubleObject:
							logger.Info("IndexDoubleObject", zap.Reflect("index_double_object", obj))
						case *IndexLongDoubleObject:
							logger.Info("IndexLongDoubleObject", zap.Reflect("index_long_object", obj))
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

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if err != nil {
		return false
	}

	return !info.IsDir()
}

func testData(filename string) string {
	return filepath.Join("test-data", filename)
}
