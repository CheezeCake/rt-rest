package web

import (
	"compress/gzip"
	"encoding/json"
	"net/http"

	"github.com/google/brotli/go/cbrotli"
)

type compressWriter interface {
	Write([]byte) (int, error)
	Close() error
}

func gzipCompressWriter(w http.ResponseWriter) compressWriter {
	w.Header().Set("Content-Encoding", "gzip")
	return gzip.NewWriter(w)
}

func brotliCompressWriter(w http.ResponseWriter) compressWriter {
	w.Header().Set("Content-Encoding", "br")
	return cbrotli.NewWriter(w, cbrotli.WriterOptions{11, 0})
}

func WriteJsonResponse(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	writer := w
	// writer := gzipCompressWriter(w)
	// writer := brotliCompressWriter(w)
	// defer writer.Close()
	return json.NewEncoder(writer).Encode(v)
}
