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
	allowed []string
)

func init() {
	pathRE = regexp.MustCompile(`^\/([[:alnum:]]*)`)
	allowed = []string{
		"slops",
		"fstests",
	}

	tpl, err = template.New("goref").Parse(`<html><head><meta name="go-import" content="github.com/didenko/{{.}}"></head></html>`)
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

	permitted := false
	for _, a := range allowed {
		if projectName == a {
			permitted = true
			break
		}
	}

	if permitted {
		tpl.Execute(w, projectName)
	} else {
		wrongPath(w, r.URL.Path)
	}
}
