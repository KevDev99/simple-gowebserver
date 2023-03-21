package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Book struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

var Books []Book

var temp *template.Template

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

func saveJson(json string, filePath string) {

	fmt.Printf("%v", json)

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer file.Close()

	ioutil.WriteFile(filePath, []byte(json), 0644)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filters, present := query["search"]

	if present && len(filters) > 0 {
		loadJson(filters[0])
	} else {
		loadJson("")
	}

	tmpl, err := template.ParseFiles("./views/index.html", "./views/header.html")

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"Books": Books, "Search": ""}

	if len(filters) > 0 {
		data["Search"] = filters[0]
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		log.Fatal(err)
	}
}

func addBookTemplate(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("./views/add-book.html", "./views/header.html")

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, Books)

	if err != nil {
		log.Fatal(err)
	}
}

func apiAddBook(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	hasTitle := r.PostForm.Has("title")
	hasThumbnailUrl := r.PostForm.Has("thumbnailUrl")

	if !hasTitle || !hasThumbnailUrl {
		http.Error(w, "Title or Thumbnail not provided.", http.StatusBadRequest)
		return
	}

	// get current books from json
	loadJson("")

	// append
	newBook := Book{Id: rand.Intn(100000), Title: r.PostFormValue("title"), ThumbnailUrl: r.PostFormValue("thumbnailUrl")}

	fmt.Printf("new book: %v", newBook)

	Books = append(Books, newBook)

	jsonBooks, err := json.Marshal(Books)

	if err != nil {
		log.Fatal(err)
	}

	saveJson(string(jsonBooks), "./data/books.json")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/add-booking", addBookTemplate)

	/* API endpoints */
	http.HandleFunc("/api/add-new-booking", apiAddBook)

	http.ListenAndServe(":8090", nil)

	// 1. web server
	// GET Books
	// POST BOoks
	// UPDATE Books
	// DELETE Books

}
