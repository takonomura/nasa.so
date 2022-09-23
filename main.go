package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/idna"
	"robpike.io/nihongo"
)

var (
	//go:embed index.html.tmpl
	html string
)

var htmlTemplate = template.Must(template.New("index").Parse(html))

func main() {
	log.Fatal(http.ListenAndServe("127.0.0.1:5000", http.HandlerFunc(handleHTTP)))
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	sub, domain := parseHost(r.Host)
	title := convertSubdomain(sub) + "なさそう"
	if path := strings.TrimPrefix(r.URL.Path, "/"); path != "" {
		title += nihongo.HiraganaString(path)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	htmlTemplate.Execute(w, struct {
		Title  string
		Domain string
	}{
		Title:  title,
		Domain: domain,
	})
}

func parseHost(host string) (sub, domain string) {
	host = strings.ToLower(host)
	if !strings.HasSuffix(host, ".nasa.so") {
		return "", "nasa.so"
	}
	return strings.TrimSuffix(host, ".nasa.so"), host
}

func convertSubdomain(sub string) string {
	if strings.HasPrefix(sub, "xn--") {
		s, _ := idna.ToUnicode(sub)
		return s
	}
	return nihongo.HiraganaString(sub)
}
