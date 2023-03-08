package data

import (
	"fmt"
)

type Runtime int

func (r Runtime) MarshalJSON() ([]byte, error) {
	json := fmt.Sprintf(`"%d mins"`, r)

	return []byte(json), nil
}
