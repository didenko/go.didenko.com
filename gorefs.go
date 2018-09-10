package gorefs

import (
	"html/template"
	"net/http"
	"regexp"
)

var (
	err     error
	tpl     *template.Template
	pathRE  *regexp.Regexp
	allowed map[string]string
)

func init() {
	pathRE = regexp.MustCompile(`^\/([[:alnum:]]*)`)
	allowed = map[string]string{
		"fst":   "https://github.com/didenko/fst",
		"pald":  "https://github.com/didenko/pald",
		"slops": "https://github.com/didenko/slops",
		"gate":  "https://gitlab.com/vldid/gate",
	}

	tpl, err = template.New("goref").Parse(`
<html><head>
	<meta name="go-import" content="go.didenko.com/{{.P}} git {{.D}}">
	<meta name="go-source" content="go.didenko.com/{{.P}} git {{.D}} {{.D}}/tree/master{/dir} {{.D}}/blob/master{/dir}/{file}#L{line}">
	<meta http-equiv="refresh" content="0; url=http://www.didenko.com">
</head></html>`)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
}

func wrongPath(w http.ResponseWriter, path string) {
	w.WriteHeader(404)
	w.Write([]byte("Wrong path: " + path))
}

func handler(w http.ResponseWriter, r *http.Request) {

	parsedPath := pathRE.FindStringSubmatch(r.URL.Path)

	if len(parsedPath) < 1 {
		wrongPath(w, r.URL.Path)
		return
	}

	projectName := parsedPath[1]

	for project, dest := range allowed {
		if projectName == project {
			tpl.Execute(w, struct{ P, D string }{project, dest})
			return
		}
	}
	wrongPath(w, r.URL.Path)
}
