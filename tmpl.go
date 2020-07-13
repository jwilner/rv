package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

const baseTemplate = "base.html"

func loadTmplMgr(tmplDir string) (*tmplMgr, error) {
	ents, err := ioutil.ReadDir(tmplDir)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadDir %q: %w", tmplDir, err)
	}
	t := &tmplMgr{tmplDir: tmplDir, templates: make(map[string]*template.Template)}
	for _, e := range ents {
		if e.Name() == baseTemplate || e.IsDir() {
			continue
		}
		if t.templates[e.Name()], err = loadTmpl(t.tmplDir, e.Name()); err != nil {
			return nil, fmt.Errorf("loadTmpl %v %v: %w", t.tmplDir, e.Name(), err)
		}
	}
	return t, nil
}

type tmplMgr struct {
	tmplDir   string
	templates map[string]*template.Template
}

func (t *tmplMgr) render(ctx context.Context, w http.ResponseWriter, name string, val interface{}) {
	tmpl := t.templates[name]
	if tmpl == nil {
		log.Printf("unable to find template: %q\n", name)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if isDebug(ctx) {
		log.Printf("reloading template %q\n", name)
		var err error
		if tmpl, err = loadTmpl(t.tmplDir, name); err != nil {
			log.Printf("unable to parse template %q: %v\n", name, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if err := tmpl.ExecuteTemplate(w, "base", val); err != nil {
		log.Printf("unable to execute template %q: %v\n", name, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func loadTmpl(tmplDir, name string) (*template.Template, error) {
	return template.New("").ParseFiles(path.Join(tmplDir, name), path.Join(tmplDir, baseTemplate))
}
