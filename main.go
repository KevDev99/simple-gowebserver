package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type Book struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

var Books []Book

func loadJson(searchString string) {
	data, err := ioutil.ReadFile("./data/books.json")

	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal([]byte(data), &Books)

	if err != nil {
		fmt.Println(err)
		return
	}

	if searchString != "" {
		var filteredBooks []Book

		for _, book := range Books {
			if strings.Contains(strings.ToLower(book.Title), strings.ToLower(searchString)) {
				filteredBooks = append(filteredBooks, book)
			}
		}

		Books = filteredBooks
	}

}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filters, present := query["search"]

	if present && len(filters) > 0 {
		loadJson(filters[0])
	} else {
		loadJson("")
	}

	tmpl, err := template.ParseFiles("./views/index.html")

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, Books)
}

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", serveTemplate)

	http.ListenAndServe(":8090", nil)

	// 1. web server
	// GET Books
	// POST BOoks
	// UPDATE Books
	// DELETE Books

}
