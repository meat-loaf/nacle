package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"html/template"
	"log"
	"path/filepath"
	"os"
	"strings"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

type Page struct {
	Title string
	Body []byte
}


var templ *template.Template

var md_render_flags = html.CommonFlags | html.HrefTargetBlank
var md_render_opts = html.RendererOptions{Flags: md_render_flags}
var md_renderer = html.NewRenderer(md_render_opts)

func md_render(html_template []byte) template.HTML{
	return template.HTML(markdown.ToHTML(html_template, nil, md_renderer))
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	print(title)
	filename := "content/" +title[1:] + ".txt"
	body, err := ioutil.ReadFile(filename)
	if (err != nil){
		fmt.Println(err)
		return nil, err
	}
	return &Page{Title:title, Body: body}, nil
}

func mainHandler(w http.ResponseWriter, r *http.Request){
	print("main handler")
	s := strings.ReplaceAll(r.URL.Path, "/", ".")
	pg, err := loadPage(s)
	if(err != nil){
		return
	}
	err = templ.ExecuteTemplate(w, s[1:]+".template.html", pg)
	//TODO seems the header is already written in this case. Why?
	if err != nil {
		fmt.Println(err);
	}
}

func root_handler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/" {
		err := templ.ExecuteTemplate(w, "main.template.html", Page{"main.", nil})
		if err != nil {
			fmt.Println(err);
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		templ.ExecuteTemplate(w, "404.template.html", Page{"ya dun goofed, buddy", nil});
	}
}

//from stack overflow...
func ParseTemplates(location string) *template.Template {
	templLocal := template.New("")
	//templLocal := template.Must(template.New("").Funcs(template.FuncMap{"markdown":md_renderer}))
	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".template.html"){
			templLocal = template.Must(
				templLocal.
				Funcs(template.FuncMap{"markdown":md_render}).
				ParseFiles(path))
		}
		return err
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return templLocal
}

func main() {
	templ = ParseTemplates("./templates")
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/", root_handler)
	//TODO handle these via reverse proxy with apache or some such instead, shouldnt serve these with go
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./style/"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
