package main

import (
	// "errors"
	// "fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// we'll build a library api
// - store a bunch of books
// - check in a book
// - check out a book
// - view all books
// - get books by id

type book struct {
	ID       string `json:"id"` // json field initializer
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// dummy data
var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}
func getBOoks(c *gin.Context){
	c.IndentedJSON(http.StatusOK, books);
}

func createBook(c *gin.Context){
	 var newBook book
	 if err:= c.BindJSON(&newBook); err != nil{
		return
	 }

	 books = append(books, newBook)
	 c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBOoks)
	router.POST("/books", createBook)
	router.Run("localhost:8080")
}
