package boot

import (
	"encoding/json"

	yaml2json "github.com/bronze1man/go-yaml2json"
)

func yamlUnmarshal(cnt []byte, v interface{}) error {
	jsonCnt, err := yaml2json.Convert(cnt)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonCnt, v)
}
