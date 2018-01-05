package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"z/pagination/pkg/pagination"
)

type Question struct {
	Title   string
	Content string
}

func main() {
	mux := http.NewServeMux()

	// static file(js css...)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/", show)

	log.Println("Starting server on :8001")
	err := http.ListenAndServe(":8001", mux)
	log.Fatal(err)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func show(w http.ResponseWriter, r *http.Request) {
	t := "./views/index.html"
	p := "./views/pager.html"

	questions := getData()
	rows := len(questions)

	pager := pagination.NewPage(r, rows)
	start, end := pager.StartEnd()

	var viewData struct {
		Q          []*Question
		P          *pagination.Page
		Start, End int
	}

	viewData.Q = questions[start-1 : end]
	viewData.P = pager
	viewData.Start, viewData.End = start, end

	tmpl, err := template.ParseFiles(t, p)
	if err != nil {
		log.Println(err.Error())
	}
	err = tmpl.ExecuteTemplate(w, "index.html", viewData)
	if err != nil {
		log.Println(err.Error())
	}
}

func getData() []*Question {
	questions := make([]*Question, 0)

	for i := 1; i <= 100; i++ {
		title := "title-" + strconv.Itoa(i)
		content := "content-" + strconv.Itoa(i)
		q := new(Question)
		q.Title = title
		q.Content = content
		questions = append(questions, q)
	}

	return questions
}
