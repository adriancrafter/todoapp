package am

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	ListCtxKey = "list"
	ReqCtxKey  = "req"
)

const (
	LogEntryCtxKey = "LogEntry"
)

const (
	tsFormat = "2006/01/02 15:04:05"
)

type (
	ReqLogger struct {
		log Logger
	}
)

func (rl *ReqLogger) Log() Logger {
	return rl.log
}

func NewReqLogger(log Logger) *ReqLogger {
	return &ReqLogger{log: log}
}

func NewReqLoggerMiddleware(log Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		rl := NewReqLogger(log)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entry := rl.NewLogEntry(r)
			ww := NewWrapResponseWriter(w)

			t1 := time.Now()
			defer func() {
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
			}()

			// Store log entry into context
			next.ServeHTTP(ww, r.WithContext(context.WithValue(r.Context(), LogEntryCtxKey, entry)))
		})
	}
}

func (rl *ReqLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	fields := map[string]string{}

	reqID := r.Header.Get("X-Request-Id")
	if reqID != "" {
		fields["req-id"] = reqID
	}

	scheme := "httpx"
	if r.TLS != nil {
		scheme = "https"
	}

	fields["scheme"] = scheme
	fields["proto"] = r.Proto
	fields["method"] = r.Method
	fields["addr"] = r.RemoteAddr
	fields["agent"] = r.UserAgent()
	fields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	fields["ts"] = time.Now().UTC().Format(tsFormat)

	sb := strings.Builder{}
	for k, v := range fields {
		sb.WriteString(fmt.Sprintf("%s: %s, ", k, v))
	}

	return &LogEntry{
		log:   rl.Log(),
		entry: &sb,
	}
}

type (
	LogEntry struct {
		log   Logger
		entry *strings.Builder
	}
)

func (le *LogEntry) Log() Logger {
	return le.log
}

func (le *LogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	le.entry.WriteString(fmt.Sprintf("%s: %d, ", "status", status))
	le.entry.WriteString(fmt.Sprintf("%s: %d, ", "bytes", bytes))
	le.entry.WriteString(fmt.Sprintf("%s: %fms", "elapsed", float64(elapsed.Nanoseconds())/1000000.0))
	le.Log().Debugf("%s", le.entry.String())
}

func (le *LogEntry) Panic(v interface{}, stack []byte) {
	le.entry.WriteString(fmt.Sprintf("%s: %s, ", "stack", string(stack)))
	le.entry.WriteString(fmt.Sprintf("%s: %s, ", "panic", fmt.Sprintf("%+v", v)))
	le.Log().Debugf("%s", le.entry.String())
}

type WrapResponseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func NewWrapResponseWriter(w http.ResponseWriter) *WrapResponseWriter {
	return &WrapResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (ww *WrapResponseWriter) WriteHeader(code int) {
	ww.statusCode = code
	ww.ResponseWriter.WriteHeader(code)
}

func (ww *WrapResponseWriter) Write(b []byte) (int, error) {
	bytesWritten, err := ww.ResponseWriter.Write(b)
	ww.bytes += bytesWritten
	return bytesWritten, err
}

func (ww *WrapResponseWriter) Status() int {
	return ww.statusCode
}

func (ww *WrapResponseWriter) BytesWritten() int {
	return ww.bytes
}
