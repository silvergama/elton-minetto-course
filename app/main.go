package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/silvergama/studying-golang/elton-minetto-course/core/beer"
	"github.com/urfave/negroni"
)

func main() {
	db, err := sql.Open("sqlite3", "../data/beer.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	service := beer.NewService(db)

	r := mux.NewRouter()

	// handlers
	n := negroni.New(
		negroni.NewLogger(),
	)

	r.Handle("/v1/beer", n.With(
		negroni.Wrap(hello(service)),
	)).Methods("GET", "OPTIONS")

	http.Handle("/", r)

	svr := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":8000",
		Handler:      http.DefaultServeMux,
		ErrorLog:     log.New(os.Stderr, "logger: ", log.Lshortfile),
	}

	err = svr.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func hello(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		all, _ := service.GetAll()
		for _, i := range all {
			fmt.Println(i)
		}
	})
}
