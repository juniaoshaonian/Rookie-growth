package webframe_2

import (
	"fmt"
	"strings"
	"testing"
)

func TestSTRINGS(t *testing.T) {
	path := "user/home"
	segs := strings.Split(path, "/")
	fmt.Println(segs)
}
