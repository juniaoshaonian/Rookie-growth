package errs

import (
	"errors"
	"fmt"
)

func NewErrPathCoverage(path string) error {
	return errors.New(fmt.Sprintf("path: %v Coverage", path))
}

func NewErrNodeConflict(node1 string, node2 string) error {
	return errors.New(fmt.Sprintf("已有%v 节点，不能同时注册%v节点", node1, node2))
}
