package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/idna"
	"robpike.io/nihongo"
)

var (
	//go:embed index.html.tmpl
	html string
	//go:embed favicon.ico icon-128.png
	publicFiles embed.FS
)

var htmlTemplate = template.Must(template.New("index").Parse(html))

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHTTP)

	fs := http.FileServer(http.FS(publicFiles))
	mux.Handle("/favicon.ico", fs)
	mux.Handle("/icon-128.png", fs)

	log.Println("Serve on :8080")
	log.Fatal(http.ListenAndServe(":8080", logRequest(mux)))
}

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := struct {
			Date       string              `json:"date"`
			RemoteAddr string              `json:"remote_addr"`
			Method     string              `json:"method"`
			URI        string              `json:"uri"`
			Host       string              `json:"host"`
			Referrer   string              `json:"referrer"`
			UserAgent  string              `json:"user_agent"`
			Headers    map[string][]string `json:"headers"`
		}{
			Date:       time.Now().Format("2006-01-02 15:04:05"),
			RemoteAddr: r.RemoteAddr,
			Method:     r.Method,
			URI:        r.RequestURI,
			Host:       r.Host,
			Referrer:   r.Referer(),
			UserAgent:  r.UserAgent(),
			Headers:    r.Header,
		}
		b, err := json.Marshal(v)
		if err != nil {
			log.Printf("marshaling log: %v", err)
		} else {
			fmt.Println(string(b))
		}
		h.ServeHTTP(w, r)
	})
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	sub, domain := parseHost(r.Host)
	title := convertSubdomain(sub) + "なさそう"
	if path := strings.TrimPrefix(r.URL.Path, "/"); path != "" {
		title += nihongo.HiraganaString(path)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if title != "なさそう" {
		w.Header().Set("X-Robots-Tag", "noindex")
	}
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
