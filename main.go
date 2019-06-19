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
)

type Page struct {
	Title string
	Body []byte
}


var templ *template.Template

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if (err != nil){
		return nil, err;
	}
	return &Page{Title:title, Body: body}, nil
}

func mainHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println(r.URL.Path);
	s := strings.ReplaceAll(r.URL.Path, "/", ".")
	fmt.Println(s)
	fmt.Println(s[1:])
	test := Page {s[1:], nil}
	fmt.Println(test)
	fmt.Println(w)

	err := templ.ExecuteTemplate(w, s[1:]+".template.html", Page{s, nil})
	//TODO seems the header is already written in this case. Why?
	if err != nil {
		fmt.Println(err)
		templ.ExecuteTemplate(w, "404.template.html", Page{"ya dun goofed, buddy", nil});
		w.WriteHeader(http.StatusNotFound)
	}
}
//from stack overflow...
func ParseTemplates(location string) *template.Template {
	templLocal := template.New("")
	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".template.html"){
			_, err = templLocal.ParseFiles(path)
			if err != nil {
				;//do something i guess?
			}
		}
		return err
	})
	if err != nil {
		panic(err)
	}
	return templLocal
}

func main() {
	templ = ParseTemplates("./templates")
	http.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
