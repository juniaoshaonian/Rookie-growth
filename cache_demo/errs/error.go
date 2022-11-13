package errs

import (
	"errors"
	"fmt"
)

func NewErrKeyNotFound(key string) error {
	return fmt.Errorf("cache: 找不到 key %s", key)

}

var ErrKeyNotFound error = errors.New("cache: key not found")
