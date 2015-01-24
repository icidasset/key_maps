package api

import (
  "compress/gzip"
  "github.com/gocraft/web"
  "io"
  "net/http"
  "strings"
)


type gzipResponseWriter struct {
  io.Writer
  web.ResponseWriter
}


func (c *BaseContext) Gzip(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
  rw.Header().Set("Vary", "Accept-Encoding")

  if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
    next(rw, req)
    return
  }

  // headers
  rw.Header().Set("Content-Encoding", "gzip")

  // gzip writer
  gz := gzip.NewWriter(rw)
  defer gz.Close()

  // gzip response writer
  gzw := gzipResponseWriter{ Writer: gz, ResponseWriter: rw }

  // next
  next(gzw, req)
}


func (grw gzipResponseWriter) Write(p []byte) (int, error) {
  if len(grw.Header().Get("Content-Type")) == 0 {
    grw.Header().Set("Content-Type", http.DetectContentType(p))
  }

  return grw.Writer.Write(p)
}
