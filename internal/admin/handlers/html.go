package handlers

import (
	"bytes"
	"embed"
	"io"
	"net/http"
	"path"
	"text/template"
)

var defaultFuncs = template.FuncMap{
	"defTitle": func(ip interface{}) string {
		v, _ := ip.(string)
		if v == "" {
			return "Dagu Admin"
		}
		return v
	},
}

//go:embed web/templates/*
var assets embed.FS
var templatePath = "web/templates/"

func useTemplate(layout string, name string) func(http.ResponseWriter, interface{}) error {
	files := append(baseTemplates(), path.Join(templatePath, layout))
	tmpl, err := template.New(name).Funcs(defaultFuncs).ParseFS(assets, files...)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, data interface{}) error {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
			//log.Printf("ERR: %v\n", err)
			//http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
		w.WriteHeader(http.StatusOK)
		_, _ = io.Copy(w, &buf)
		return nil
	}
}

func baseTemplates() []string {
	var templateFiles = []string{
		"base.gohtml",
	}
	var ret []string
	for _, t := range templateFiles {
		ret = append(ret, path.Join(templatePath, t))
	}
	return ret
}
