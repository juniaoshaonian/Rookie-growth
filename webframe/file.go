package webframe

import (
	lru "github.com/hashicorp/golang-lru"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileUploader struct {
	FileField   string
	Dst         string
	DstPathFunc func(fh *multipart.FileHeader)
}

func (f *FileUploader) Handle() HanleFunc {
	return func(c *Context) {
		file, header, err := c.Req.FormFile(f.FileField)
		if err != nil {
			c.ResponseCode = 500
			c.RespsonseDate = []byte("not found")
			return
		}
		dst, err := os.OpenFile(filepath.Join(f.DstPathFunc(header), f.Dst), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			c.ResponseCode = 500
			c.RespsonseDate = []byte("not found")
			return
		}
		io.CopyBuffer(dst, file, nil)
	}
}

type FileDownloader struct {
	dir string
}

func (f *FileDownloader) Handler(ctx *Context) {
	filename, err := ctx.QueryValue("File")
	if err != nil {
		ctx.ResponseCode = 500
		ctx.RespsonseDate = []byte("not found")
		return
	}
	path := filepath.Join(f.dir, filepath.Clean(filename))
	fn := filepath.Base(path)
	header := ctx.Resp.Header()
	header.Set("Content-Disposition", "attachment;filename="+fn)
	header.Set("Content-Description", "File Transfer")
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Transfer-Encoding", "binary")
	header.Set("Expires", "0")
	header.Set("Cache-Control", "must-revalidate")
	header.Set("Pragma", "public")
	http.ServeFile(ctx.Resp, ctx.Req, path)
}

type StaticResourceHandler struct {
	dir                     string
	extensionContentTypeMap map[string]string
	cache                   *lru.Cache
	maxFileSize             int
}

type StaticResourceHandlerOption func(handler *StaticResourceHandler)

func WithMoreExtention(extention map[string]string) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.extensionContentTypeMap = extention
	}
}

func NewStaticResourceHandelr(dir string, opts ...StaticResourceHandlerOption) *StaticResourceHandler {
	s := &StaticResourceHandler{
		dir: dir,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *StaticResourceHandler) Handle(ctx *Context) {
	filename, ok := ctx.pathParams["file"]
	if !ok {
		ctx.ResponseCode = 500
		ctx.RespsonseDate = []byte("not found")
		return
	}
	path := filepath.Join(s.dir, filepath.Clean(filename))
	fn := filepath.Base(path)
	val, ok := s.cache.Get(fn)
	header := ctx.Resp.Header()
	if ok {
		ctx.ResponseCode = 200
		ctx.RespsonseDate = val.(*cacheItem).data
		header.Set("Content-Type", val.(*cacheItem).contentType)
		return
	}
	f, err := os.Open(path)
	if err != nil {
		ctx.RespsonseDate = []byte("文件不存在")
		ctx.ResponseCode = 500
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		ctx.RespsonseDate = []byte("文件不存在")
		ctx.ResponseCode = 500
		return
	}
	newItem := &cacheItem{
		data:        data,
		contentType: s.extensionContentTypeMap[filepath.Ext(fn)],
	}

	if len(data) <= s.maxFileSize {
		s.cache.Add(fn, newItem)
	}
	header.Set("Content-Type", newItem.contentType)
	ctx.ResponseCode = 200
	ctx.RespsonseDate = data
}

type cacheItem struct {
	data        []byte
	contentType string
}
