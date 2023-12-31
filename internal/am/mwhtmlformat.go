package am

import (
	"bytes"
	"net/http"
	"strings"
)

func HTMLFormat(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rr := &responseRecorder{ResponseWriter: w, buf: bytes.NewBuffer(nil)}

		next.ServeHTTP(rr, r)

		htmlContent := rr.buf.String()
		formattedHTML := formatHTML(htmlContent)
		w.Write([]byte(formattedHTML))
	})
}

func formatHTML(html string) string {
	html = strings.ReplaceAll(html, "\n", "")
	html = strings.ReplaceAll(html, "  ", "")

	return html
}

type responseRecorder struct {
	http.ResponseWriter
	buf *bytes.Buffer
}
