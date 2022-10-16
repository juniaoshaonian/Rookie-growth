package orm

import (
	"orm/internal/model"
	"strings"
)

type builder struct {
	model   *model.Model
	sb      strings.Builder
	args    []any
	dialect Dialect
}
