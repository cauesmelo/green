package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int

func (r Runtime) MarshalJSON() ([]byte, error) {
	json := fmt.Sprintf(`"%d mins"`, r)

	return []byte(json), nil
}

func (r *Runtime) UnmarshalJSON(json []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(json))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}
