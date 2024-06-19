package server

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const WSCODE = `
<script>
    let socket = new WebSocket("ws://127.0.0.1:8080/ws");
    socket.onmessage = (msg) => {
        if (msg.data === "RELOAD") {
            location.reload();
        }
    }
</script>
`

func index(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path

	if isHtmlPage(urlPath) {
		name := filepath.Join(htmlDir, getHtmlPageName(urlPath))
		b, err := os.ReadFile(name)
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		b = insertWsCode(b)
		w.Write(b)
		return
	}

	fullPath := filepath.Join(htmlDir, urlPath)
	http.ServeFile(w, r, fullPath)
}

func isHtmlPage(urlPath string) bool {
	last := urlPath[len(urlPath)-1:]
	ext := path.Ext(urlPath)
	return last == "/" || ext == ".html"
}

func getHtmlPageName(urlPath string) string {
	last := urlPath[len(urlPath)-1:]
	if last == "/" {
		return urlPath + "index.html"
	}
	return urlPath
}

func insertWsCode(b []byte) []byte {
	s := string(b)
	idx := strings.LastIndex(s, "</body>")
	if idx > 0 {
		var bb bytes.Buffer
		bb.Write([]byte(s[:idx]))
		bb.Write([]byte(WSCODE))
		bb.Write([]byte(s[idx:]))
		b = bb.Bytes()
	}
	return b
}
