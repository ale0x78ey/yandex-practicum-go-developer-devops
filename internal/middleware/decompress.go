package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

func GzipDecoder() MiddlewareFunc {
	pool := gzipDecompressPool()

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Content-Encoding"), gzipScheme) {
				next.ServeHTTP(w, r)
				return
			}

			i := pool.Get()
			gz, ok := i.(*gzip.Reader)
			if !ok || gz == nil {
				http.Error(w, i.(error).Error(), http.StatusInternalServerError)
				return
			}
			defer pool.Put(gz)

			if err := gz.Reset(r.Body); err != nil {
				if err == io.EOF {
					next.ServeHTTP(w, r)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer gz.Close()
			r.Body = gz
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func gzipDecompressPool() sync.Pool {
	return sync.Pool{New: func() interface{} { return new(gzip.Reader) }}
}
