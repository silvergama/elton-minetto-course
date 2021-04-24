package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/silvergama/studying-golang/elton-minetto-course/app/handlers"
	"github.com/silvergama/studying-golang/elton-minetto-course/core/beer"
	"github.com/urfave/negroni"
)

func main() {
	db, err := sql.Open("sqlite3", "data/beer.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	service := beer.NewService(db)

	r := mux.NewRouter()

	// midlewares -  código que vai ser executado em todas as requests
	// aqui podemos colocar logs, inclusão e validação de cabeçalho, etc.
	n := negroni.New(
		negroni.NewLogger(),
	)

	// handlers
	handlers.MakeBeerHandlers(r, n, service)

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
