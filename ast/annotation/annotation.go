package annotation

import (
	"go/ast"
	"strings"
)

type Annotations[N ast.Node] struct {
	Node N
	Ans  []Annotation
}
type Annotation struct {
	Key   string
	Value string
}

func NewAnnotations[N ast.Node](n N, group *ast.CommentGroup) Annotations[N] {
	if group == nil {
		return Annotations[N]{Node: n}
	}
	ans := make([]Annotation, 0, len(group.List))
	for i := 0; i < len(group.List); i++ {
		text, ok := extracContent(group.List[i])
		if ok {
			if strings.HasPrefix(text, "@") {
				regs := strings.SplitN(text[1:], " ", 2)
				if len(regs) != 2 {
					ans = append(ans, Annotation{Key: regs[0], Value: ""})
					continue
				}
				ans = append(ans, Annotation{Key: regs[0], Value: regs[1]})
			}
		} else {
			continue
		}
	}
	return Annotations[N]{
		Ans:  ans,
		Node: n,
	}
}

func extracContent(c *ast.Comment) (string, bool) {
	text := c.Text
	if strings.HasPrefix(text, "// ") {
		return text[3:], true
	} else if strings.HasPrefix(text, "/* ") {
		return text[3 : len(text)-2], true
	}
	return "", false
}
