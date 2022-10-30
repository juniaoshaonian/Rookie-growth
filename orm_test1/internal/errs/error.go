package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly   = errors.New("orm: 只支持一级指针作为输入，例如 *User")
	ErrNoRows        = errors.New("orm: 未找到数据")
	ErrInsertZeroRow = errors.New("orm: 插入0行")
)

func NewErrField(name string) error {
	return errors.New(fmt.Sprintf("orm: %v字段不存在", name))
}

func NewErrUnknownColumn(name string) error {
	return errors.New(fmt.Sprintf("orm %v未知列", name))
}
func NewErrUnsupportedSelectable(exp any) error {
	return fmt.Errorf("orm: 不支持的目标列 %v", exp)
}
