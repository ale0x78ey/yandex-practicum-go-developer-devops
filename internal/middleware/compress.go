package middleware

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type (
	GzipConfig struct {
		Level int
	}

	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
)

var (
	DefaultGzipConfig = GzipConfig{
		Level: gzip.BestSpeed,
	}
)

const (
	gzipScheme = "gzip"
)

func GzipEncoder() MiddlewareFunc {
	return GzipEncoderWithConfig(DefaultGzipConfig)
}

func GzipEncoderWithConfig(config GzipConfig) MiddlewareFunc {
	if config.Level == 0 {
		config.Level = DefaultGzipConfig.Level
	}

	pool := gzipCompressPool(config)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), gzipScheme) {
				next.ServeHTTP(w, r)
				return
			}

			i := pool.Get()
			gz, ok := i.(*gzip.Writer)
			if !ok {
				http.Error(w, i.(error).Error(), http.StatusInternalServerError)
			}
			gz.Reset(w)

			defer pool.Put(gz)
			defer gz.Close()

			w.Header().Set("Content-Encoding", gzipScheme)
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		}

		return http.HandlerFunc(fn)
	}
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipCompressPool(config GzipConfig) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			w, err := gzip.NewWriterLevel(ioutil.Discard, config.Level)
			if err != nil {
				return err
			}
			return w
		},
	}
}
