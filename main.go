package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles(
	"templates/notFound.html",
	"templates/header.html",
	"templates/footer.html",
	"templates/index.html",
	"templates/reports.html"))

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

//Sewrve the home page (index)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := Page{
		PageTitle: "Home",
	}
	display(w, "index", data)
}

func reportsHandler(w http.ResponseWriter, r *http.Request) {
	data := TodoPageData{
		PageTitle: "TODO List",
		ListTitle: "TODO List",
		Todos: []Todo{
			{Title: "Project Setup", Done: true},
			{Title: "Setup CI", Done: false},
			{Title: "Setup CD", Done: true},
			{Title: "Make Test Cases", Done: false},
			{Title: "Flushout Navbar", Done: false},
			{Title: "Create 404", Done: true},
			{Title: "Add How to Use Page", Done: false},
		},
	}
	display(w, "todo", data)
}

func testPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TestPage!")
}

func randomPageHandler(w http.ResponseWriter, r *http.Request) {
	// If empty show the home page
	// If a number simulate a dice with that many sides
	// If static page show the static package
	// Else show the 404 page
	if r.URL.Path == "/" {
		homeHandler(w, r)
	} else if strings.HasSuffix(r.URL.Path[1:], ".html") {
		http.ServeFile(w, r, "static/html/"+r.URL.Path[1:])
	} else {
		fmt.Println("Sorry but it seems this page does not exist...")
		errorHandler(w, r, http.StatusNotFound)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		display(w, "404", &Page{PageTitle: "404"})
	} else {
		http.ServeFile(w, r, "static/html/issue.html")
	}
}

func getPort() string {
	if value, ok := os.LookupEnv("PORT"); ok {
		return ":" + value
	}
	return ":8080"
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/custom/", testPageHandler)
	http.HandleFunc("/reports/", reportsHandler)

	http.HandleFunc("/", randomPageHandler)

	port := getPort()
	fmt.Println("Now listening to port " + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Page structure
type Page struct {
	PageTitle string
}

// Todo struct
type Todo struct {
	Title string
	Done  bool
}

// TodoPageData with titles and a list of Todos
type TodoPageData struct {
	PageTitle string
	ListTitle string
	Todos     []Todo
}
