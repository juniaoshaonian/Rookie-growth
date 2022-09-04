package http

import (
	"io"
	"text/template"
)

// 这部分和课堂的很像，但是有一些地方被我改掉了
const serviceTpl = `
{{- $service := .GenName  -}}
package {{ .Package}}

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type {{ $service }} struct {
    Endpoint string
    Path string
	Client http.Client
}
{{range $idx,$method := .Methods}}
func (s *{{$service}}) {{$method.Name}}(ctx context.Context, req *{{$method.ReqTypeName}}) (*{{$method.RespTypeName}}, error) {
	url := s.Endpoint + s.Path + "/{{$method.Name}}"
	bs, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body := &bytes.Buffer{}
	body.Write(bs)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	httpResp, err := s.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	bs, err = ioutil.ReadAll(httpResp.Body)
	resp := &{{$method.RespTypeName}}{}
	err = json.Unmarshal(bs, resp)
	return resp, err
}
{{end}}
`

func Gen(writer io.Writer, def ServiceDefinition) error {
	tpl := template.New("service")
	tpl, err := tpl.Parse(serviceTpl)
	if err != nil {
		return err
	}
	// 还可以进一步调用 format.Source 来格式化生成代码
	return tpl.Execute(writer, def)
}
