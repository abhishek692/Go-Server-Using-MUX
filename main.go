package main

import (
	"encoding/json"
	"fmt"
	"gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// creating a struct for books

type Book struct {
	ID     string  `json:"id"`
	ISBN   string  `json:"isbn"`
	TITLE  string  `json:"title"`
	AUTHOR *AUTHOR `json:"author"`
}

type AUTHOR struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

// initialize books variable as a slie of book type

var books []Book

// creating routeHandler functions. Any routeHandler function takes a request and a response

// get all Books function()

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

//get a single book

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get parameters from the request URL

	// finding the requested id from the different books we have in our list.

	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item) // item is getting encoded into w to be sent as response
			return
		}

	}
	json.NewEncoder(w).Encode(&Book{})

}

func createBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // setting the header to "application/json" type so that it gets interpretated as json file

	var book Book

	_ = json.NewDecoder(r.Body).Decode(&book)

	book.ID = strconv.Itoa(rand.Intn(10000000)) // creating anew ID for the book

	books = append(books, book)

	json.NewEncoder(w).Encode(book)

}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Conten-Type", "applicaton/json")

	params := mux.Vars(r) // extracting information from the request using mux package

	for index, item := range books {

		if item.ID == params["id"] {

			books = append(books[:index], books[index+1:]...) // deleting that book

			var book Book                             // creating a new book to update
			_ = json.NewDecoder(r.Body).Decode(&book) // decode the value of the book from the request, into book variable
			book.ID = params["id"]                    // since we are updating, we dont have to change the id
			books = append(books, book)               // appending the book to the books slice
			json.NewEncoder(w).Encode(book)           // encoding the book into the response
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	found := false

	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			found = true
			break
		}
	}

	if found == false {
		fmt.Println("book not found") // how to display on postman
		return
	}
	json.NewEncoder(w).Encode(books) // this means that the server will respond with all the books after deleting the book we requested
}

func main() {
	fmt.Println("Starting the server...")

	r := mux.NewRouter() // creating a new router using the mux pacakge

	// adding books in the slice
	books = append(books, Book{ID: "1", ISBN: "45673", TITLE: "Atomic Habits", AUTHOR: &AUTHOR{FirstName: "James", LastName: "Chadwick"}})
	books = append(books, Book{ID: "2", ISBN: "14573", TITLE: "Way of the superior man", AUTHOR: &AUTHOR{FirstName: "David", LastName: "Snyder"}})
	books = append(books, Book{ID: "3", ISBN: "90675", TITLE: "All along", AUTHOR: &AUTHOR{FirstName: "JK", LastName: "Singh"}})
	books = append(books, Book{ID: "4", ISBN: "12345", TITLE: "The Alchemist", AUTHOR: &AUTHOR{FirstName: "Naval", LastName: "Gupta"}})

	// creating route handlers to handle various http requests

	r.HandleFunc("/api/books", getBooks).Methods("GET")           // to get all the books
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")       // to get a particular book
	r.HandleFunc("/api/books", createBook).Methods("POST")        // to create a new book
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")    // to update a book
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE") // delete a book

	// starting the server

	log.Fatal(http.ListenAndServe(":8000", r)) // listening on port 8080 where r is the router and log.Fatal to indicate any problem if the server is unable to start
}

// go build && ./server.exe      run this command to build and run the server at once
