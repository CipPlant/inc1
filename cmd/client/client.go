package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	post()
	get()
}

func get() {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/id?c0e2fd12-1105-4cbf-b8d8-99881602ad25", nil)

	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Header.Get("Location"))
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
}

func post() {
	client := http.Client{}
	var URL string
	fmt.Scanf("%s\n", &URL)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("URL", URL)
	fmt.Println("!!!!!!!!!!!!!!!!", URL)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)
	fmt.Println(req.Header)
}
