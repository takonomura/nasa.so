package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/idna"
	"robpike.io/nihongo"
)

var indexTemplate = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html lang="ja">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{ .Title }}</title>

	<link rel="icon" sizes="128x128" href="icon-128.png">

	<meta property="og:title" content="{{ .Title }}">
	<meta property="og:type" content="website">
	<meta property="og:url" content="https://{{ .Domain }}/">
	<meta property="og:image" content="https://{{ .Domain }}/icon-128.png">

	<style>
		html, body {
			width: 100%;
			height: 100%;
		}

		.centered {
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100%;
		}

		h1 {
			font-size: 6rem;
		}
	</style>
</head>
<body>
	<div class="centered">
		<h1>{{ .Title }}</h1>
	</div>
</body>
</html>
`))

func main() {
	log.Fatal(http.ListenAndServe("127.0.0.1:5000", http.HandlerFunc(handleHTTP)))
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	sub, domain := parseHost(r.Host)
	title := convertSubdomain(sub) + "なさそう"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	indexTemplate.Execute(w, struct {
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
