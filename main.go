package main

import (
	"log"
	"net/http"
)

func main() {
	st, err := NewStore("data")
	if err != nil {
		log.Fatal(err)
	}
	app := &App{St: st}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/resources", http.StatusSeeOther)
	})
	mux.HandleFunc("/resources", app.ResourcesPage)
	mux.HandleFunc("/sources", app.SourcesPage)
	mux.HandleFunc("/history", app.HistoryPage)
	mux.HandleFunc("/summary", app.SummaryPage)
	mux.HandleFunc("/details", app.DetailsPage)

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
