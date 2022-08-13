// Сервер должен быть доступен по адресу: http://localhost:8080.
// Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
// Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
// Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
// Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	mux := http.NewServeMux()
	var err error
	connStr := "user=postgres password=200112 dbname=forURL sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db.Stats().OpenConnections)

	mux.HandleFunc("/post", handlerF)
	mux.HandleFunc("/get/id", handlerF1)
	mux.HandleFunc("/", handlerF2)

	http.ListenAndServe("localhost:8080", mux)
}

func handlerF2(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO!!!"))
	fmt.Println("HELLO FROM HF2")
	// q := r.URL.RawQuery
	// http.Redirect(w, r, q, http.StatusTemporaryRedirect)
}

type URLs struct {
	fullURL string
	cutURL  string
}

func checkErr(err error) error {
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func handlerF1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HELLO FROM HF1", time.Now())
	query := r.URL.RawQuery

	switch {
	case r.Method == http.MethodGet:
		rows, err := db.Query("SELECT fullURL, ID FROM URLS WHERE ID = $1", query)
		checkErr(err)
		var URLs = URLs{}
		for rows.Next() {
			err = rows.Scan(&URLs.fullURL, &URLs.cutURL)

			checkErr(err)
			if query == URLs.cutURL {
				// http.Redirect(w, r, URLs.fullURL, http.StatusTemporaryRedirect)
				w.Header().Set("Location", URLs.fullURL)
				// w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}
		http.Error(w, "Not found", http.StatusNotFound)
	}
	http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
}

func handlerF(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HELLO FROM HF", time.Now())

	switch {
	case r.Method == http.MethodPost:
		rows, err := db.Query("SELECT fullURL, ID FROM URLS WHERE fullURL = $1", r.Header.Get("URL"))
		checkErr(err)
		var URLs = URLs{}
		for rows.Next() {
			err = rows.Scan(&URLs.fullURL, &URLs.cutURL)
			checkErr(err)
			if URLs.fullURL == r.Header.Get("URL") {
				checkErr(err)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("localhost:8080/" + URLs.cutURL))
				return
			}
		}
		res, err := db.Exec("INSERT INTO URLS (fullURL, ID) VALUES ($1, $2)",
			r.Header.Get("URL"),
			RandStringBytes(5),
		)
		checkErr(err)
		_ = res
		nextNewRows, err := db.Query("SELECT ID FROM URLS")
		checkErr(err)
		for nextNewRows.Next() {
			err = nextNewRows.Scan(&URLs.cutURL)
			checkErr(err)
		}
		nextNewRows.Scan(&URLs.cutURL)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("localhost:8080/" + URLs.cutURL))
	default:
		http.Error(w, "Only POST requests are allowed!", http.StatusUnauthorized)
	}
}
