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
	s := strings.ReplaceAll(r.URL.Path, "/", ".")
	err := templ.ExecuteTemplate(w, s[1:]+".template.html", Page{s, nil})
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
	err := filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".template.html"){
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
	http.HandleFunc("/main", mainHandler)
	http.HandleFunc("/", root_handler)
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
